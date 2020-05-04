[CmdletBinding()]
param (
    [switch]
    $Fix
)

$script:PSDefaultParameterValues = @{
    '*:Confirm'           = $false
    '*:ErrorAction'       = 'Stop'
}

. (Join-Path -Path $PSScriptRoot -ChildPath 'commons.ps1' -Resolve)

if ($Fix) {
    # Check gofmt
    echo "==> Fixing gofmt deviations..."

    # This filter should match the search filter in ../GNUMakefile
    $null = Get-ChildItem -Path $SOURCE_DIR -Recurse -Filter '*.go' `
    | Select-Object -ExpandProperty FullName `
    | Select-String -NotMatch -SimpleMatch vendor `
    | ForEach-Object -Process { gofmt.exe -s -w $_ }
}
else {
    # Check gofmt
    echo "==> Checking that code complies with gofmt requirements..."

    # This filter should match the search filter in ../GNUMakefile
    $gofmt_files= Get-ChildItem -Path $SOURCE_DIR -Recurse -Filter '*.go' `
    | Select-Object -ExpandProperty FullName `
    | Select-String -NotMatch -SimpleMatch vendor `
    | ForEach-Object -Process { gofmt.exe -s -l $_ }

    if ($gofmt_files) {
        echo 'gofmt needs running on the following files:'
        echo "${gofmt_files}"
        echo "You can use this command and pass the -Fix parameter to reformat code."
        exit 1
    }
}
