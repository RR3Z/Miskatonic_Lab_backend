param(
    [ValidateSet("setup", "local", "e2e", "clerk", "migrations", "all")]
    [string]$Mode = "all"
)

$ErrorActionPreference = "Stop"
$ProjectRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$EnvPath = Join-Path $ProjectRoot ".env"
$RunDirectory = Join-Path ([System.IO.Path]::GetTempPath()) ("miskatonic-tests-" + [guid]::NewGuid())
$BackendProcess = $null
$TunnelProcess = $null
$DevTunnelExecutable = $null
$RunStartedAt = Get-Date
$RunSucceeded = $false

function Import-EnvironmentFile {
    if (-not (Test-Path $EnvPath)) {
        throw ".env is required; copy .env.example and fill the documented values"
    }

    Get-Content $EnvPath | ForEach-Object {
        $line = $_.Trim()
        if ($line -eq "" -or $line.StartsWith("#") -or -not $line.Contains("=")) {
            return
        }
        $name, $value = $line.Split("=", 2)
        $name = $name.Trim()
        if ($name -and -not [Environment]::GetEnvironmentVariable($name, "Process")) {
            [Environment]::SetEnvironmentVariable($name, $value.Trim().Trim('"').Trim("'"), "Process")
        }
    }
}

function Require-Environment {
    param([string[]]$Names)

    foreach ($name in $Names) {
        if ([string]::IsNullOrWhiteSpace([Environment]::GetEnvironmentVariable($name, "Process"))) {
            switch ($name) {
                "TEST_BACKEND_PORT" { throw "TEST_BACKEND_PORT must be set in .env (normally TEST_BACKEND_PORT=8001)" }
                "DEVTUNNEL_TUNNEL_ID" { throw "DEVTUNNEL_TUNNEL_ID must be set in .env. Create one with: devtunnel create --allow-anonymous" }
                "CLERK_WEBHOOK_PUBLIC_URL" { throw "CLERK_WEBHOOK_PUBLIC_URL must be set in .env to the HTTPS URL for the configured Dev Tunnel port" }
                "MIGRATION_SMOKE_DATABASE_URL" { throw "MIGRATION_SMOKE_DATABASE_URL must be set in .env to a separate loopback database ending in _migration_smoke_test" }
            }
            throw "$name must be set in .env"
        }
    }
}

function Require-Command {
    param([string]$Name, [string]$InstallHint)

    if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
        throw "$Name is required. $InstallHint"
    }
}

function Ensure-DevTunnelCLI {
    if ($script:DevTunnelExecutable) {
        return
    }

    $existing = Get-Command "devtunnel" -ErrorAction SilentlyContinue
    if ($existing) {
        $script:DevTunnelExecutable = $existing.Source
        return
    }

    $installDirectory = Join-Path $env:LOCALAPPDATA "MiskatonicLab\tools"
    $downloadPath = Join-Path $installDirectory "devtunnel.exe"
    if (-not (Test-Path $downloadPath)) {
        New-Item -ItemType Directory -Path $installDirectory -Force | Out-Null
        Write-Output "Dev Tunnel CLI is not installed; downloading it to $installDirectory"
        Invoke-WebRequest -Uri "https://aka.ms/TunnelsCliDownload/win-x64" -OutFile $downloadPath
    }
    $script:DevTunnelExecutable = $downloadPath
}

function Invoke-DevTunnel {
    param([string[]]$Arguments)

    Ensure-DevTunnelCLI
    & $script:DevTunnelExecutable @Arguments
}

function Ensure-DevTunnelLogin {
    Ensure-DevTunnelCLI
    $status = (& $script:DevTunnelExecutable user show 2>&1 | Out-String)
    $loggedIn = $LASTEXITCODE -eq 0 -and $status -notmatch "(?i)not\s+logged|sign\s+in|no\s+user"
    if ($loggedIn) {
        return
    }

    Write-Output "Dev Tunnel login is required; opening the interactive login flow."
    & $script:DevTunnelExecutable user login
    Assert-LastExitCode -Operation "Dev Tunnel login"
    $status = (& $script:DevTunnelExecutable user show 2>&1 | Out-String)
    if ($LASTEXITCODE -ne 0 -or $status -match "(?i)not\s+logged|sign\s+in|no\s+user") {
        throw "Dev Tunnel login did not complete; run: devtunnel user login"
    }
}

