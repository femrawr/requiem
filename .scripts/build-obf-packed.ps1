$origin = Get-Location

try {
    if ((Get-Location).Path -like "*.scripts*") {
        Set-Location ..\
    }

    Set-Location .\requiem\

    garble -tiny -seed=random build -trimpath -buildvcs=false -ldflags="-s -w -H windowsgui -buildid=" -o .\requiem.exe

    upx -3 -o .\requiem.exe.upx .\requiem.exe

    Remove-Item -Path .\requiem.exe

    Rename-Item -Path .\requiem.exe.upx -NewName requiem.exe
}
finally {
    Set-Location $origin
}
