# Adds a `tt` alias pointing at the built binary to your PowerShell profile (Windows).
$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent $PSScriptRoot
$bin = Join-Path $repoRoot 'tt.exe'

# Build the binary if it isn't there yet.
if (-not (Test-Path $bin)) {
    Write-Host 'Building tt...'
    Push-Location $repoRoot
    try {
        go build -o tt.exe main.go
    } finally {
        Pop-Location
    }
}

$profilePath = $PROFILE
$profileDir = Split-Path -Parent $profilePath
if (-not (Test-Path $profileDir)) {
    New-Item -ItemType Directory -Path $profileDir -Force | Out-Null
}
if (-not (Test-Path $profilePath)) {
    New-Item -ItemType File -Path $profilePath -Force | Out-Null
}

$aliasLine = "Set-Alias tt '$bin'"

# Idempotent: refresh an existing tt alias instead of appending a duplicate.
if (Select-String -Path $profilePath -Pattern '^\s*Set-Alias\s+tt\b' -Quiet) {
    $kept = Get-Content $profilePath | Where-Object { $_ -notmatch '^\s*Set-Alias\s+tt\b' }
    Set-Content -Path $profilePath -Value $kept
    Add-Content -Path $profilePath -Value $aliasLine
    Write-Host "Updated existing tt alias in $profilePath"
} else {
    Add-Content -Path $profilePath -Value "`n# Added by ttool install script`n$aliasLine"
    Write-Host "Added tt alias to $profilePath"
}

Write-Host "Done. Restart PowerShell or run: . `$PROFILE"
