#!/usr/bin/env bash

declare -a os_list=("linux" "darwin" "windows")
arc=amd64
file_name=${1:-"wskdeploy"}

for os in "${os_list[@]}"
do
    build_file_name=$file_name
    if [ "$os" == "windows" ]; then
        build_file_name="$build_file_name.exe"
    fi
    cd $TRAVIS_BUILD_DIR
    GOOS=$os GOARCH=$arc go build -o build/$os/$build_file_name
    cd build/$os
    zip -r "$TRAVIS_BUILD_DIR/$build_file_name-$TRAVIS_TAG-$os-$arc.zip" $build_file_name
done