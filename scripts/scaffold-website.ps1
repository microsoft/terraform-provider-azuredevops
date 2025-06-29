[CmdletBinding()]
param (
    [Parameter()]
    [string]
    [ValidateNotNullOrEmpty()]
    $BrandName,

    [Parameter()]
    [string]
    [ValidateNotNullOrEmpty()]
    $ResourceName,

    [Parameter()]
    [string]
    [ValidateNotNullOrEmpty()]
    [ValidateSet('resource', 'data')]
    $ResourceType,

    [Parameter()]
    [string]
    [ValidateNotNullOrEmpty()]
    $ResourceId
)

$script:PSDefaultParameterValues = @{
    '*:Confirm'           = $false
    '*:ErrorAction'       = 'Stop'
}


Write-Host "==> Scaffolding Documentation..."

& go run azuredevops/internal/website-scaffold/main.go `
    -name $ResourceName `
    -brand-name $BrandName `
    -type $ResourceType `
    -resource-id $ResourceId `
    -website-path "./website/"

Write-Host "==> Done."
