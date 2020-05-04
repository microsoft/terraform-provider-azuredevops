[CmdletBinding()]
param (
)

$SCRIPTS_DIR=$PSScriptRoot
$BUILD_DIR=[IO.Path]::GetFullPath([IO.Path]::Combine($SCRIPTS_DIR, '..', 'bin'))
$SOURCE_DIR=[IO.Path]::GetFullPath([IO.Path]::Combine($SCRIPTS_DIR, '..'))
$PROVIDER_NAME_FILE=[IO.Path]::GetFullPath([IO.Path]::Combine($SCRIPTS_DIR, '..', 'PROVIDER_NAME.txt'))
$PROVIDER_VERSION_FILE=[IO.Path]::GetFullPath([IO.Path]::Combine($SCRIPTS_DIR, '..', 'PROVIDER_VERSION.txt'))
