#!/bin/bash

__build_tbg() {
    arch=$1
    mkdir -p ./bin
    version=$(git describe --tags --dirty --always)
    
    if [[ "$arch" != "arm" && "$arch" != "arm64" && "$arch" != "amd64" && "$arch" != "386" && "$arch" != "dev" ]]; then
        echo "Invalid GOARCH value: \"$arch\"."
        echo "Please provide one of: arm, arm64, amd64, or 386."
        exit 1
    fi
    
    go mod tidy
    if [[ "$arch" == "dev" ]]; then
        GOOS=windows go build -ldflags "-X main.TbgVersion=$version-dev"
        mv ./tbg.exe ./bin/tbg.exe
    else
        GOOS=windows GOARCH=$arch go build -ldflags "-s -w -X main.TbgVersion=$version"
        mv ./tbg.exe ./bin/tbg-$version.windows.$arch.exe
        sha256sum ./bin/tbg-$version.windows.$arch.exe | cut -d ' ' -f 1 > ./bin/tbg-$version.windows.$arch.exe.sha256
    fi
}

arg=$1
if [[ "$arg" == "release" ]]; then
    __build_tbg "arm"
    __build_tbg "arm64"
    __build_tbg "amd64"
    __build_tbg "386"
elif [[ -z "$arg" ]]; then
    __build_tbg "dev"
else
    echo "Invalid argument: \"$arg\". Please use either \"release\" or leave the argument empty."
    exit 1
fi
