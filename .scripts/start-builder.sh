#!/bin/bash

origin=$(pwd)

cleanup() {
    cd "$origin" 2>/dev/null
}

trap cleanup EXIT INT TERM

if [[ "$(pwd)" == *"scripts"* ]]; then
    cd .. || exit 1
fi

cd "./builder/server" || exit 1

go build -trimpath -buildvcs=false -ldflags="-s -w -buildid=" -o ./builder.exe

./builder.exe
