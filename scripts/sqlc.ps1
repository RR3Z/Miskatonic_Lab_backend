param(
    [Parameter(Position = 0)]
    [ValidateSet("generate", "version")]
    [string]$Command = "generate"
)

$ErrorActionPreference = "Stop"

$ProjectRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$SqlcVersion = "v1.31.1"

Push-Location $ProjectRoot
try {
    if (Get-Command sqlc -ErrorAction SilentlyContinue) {
        & sqlc $Command
    } else {
        & go run "github.com/sqlc-dev/sqlc/cmd/sqlc@$SqlcVersion" $Command
    }

    exit $LASTEXITCODE
} finally {
    Pop-Location
}