function Set-EnvironmentFileValue {
    param([string]$Name, [string]$Value)

    $lines = @(Get-Content $EnvPath)
    $pattern = "^\s*" + [regex]::Escape($Name) + "="
    $replacement = "$Name=$Value"
    $updated = $false
    for ($index = 0; $index -lt $lines.Count; $index++) {
        if ($lines[$index] -match $pattern) {
            $lines[$index] = $replacement
            $updated = $true
            break
        }
    }
    if (-not $updated) {
        $lines += $replacement
    }
    [System.IO.File]::WriteAllLines($EnvPath, $lines, [System.Text.UTF8Encoding]::new($false))
    [Environment]::SetEnvironmentVariable($Name, $Value, "Process")
}

function Test-PlaceholderTunnelValue {
    param([string]$Value)

    return [string]::IsNullOrWhiteSpace($Value) -or $Value.Trim() -match "(?i)replace_me|<[^>]+>"
}

function Assert-LastExitCode {
    param([string]$Operation)

    if ($LASTEXITCODE -ne 0) {
        throw "$Operation failed with exit code $LASTEXITCODE"
    }
}

function Write-TestStage {
    param([string]$Name)

    Write-Host "`n[$Name]"
}

function Invoke-QuietCommand {
    param(
        [string]$Operation,
        [string]$FilePath,
        [string[]]$Arguments
    )

    $previousErrorActionPreference = $ErrorActionPreference
    $ErrorActionPreference = "Continue"
    try {
        $output = @(& $FilePath @Arguments 2>&1)
    }
    finally {
        $ErrorActionPreference = $previousErrorActionPreference
    }
    if ($LASTEXITCODE -ne 0) {
        Write-Host "  FAIL  $Operation"
        Write-Host "`n--- command output ---"
        $output | ForEach-Object { Write-Host $_ }
        throw "$Operation failed with exit code $LASTEXITCODE"
    }
    return $output
}

function Invoke-GoTestSuite {
    param([string]$Name, [string[]]$Packages)

    $watch = [System.Diagnostics.Stopwatch]::StartNew()
    $previousErrorActionPreference = $ErrorActionPreference
    $ErrorActionPreference = "Continue"
    try {
        $output = @(& go test @Packages 2>&1)
    }
    finally {
        $ErrorActionPreference = $previousErrorActionPreference
    }
    $exitCode = $LASTEXITCODE
    $watch.Stop()
    if ($exitCode -ne 0) {
        Write-Host "  FAIL  $Name"
        Write-Host "`n--- go test output ---"
        $output | ForEach-Object { Write-Host $_ }
        throw "$Name failed with exit code $exitCode"
    }

    $packageCount = @($output | Where-Object { "$_" -match "^ok\s+" }).Count
    Write-Host ("  PASS  {0} | {1} test packages | {2:N1}s" -f $Name, $packageCount, $watch.Elapsed.TotalSeconds)
}

function Stop-ProcessTree {
    param([int]$ProcessId)

    $children = @(Get-CimInstance Win32_Process -Filter "ParentProcessId = $ProcessId" -ErrorAction SilentlyContinue)
    foreach ($child in $children) {
        Stop-ProcessTree -ProcessId $child.ProcessId
    }
    Stop-Process -Id $ProcessId -Force -ErrorAction SilentlyContinue
}

function Test-PortAvailable {
    param([int]$Port)

    if (Get-NetTCPConnection -LocalPort $Port -State Listen -ErrorAction SilentlyContinue) {
        throw "localhost:$Port is already in use; the test runner will not stop a process it did not create"
    }
}

function Get-TestBackendPort {
    Require-Environment @("TEST_BACKEND_PORT")
    $port = 0
    if (-not [int]::TryParse($env:TEST_BACKEND_PORT, [ref]$port) -or $port -lt 1 -or $port -gt 65535) {
        throw "TEST_BACKEND_PORT must be an integer between 1 and 65535"
    }
    return $port
}

function Validate-TestConfiguration {
    param([bool]$RequireE2E, [bool]$RequireClerk, [bool]$RequireSmoke, [bool]$AllowTunnelBootstrap = $false)

    Require-Environment @("TEST_DATABASE_URL")
    Require-Command "go" "Install Go and add it to PATH."
    Require-Command "docker" "Install Docker Desktop and start its engine."
    Require-Command "migrate" "Install golang-migrate CLI and add it to PATH."
    if ($RequireSmoke) {
        Require-Environment @("MIGRATION_SMOKE_DATABASE_URL")
        $null = Invoke-QuietCommand -Operation "validate migration smoke database" -FilePath "go" -Arguments @("run", "./cmd/testdb-admin", "-validate")
    }

    $requiresBackend = $RequireE2E -or $RequireClerk
    if (-not $requiresBackend) {
        return
    }

    Require-Environment @(
        "CLERK_SECRET_KEY",
        "CLERK_WEBHOOK_SIGNING_SECRET",
        "CLERK_AUTHORIZED_PARTIES",
        "TEST_BACKEND_PORT"
    )

    if ($RequireE2E) {
        Require-Environment @("E2E_TEST1_MAIL", "E2E_TEST2_MAIL")
    }

    if (-not $RequireClerk) {
        return
    }

    Ensure-DevTunnelCLI
    if (-not $AllowTunnelBootstrap) {
        Require-Environment @("DEVTUNNEL_TUNNEL_ID", "CLERK_WEBHOOK_PUBLIC_URL")
    }
}

