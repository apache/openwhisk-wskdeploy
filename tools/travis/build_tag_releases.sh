#!/usr/bin/env bash

declare -a os_list=("linux" "darwin" "windows")
arc=amd64
build_file_name=${1:-"wskdeploy"}

for os in "${os_list[@]}"
do
    wskdeploy=$build_file_name
    if [ "$os" == "windows" ]; then
        wskdeploy="$wskdeploy.exe"
    fi
    cd $TRAVIS_BUILD_DIR
    GOOS=$os GOARCH=$arc go build -o build/$os/$wskdeploy
    cd build/$os
    zip -r "$TRAVIS_BUILD_DIR/$build_file_name-$TRAVIS_TAG-$os-$arc.zip" $wskdeploy
done