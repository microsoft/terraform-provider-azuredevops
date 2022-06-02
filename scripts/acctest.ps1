[CmdletBinding()]
param (
    [Parameter()]
    [ValidateNotNullOrEmpty()]
    [string]
    $TestFilter = '^TestAcc',

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

Write-Host "Executing acceptance tests"
Push-Location -Path $SOURCE_DIR
try {
    # This is similar to the unit test command aside from the following:
    #   - TF_ACC=1 is a flag that will enable the acceptance tests. This flag is
    #     documented here:
    #       https://www.terraform.io/docs/extend/testing/acceptance-tests/index.html#running-acceptance-tests
    #
    #   - A `-run` parameter is used to target *only* tests starting with `TestAcc`. This prefix is
    #     recommended by Hashicorp and is documented here:
    #       https://www.terraform.io/docs/extend/testing/acceptance-tests/index.html#test-files
    #
    # Using build tags as test filter: https://stackoverflow.com/a/24036237
    $env:TF_ACC=1
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
    'TF_ACC', 'TF_SCHEMA_PANIC_ON_ERROR', 'GO111MODULE' `
    | ForEach-Object -Process {Remove-Item -Path "Env:$_" }
    Pop-Location
}
