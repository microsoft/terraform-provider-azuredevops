[CmdletBinding()]
param (
    [Parameter()]
    [string]
    $PluginsDirectory
)

$script:PSDefaultParameterValues = @{
    '*:Confirm'           = $false
    '*:ErrorAction'       = 'Stop'
}

function Install-Provider {
    [CmdletBinding()]
    param (
        [Parameter(Mandatory, Position=0, ValueFromPipeline, ValueFromPipelineByPropertyName)]
        [ValidateNotNullOrEmpty()]
        [string]
        $LiteralPath
    )

    PROCESS {
        if (Test-Path -Path $LiteralPath) {
            Write-Verbose -Message "Terraform Plugins directory [$LiteralPath] already exists"
        }
        else {
            Write-Verbose -Message "Creating Terraform Plugins directory [$LiteralPath]"
            $null = New-Item -Path $LiteralPath -ItemType Directory
        }

        Write-Host "Installing provider to $LiteralPath"
        Copy-Item -Path (Join-Path -Path $BUILD_DIR -ChildPath '*') -Destination $LiteralPath -Force
    }
}

if (-not $PluginsDirectory) {
    . (Join-Path -Path $PSScriptRoot -ChildPath 'commons.ps1' -Resolve)

    # https://www.terraform.io/docs/plugins/basics.html
    # https://www.terraform.io/docs/extend/how-terraform-works.html#discovery
    if ($env:OS -like '*Windows*') {
        $PluginsBase = $env:APPDATA
    }
    else {
        $PluginsBase = $HOME
    }
    $PluginsDirectory = [System.IO.Path]::Combine($PluginsBase, '.terraform.d', 'plugins')
    Install-Provider -LiteralPath $PluginsDirectory

    ## Terraform >= v0.13 requires different layout
    $ProviderName=Get-Content -LiteralPath "$PROVIDER_NAME_FILE"
    $ProviderVersion=Get-Content -LiteralPath "$PROVIDER_VERSION_FILE"
    $ProviderRegistry='registry.terraform.io'
    $ProviderOrganization='terraform-providers'

    $PluginsDirectory = [System.IO.Path]::Combine($PluginsDirectory, $ProviderRegistry, $ProviderOrganization, $ProviderName, $ProviderVersion, "${OS}_${PROC}")
    Install-Provider -LiteralPath $PluginsDirectory
}
else {
    Install-Provider -LiteralPath $PluginsDirectory
}