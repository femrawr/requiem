$origin = gl

try {
    if ((gl).Path -like "*.scripts*") {
        cd ..\
    }

    cd .\builder\
    cd .\server\

    go build -trimpath -buildvcs=false -ldflags="-s -w -buildid=" -o .\builder.exe

    .\builder.exe
}
finally {
    cd $origin
}
