[CmdletBinding()]
param (
    [Parameter()]
    [switch]
    $SkipTests,

    [Parameter()]
    [switch]
    $Install,

    [Parameter()]
    [switch]
    $DebugBuild,

    [Parameter()]
    [ValidateSet('', 'readonly', 'vendor')]
    [string]
    $GoMod = 'vendor'
)

$script:PSDefaultParameterValues = @{
    '*:Confirm'     = $false
    '*:ErrorAction' = 'Stop'
}

. (Join-Path -Path $PSScriptRoot -ChildPath 'commons.ps1' -Resolve)

function clean() {
    Write-Host "Cleaning $BUILD_DIR"
    if (Test-Path -Path $BUILD_DIR) {
        Remove-Item -Recurse -Force -Path $BUILD_DIR
    }
    $null = New-Item -ItemType Container -Path $BUILD_DIR
}

function compile() {
    $NAME = Get-Content -Raw -Path $PROVIDER_NAME_FILE
    $VERSION = Get-Content -Raw -Path $PROVIDER_VERSION_FILE

    $BUILD_ARTIFACT = "terraform-provider-${NAME}_v${VERSION}"
    if ($env:OS -like '*Windows*') {
        $BUILD_ARTIFACT += '.exe'
    }
    Write-Host "Attempting to build $BUILD_ARTIFACT"
    Push-Location -Path $SOURCE_DIR
    try {
        $env:GO111MODULE = 'on'

        go mod download
        if ($LASTEXITCODE) {
            throw "Failed to download modules"
        }

        $argv = @(
            'build',
            "-mod=$(if ('' -ne $GoMod) { $GoMod } else { $null })",
            '-o',
            "$BUILD_DIR/$BUILD_ARTIFACT"
        )
        if ($DebugBuild) {
            $argv += @( '-gcflags="all=-N -l"' )
        }
        go @argv
        if ($LASTEXITCODE) {
            throw "Build failed"
        }
    }
    finally {
        'GO111MODULE' `
        | ForEach-Object -Process { Remove-Item -Path "Env:$_" }
        Pop-Location
    }
}

function clean_and_build() {
    clean
    compile
    if (-not $SkipTests) {
        & (Join-Path -Path $PSScriptRoot -ChildPath 'unittest.ps1' -Resolve) -GoMod $GoMod
    }
    Write-Host "Build finished successfully"
    if ($Install) {
        & (Join-Path -Path $PSScriptRoot -ChildPath 'local-install.ps1' -Resolve)
    }
}

clean_and_build