function Initialize-TestDatabase {
    Write-TestStage "Test database"
    $null = Invoke-QuietCommand -Operation "start postgres-test" -FilePath "docker" -Arguments @("compose", "up", "-d", "--wait", "postgres-test")
    $null = Invoke-QuietCommand -Operation "prepare test database" -FilePath "powershell" -Arguments @("-NoProfile", "-ExecutionPolicy", "Bypass", "-File", "./scripts/testdb.ps1", "prepare")
    Write-Host "  PASS  Docker test database is ready"
}

function Reset-MigrationSmokeDatabase {
    Write-TestStage "Migration smoke database"
    $null = Invoke-QuietCommand -Operation "reset migration smoke database" -FilePath "go" -Arguments @("run", "./cmd/testdb-admin", "-reset")
    Write-Host "  PASS  Disposable smoke database was reset"
}

function Start-TestBackend {
    param([int]$Port)

    Test-PortAvailable -Port $Port
    New-Item -ItemType Directory -Path $RunDirectory -Force | Out-Null
    $storageDirectory = Join-Path $RunDirectory "portraits"
    New-Item -ItemType Directory -Path $storageDirectory -Force | Out-Null

    $env:DATABASE_URL = $env:TEST_DATABASE_URL
    $env:PORT = [string]$Port
    $env:E2E_BASE_URL = "http://localhost:$Port"
    $env:CLERK_LOCAL_BACKEND_URL = $env:E2E_BASE_URL
    $env:PUBLIC_BACKEND_URL = $env:E2E_BASE_URL
    $env:PORTRAIT_STORAGE_DIR = $storageDirectory

    $stdout = Join-Path $RunDirectory "backend.stdout.log"
    $stderr = Join-Path $RunDirectory "backend.stderr.log"
    $script:BackendProcess = Start-Process -FilePath "go" -ArgumentList @("run", "./cmd") -WorkingDirectory $ProjectRoot -WindowStyle Hidden -PassThru -RedirectStandardOutput $stdout -RedirectStandardError $stderr

    $deadline = (Get-Date).AddSeconds(30)
    while ((Get-Date) -lt $deadline) {
        if ($script:BackendProcess.HasExited) {
            $output = if (Test-Path $stdout) { Get-Content $stdout -Raw } else { "" }
            $errors = if (Test-Path $stderr) { Get-Content $stderr -Raw } else { "" }
            throw "temporary backend exited before listening.`n$output`n$errors"
        }
        if (Get-NetTCPConnection -LocalPort $Port -State Listen -ErrorAction SilentlyContinue) {
            return
        }
        Start-Sleep -Milliseconds 250
    }
    throw "temporary backend did not listen on localhost:$Port within 30 seconds"
}

function Assert-PersistentTunnelConfiguration {
    $port = Get-TestBackendPort
    $tunnelID = $env:DEVTUNNEL_TUNNEL_ID.Trim()
    $publicURL = [Uri]$env:CLERK_WEBHOOK_PUBLIC_URL.Trim()
    if ($publicURL.Scheme -ne "https" -or -not $publicURL.Host.Contains("-$port.") -or -not $publicURL.Host.EndsWith(".devtunnels.ms")) {
        throw "CLERK_WEBHOOK_PUBLIC_URL must be an HTTPS Dev Tunnel URL for TEST_BACKEND_PORT"
    }

    Ensure-DevTunnelLogin
    $null = Invoke-QuietCommand -Operation "check persistent Dev Tunnel port" -FilePath $script:DevTunnelExecutable -Arguments @("port", "show", $tunnelID, "-p", $port)
}

function Get-DefaultDevTunnelID {
    Ensure-DevTunnelLogin
    $output = [string]((& $script:DevTunnelExecutable show 2>&1 | Out-String))
    if ($LASTEXITCODE -ne 0) {
        return $null
    }
    $match = [regex]::Match($output, "(?im)^Tunnel ID\s*:\s*(\S+)")
    if ($match.Success) {
        return $match.Groups[1].Value
    }
    return $null
}

