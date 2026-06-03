param(
    [Parameter(Position = 0)]
    [ValidateSet("generate", "version")]
    [string]$Command = "generate"
)

$ErrorActionPreference = "Stop"

$ProjectRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$SqlcVersion = "v1.31.1"
$GeneratedGoOut = Join-Path $ProjectRoot "pkg/repository/db"

Push-Location $ProjectRoot
try {
    if ($Command -eq "generate") {
        $outputPath = [System.IO.Path]::GetFullPath($GeneratedGoOut)
        $projectRootPath = [System.IO.Path]::GetFullPath($ProjectRoot)

        if (-not $outputPath.StartsWith($projectRootPath, [System.StringComparison]::OrdinalIgnoreCase)) {
            throw "Refusing to remove generated output outside project root: $outputPath"
        }

        if (Test-Path -LiteralPath $GeneratedGoOut) {
            Remove-Item -LiteralPath $GeneratedGoOut -Recurse -Force
        }

        New-Item -ItemType Directory -Path $GeneratedGoOut -Force | Out-Null
    }

    if (Get-Command sqlc -ErrorAction SilentlyContinue) {
        & sqlc $Command
    } else {
        & go run "github.com/sqlc-dev/sqlc/cmd/sqlc@$SqlcVersion" $Command
    }

    exit $LASTEXITCODE
} finally {
    Pop-Location
}
