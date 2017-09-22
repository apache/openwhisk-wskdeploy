#!/usr/bin/env bash

declare -a os_list=("linux" "darwin" "windows")
declare -a arc_list=("amd64" "386")
build_file_name=${1:-"wskdeploy"}

for os in "${os_list[@]}"
do
    for arc in "${arc_list[@]}"
    do
        wskdeploy=$build_file_name
        os_name=$os
        if [ "$os" == "windows" ]; then
            wskdeploy="$wskdeploy.exe"
        fi
        if [ "$os" == "darwin" ]; then
            os_name="mac"
        fi
        cd $TRAVIS_BUILD_DIR
        GOOS=$os GOARCH=$arc go build -o build/$os/$wskdeploy
        cd build/$os
        if [[ "$os" == "linux" ]]; then
            tar -czvf "$TRAVIS_BUILD_DIR/$build_file_name-$TRAVIS_TAG-$os_name-$arc.tgz" $wskdeploy
        else
            zip -r "$TRAVIS_BUILD_DIR/$build_file_name-$TRAVIS_TAG-$os_name-$arc.zip" $wskdeploy
        fi
    done
done