function Resolve-DevTunnelPublicURL {
    param([string]$TunnelID, [int]$Port)

    New-Item -ItemType Directory -Path $RunDirectory -Force | Out-Null
    $stdout = Join-Path $RunDirectory "tunnel-setup.stdout.log"
    $stderr = Join-Path $RunDirectory "tunnel-setup.stderr.log"
    $setupTunnel = Start-Process -FilePath $script:DevTunnelExecutable -ArgumentList @("host", $TunnelID) -WorkingDirectory $ProjectRoot -WindowStyle Hidden -PassThru -RedirectStandardOutput $stdout -RedirectStandardError $stderr
    try {
        $deadline = (Get-Date).AddSeconds(15)
        while ((Get-Date) -lt $deadline) {
            $output = if (Test-Path $stdout) { [string](Get-Content $stdout -Raw) } else { "" }
            if (Test-Path $stderr) {
                $output += [string](Get-Content $stderr -Raw)
            }
            $match = [regex]::Match($output, "https://[^\s/]+-$Port\.[^\s/]+devtunnels\.ms")
            if ($match.Success) {
                return $match.Value
            }
            if ($setupTunnel.HasExited) {
                throw "Dev Tunnel exited while resolving its public URL.`n$output"
            }
            Start-Sleep -Milliseconds 250
        }
        throw "Dev Tunnel did not report a public URL within 15 seconds"
    }
    finally {
        Stop-ProcessTree -ProcessId $setupTunnel.Id
    }
}

function Initialize-PersistentTestTunnel {
    $tunnelID = [Environment]::GetEnvironmentVariable("DEVTUNNEL_TUNNEL_ID", "Process")
    $publicURL = [Environment]::GetEnvironmentVariable("CLERK_WEBHOOK_PUBLIC_URL", "Process")
    $hasTunnelID = -not (Test-PlaceholderTunnelValue -Value $tunnelID)
    $hasPublicURL = -not (Test-PlaceholderTunnelValue -Value $publicURL)
    if ($hasTunnelID -and $hasPublicURL) {
        try {
            Assert-PersistentTunnelConfiguration
            return
        }
        catch {
            if ($_.Exception.Message -notmatch "check persistent Dev Tunnel port") {
                throw
            }
            $defaultTunnelID = Get-DefaultDevTunnelID
            if ($defaultTunnelID -and $defaultTunnelID -ne $tunnelID) {
                Write-Warning "Configured Dev Tunnel '$tunnelID' is stale; recovering the newly created default tunnel '$defaultTunnelID'."
                $tunnelID = $defaultTunnelID
                $publicURL = $null
                $hasTunnelID = $true
                $hasPublicURL = $false
            }
            else {
                Write-Warning "Configured Dev Tunnel '$tunnelID' no longer exists or has no port $env:TEST_BACKEND_PORT; provisioning a replacement."
                $hasTunnelID = $false
                $hasPublicURL = $false
            }
        }
    }

    if (-not $hasTunnelID -and $hasPublicURL) {
        throw "CLERK_WEBHOOK_PUBLIC_URL is configured but DEVTUNNEL_TUNNEL_ID is empty; clear both values and run test:setup again"
    }

    Ensure-DevTunnelLogin
    $port = Get-TestBackendPort
    if (-not $hasTunnelID) {
        $tunnelID = "miskatonic" + [guid]::NewGuid().ToString("N").Substring(0, 12)
        Invoke-DevTunnel -Arguments @("create", $tunnelID, "--allow-anonymous")
        Assert-LastExitCode -Operation "create persistent Dev Tunnel"
        Invoke-DevTunnel -Arguments @("port", "create", $tunnelID, "-p", $port, "--protocol", "http")
        Assert-LastExitCode -Operation "create Dev Tunnel HTTP port"
    }

    Set-EnvironmentFileValue -Name "DEVTUNNEL_TUNNEL_ID" -Value $tunnelID
    $publicURL = Resolve-DevTunnelPublicURL -TunnelID $tunnelID -Port $port
    Set-EnvironmentFileValue -Name "CLERK_WEBHOOK_PUBLIC_URL" -Value $publicURL
    Write-Output "Created persistent Dev Tunnel. Configure Clerk Dashboard webhook URL: $publicURL/webhooks/clerk/user"
}

