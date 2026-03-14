$origin = gl

try {
    if ((gl).Path -like "*.scripts*") {
        cd ..\
    }

    cd .\requiem\

    go build -trimpath -buildvcs=false -ldflags="-s -w -H windowsgui -buildid=" -o .\requiem.exe

    start .\requiem.exe -verb RunAs
}
finally {
    cd $origin
}
