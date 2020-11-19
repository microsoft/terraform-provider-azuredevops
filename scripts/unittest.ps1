[CmdletBinding()]
param (
    [Parameter()]
    [ValidateNotNullOrEmpty()]
    [string]
    $TestFilter,

    [Parameter()]
    [ValidateNotNullOrEmpty()]
    [string[]]
    $Tag = 'all',

    [Parameter()]
    [ValidateSet('', 'readonly', 'vendor')]
    [string]
    $GoMod = 'vendor'
)

$script:PSDefaultParameterValues = @{
    '*:Confirm'           = $false
    '*:ErrorAction'       = 'Stop'
}

. (Join-Path -Path $PSScriptRoot -ChildPath 'commons.ps1' -Resolve)

Write-Host "Executing unit tests"
Push-Location -Path $SOURCE_DIR
try {
    Remove-Item -ErrorAction:Ignore 'ENV:TF_ACC'
    $env:TF_SCHEMA_PANIC_ON_ERROR=1
    $env:GO111MODULE='on'

    $argv = @(
        'test',
        "-mod=$(if ('' -ne $GoMod) { $GoMod } else { $null })",
        '-v'
    )
    if ($TestFilter) {
        $argv += @('-run', $TestFilter)
    }
    if ($Tag -and 0 -lt $Tag.Length) {
        $argv += @('-tags', [string]::Join(' ', $Tag))
    }
    go @argv ./...
    if ($LASTEXITCODE) {
        throw "Build finished in error due to failed tests"
    }
}
finally {
    'TF_SCHEMA_PANIC_ON_ERROR', 'GO111MODULE' `
    | ForEach-Object -Process {Remove-Item -Path "Env:$_" }
    Pop-Location
}
