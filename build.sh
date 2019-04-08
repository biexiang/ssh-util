#!/usr/bin/env bash

TOOL="ssh-util"
VERSION="0.01"
Date=`date +%F%T%z`

function build() {
    os=$1
    arch=$2
    alias=$3
    pkg="${TOOL}_${os}_${arch}_${alias}_${VERSION}"

    echo "build ${pkg} ... "
    mkdir -p "./pack/${pkg}"
    CGO_ENABLED=0 GOOS=${os} GOARCH=${arch} go build -o "./pack/${pkg}/${TOOL}" -ldflags "-X main.Version=${VERSION} -X main.Date=${Date}" ./src/main.go
    cp ./confg.yaml "./pack/${pkg}/config.yaml"
}

build darwin amd64 macOS
build linux amd64 linux
build windows amd64 windows