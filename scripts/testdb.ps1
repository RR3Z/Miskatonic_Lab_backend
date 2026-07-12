param(
    [Parameter(Mandatory = $true, Position = 0)]
    [ValidateSet("migrate", "seed", "prepare")]
    [string]$Command
)

$ErrorActionPreference = "Stop"
$ProjectRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$EnvPath = Join-Path $ProjectRoot ".env"

if (Test-Path $EnvPath) {
    Get-Content $EnvPath | ForEach-Object {
        $line = $_.Trim()
        if ($line -eq "" -or $line.StartsWith("#") -or -not $line.Contains("=")) { return }
        $name, $value = $line.Split("=", 2)
        if ($name.Trim() -eq "TEST_DATABASE_URL" -and -not $env:TEST_DATABASE_URL) {
            $env:TEST_DATABASE_URL = $value.Trim().Trim('"').Trim("'")
        }
    }
}

if (-not $env:TEST_DATABASE_URL) {
    $env:TEST_DATABASE_URL = "postgres://miskatonic_user:miskatonic_password@localhost:5433/miskatonic_lab_test?sslmode=disable"
}

$parsed = [System.Uri]$env:TEST_DATABASE_URL
$hostName = $parsed.Host.ToLowerInvariant()
$databaseName = $parsed.AbsolutePath.TrimStart('/')
if (($hostName -ne "localhost" -and $hostName -ne "127.0.0.1" -and $hostName -ne "::1") -or -not $databaseName.EndsWith("_test")) {
    throw "TEST_DATABASE_URL must target a local database whose name ends with _test"
}

$goBin = Join-Path (go env GOPATH) "bin"
$migrate = Join-Path $goBin "migrate.exe"
if (-not (Test-Path $migrate)) {
    throw "migrate command not found; run: go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3"
}

function Invoke-TestMigrations {
    & $migrate -path "migrations" -database $env:TEST_DATABASE_URL up
    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
}

Push-Location $ProjectRoot
try {
    switch ($Command) {
        "migrate" { Invoke-TestMigrations }
        "seed" { go run ./cmd/seed-testdb }
        "prepare" {
            Invoke-TestMigrations
            go run ./cmd/seed-testdb
        }
    }
    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
} finally {
    Pop-Location
}