function Start-TestTunnel {
    New-Item -ItemType Directory -Path $RunDirectory -Force | Out-Null
    $stdout = Join-Path $RunDirectory "tunnel.stdout.log"
    $stderr = Join-Path $RunDirectory "tunnel.stderr.log"
    $script:TunnelProcess = Start-Process -FilePath $script:DevTunnelExecutable -ArgumentList @("host", $env:DEVTUNNEL_TUNNEL_ID) -WorkingDirectory $ProjectRoot -WindowStyle Hidden -PassThru -RedirectStandardOutput $stdout -RedirectStandardError $stderr

    Start-Sleep -Seconds 2
    if ($script:TunnelProcess.HasExited) {
        $output = if (Test-Path $stdout) { Get-Content $stdout -Raw } else { "" }
        $errors = if (Test-Path $stderr) { Get-Content $stderr -Raw } else { "" }
        throw "Dev Tunnel exited before hosting the configured tunnel.`n$output`n$errors"
    }
}

function Run-LocalTests {
    $env:E2E_TESTS = ""
    $env:CLERK_INTEGRATION_TESTS = ""
    $env:MIGRATION_SMOKE_TESTS = ""
    Write-TestStage "Local Go tests"
    Invoke-GoTestSuite -Name "Local Go tests" -Packages @("./...")
}

function Run-ClerkTests {
    $env:CLERK_INTEGRATION_TESTS = "1"
    Write-TestStage "Clerk webhook integration"
    Invoke-GoTestSuite -Name "Clerk webhook integration" -Packages @("./tests/user/integration")
}

function Run-E2ETests {
    $env:E2E_TESTS = "1"
    Write-TestStage "Live HTTP and WebSocket E2E"
    Invoke-GoTestSuite -Name "Live HTTP and WebSocket E2E" -Packages @("./tests/e2e")
}

function Run-MigrationSmokeTests {
    $env:MIGRATION_SMOKE_TESTS = "1"
    Write-TestStage "Migration rollback smoke"
    Invoke-GoTestSuite -Name "Migration rollback smoke" -Packages @("./tests/migrations")
}

Push-Location $ProjectRoot
try {
    Import-EnvironmentFile

    if ($Mode -eq "setup") {
        Write-Host "Miskatonic Lab test setup"
        Write-TestStage "Configuration"
        Validate-TestConfiguration -RequireE2E $true -RequireClerk $true -RequireSmoke $true -AllowTunnelBootstrap $true
        Initialize-PersistentTestTunnel
        Assert-PersistentTunnelConfiguration
        Write-Host "  PASS  Setup is ready; ensure Clerk Dashboard uses CLERK_WEBHOOK_PUBLIC_URL/webhooks/clerk/user"
        $script:RunSucceeded = $true
        exit 0
    }

    Write-Host "Miskatonic Lab test runner - $Mode"
    $requiresE2E = $Mode -in @("e2e", "all")
    $requiresClerk = $Mode -in @("clerk", "all")
    $requiresSmoke = $Mode -in @("migrations", "all")
    Write-TestStage "Configuration"
    Validate-TestConfiguration -RequireE2E $requiresE2E -RequireClerk $requiresClerk -RequireSmoke $requiresSmoke
    if ($requiresClerk) {
        Assert-PersistentTunnelConfiguration
    }
    Write-Host "  PASS  Environment and tools are ready"
    Initialize-TestDatabase

    if ($Mode -eq "migrations" -or $Mode -eq "all") {
        Reset-MigrationSmokeDatabase
    }

    if ($Mode -eq "local" -or $Mode -eq "all") {
        Run-LocalTests
    }

    if ($requiresE2E -or $requiresClerk) {
        Write-TestStage "Isolated test backend"
        $backendPort = Get-TestBackendPort
        Start-TestBackend -Port $backendPort
        Write-Host "  PASS  Backend listens on localhost:$backendPort"
    }

    if ($Mode -eq "clerk" -or $Mode -eq "all") {
        Write-TestStage "Dev Tunnel"
        Start-TestTunnel
        Write-Host "  PASS  Persistent tunnel is hosted"
        Run-ClerkTests
    }

    if ($Mode -eq "e2e" -or $Mode -eq "all") {
        Run-E2ETests
    }

    if ($Mode -eq "migrations" -or $Mode -eq "all") {
        Run-MigrationSmokeTests
    }
    $script:RunSucceeded = $true
}
finally {
    if ($null -ne $TunnelProcess) {
        Stop-ProcessTree -ProcessId $TunnelProcess.Id
    }
    if ($null -ne $BackendProcess) {
        Stop-ProcessTree -ProcessId $BackendProcess.Id
    }
    Remove-Item -LiteralPath $RunDirectory -Recurse -Force -ErrorAction SilentlyContinue
    Pop-Location
    if ($script:RunSucceeded) {
        $duration = (Get-Date) - $RunStartedAt
        $durationText = $duration.ToString('m\:ss')
        Write-Host ("`n[OK] {0} completed in {1}" -f $Mode, $durationText)
    }
}
