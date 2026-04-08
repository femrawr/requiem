$origin = Get-Location

try {
    if ((Get-Location).Path -like "*.scripts*") {
        Set-Location ..\
    }

    Set-Location .\builder\
    Set-Location .\server\

    go build -trimpath -buildvcs=false -ldflags="-s -w -buildid=" -o .\builder.exe

    .\builder.exe
}
finally {
    Set-Location $origin
}
