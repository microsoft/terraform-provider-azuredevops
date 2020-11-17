[CmdletBinding()]
param (
)

$script:PSDefaultParameterValues = @{
    '*:Confirm'           = $false
    '*:ErrorAction'       = 'Stop'
}

. (Join-Path -Path $PSScriptRoot -ChildPath 'commons.ps1' -Resolve)


echo "[INFO] Linting Go Files... If this fails, run 'golint ./... | grep -v 'vendor' ' to see errors"

Push-Location -Path $SOURCE_DIR
try {
    go get -u golang.org/x/lint/golint 2>$null
    if ($LASTEXITCODE) {
        throw "Failed to install or update golint"
    }

    go list ./... `
    | Select-String -NotMatch -SimpleMatch 'vendor' `
    | ForEach-Object -Process { 
        $package = $_
        golint.exe -set_exit_status $package
        if ($LASTEXITCODE) {
            Write-Error -Message "Linting failed for package: $package"
        }
    }
}
finally {
    Pop-Location
}
