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

. (Join-Path -Path $PSScriptRoot -ChildPath 'commons.ps1' -Resolve)

if (-not $PluginsDirectory) {
    # https://www.terraform.io/docs/plugins/basics.html
    # https://www.terraform.io/docs/extend/how-terraform-works.html#discovery
    if ($env:OS -like '*Windows*') {
        $PluginsDirectory = [System.IO.Path]::Combine($env:APPDATA, 'terraform.d', 'plugins')
    }
    else {
        $PluginsDirectory = [System.IO.Path]::Combine($HOME, '.terraform.d', 'plugins')
    }
}

if (Test-Path -Path $PluginsDirectory) {
    Write-Verbose -Message "Terraform Plugins directory [$PluginsDirectory] already exists"
}
else {
    Write-Verbose -Message "Creating Terraform Plugins directory [$PluginsDirectory]"
    $null = New-Item -Path $PluginsDirectory -ItemType Directory
}

Write-Host "Installing provider to $PluginsDirectory"
Copy-Item -Path (Join-Path -Path $BUILD_DIR -ChildPath '*') -Destination $PluginsDirectory -Force
