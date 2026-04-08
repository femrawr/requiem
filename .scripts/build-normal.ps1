$origin = Get-Location

try {
    if ((Get-Location).Path -like "*.scripts*") {
        Set-Location ..\
    }

    Set-Location .\requiem\

    go build -trimpath -buildvcs=false -ldflags="-s -w -H windowsgui -buildid=" -o .\requiem.exe
}
finally {
    Set-Location $origin
}
