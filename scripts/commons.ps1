[CmdletBinding()]
param (
)

$SCRIPTS_DIR=$PSScriptRoot
$BUILD_DIR=[IO.Path]::GetFullPath([IO.Path]::Combine($SCRIPTS_DIR, '..', 'bin'))
$SOURCE_DIR=[IO.Path]::GetFullPath([IO.Path]::Combine($SCRIPTS_DIR, '..'))
$PROVIDER_NAME_FILE=[IO.Path]::GetFullPath([IO.Path]::Combine($SCRIPTS_DIR, '..', 'PROVIDER_NAME.txt'))
$PROVIDER_VERSION_FILE=[IO.Path]::GetFullPath([IO.Path]::Combine($SCRIPTS_DIR, '..', 'PROVIDER_VERSION.txt'))

if ($PSVersionTable.PSVersion.Major -lt 6) {
    $OS = 'windows'
}
else {
    $OS=$PSVersionTable.OS.Split(' ')[0].ToLower()
}
if ("$OS" -eq 'linux') {
  $match = $(lscpu) | Select-String 'Architecture:\s*(?<PROC>.+)$'
  if ($match) {
    $PROC = switch -Wildcard ($match.Matches[0].Groups['PROC'].Value) {
      'x86_64' {
        'amd64'
      }
      '*arm*' {
        'arm'
      }
      '*aarch64*' {
        'arm'
      }
      default {
        '386'
      }
    }
  }
  if ([string]::IsNullOrWhiteSpace($PROC)) {
    $PROC=if (Get-Content /proc/cpuinfo | Select-String 'model name\s*:\s*ARM' | Select-Object -First 1) { 'arm' }
  }
  if ([string]::IsNullOrWhiteSpace($PROC)) {
    $PROC=if (Get-Content /proc/cpuinfo | Select-String 'flags\s*:\s*.* lm ' | Select-Object -First 1) { 'amd64' } else { '386' }
  }
}
else {
  $PROC="amd64"
}
if ($PROC -like "*arm*") {
    # terraform downloads use "arm" not full arm type
    $PROC = 'arm'
}
