param(
    [Parameter(Mandatory = $true, Position = 0)]
    [ValidateSet("down", "up", "up-all", "version")]
    [string]$Command,

    [Parameter(Position = 1)]
    [int]$Steps
)

$ErrorActionPreference = "Stop"

$ProjectRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$MigrationsPath = "migrations"
$EnvPath = Join-Path $ProjectRoot ".env"

if (Test-Path $EnvPath) {
    Get-Content $EnvPath | ForEach-Object {
        $line = $_.Trim()

        if ($line -eq "" -or $line.StartsWith("#") -or -not $line.Contains("=")) {
            return
        }

        $name, $value = $line.Split("=", 2)
        $name = $name.Trim()
        $value = $value.Trim().Trim('"').Trim("'")

        if ($name -and -not [Environment]::GetEnvironmentVariable($name, "Process")) {
            [Environment]::SetEnvironmentVariable($name, $value, "Process")
        }
    }
}

if ($env:DATABASE_URL) {
    $DatabaseUrl = $env:DATABASE_URL
} else {
    $PostgresHost = if ($env:POSTGRES_HOST) { $env:POSTGRES_HOST } else { "localhost" }
    $PostgresPort = if ($env:POSTGRES_PORT) { $env:POSTGRES_PORT } else { "5432" }
    $PostgresDb = if ($env:POSTGRES_DB) { $env:POSTGRES_DB } else { "miskatonic_lab" }
    $PostgresUser = if ($env:POSTGRES_USER) { $env:POSTGRES_USER } else { "miskatonic_user" }
    $PostgresPassword = if ($env:POSTGRES_PASSWORD) { $env:POSTGRES_PASSWORD } else { "miskatonic_password" }
    $PostgresSslMode = if ($env:POSTGRES_SSLMODE) { $env:POSTGRES_SSLMODE } else { "disable" }

    $DatabaseUrl = "postgres://$PostgresUser`:$PostgresPassword@$PostgresHost`:$PostgresPort/$PostgresDb`?sslmode=$PostgresSslMode"
}

if (-not (Get-Command migrate -ErrorAction SilentlyContinue)) {
    Write-Error "migrate command not found. Install golang-migrate CLI and make sure it is available in PATH."
}

function Invoke-Migrate {
    param(
        [Parameter(ValueFromRemainingArguments = $true)]
        [string[]]$Arguments
    )

    Push-Location $ProjectRoot
    try {
        & migrate -path $MigrationsPath -database $DatabaseUrl @Arguments
        return $LASTEXITCODE
    } finally {
        Pop-Location
    }
}

switch ($Command) {
    "down" {
        if ($Steps -le 0) {
            Write-Error "Usage: .\scripts\migrate.ps1 down <steps>"
        }

        $ExitCode = Invoke-Migrate down $Steps
    }
    "up" {
        if ($Steps -le 0) {
            Write-Error "Usage: .\scripts\migrate.ps1 up <steps>"
        }

        $ExitCode = Invoke-Migrate up $Steps
    }
    "up-all" {
        $ExitCode = Invoke-Migrate up
    }
    "version" {
        $ExitCode = Invoke-Migrate version
    }
}

exit $ExitCode
