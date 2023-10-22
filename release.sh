#!/bin/bash
DIR=$(dirname $0)
EXECUTABLE="ScreenToArtnet"

build() {
    local os=$1
    local arch=$2
    local version=$3
    local out_dir="$DIR/bin/$os"
    local ext=""
    if [ $os = "windows" ]
    then
        ext=".exe"
    fi
    if [ -n "$version" ]
    then
        version="-$version"
    fi
    export GOOS=$os
    export GOARCH=$arch
    mkdir -p $out_dir
    go build -o "$out_dir/$EXECUTABLE$ext"
    tar -zcvf $EXECUTABLE-$os$version.tar.gz --directory $(realpath $out_dir) . ../../README.md ../../LICENSE ../../config.json
}

version=$1

build "windows" "amd64" $version
build "linux" "amd64" $version
build "darwin" "amd64" $version