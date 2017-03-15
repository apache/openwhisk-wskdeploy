#!/usr/bin/env bash

declare -a os_list=("linux" "darwin" "windows")
arc=amd64
build_file_name=${1:-"$wskdeploy"}

for os in "${os_list[@]}"
do
    GOOS=$os GOARCH=$arc go build -o $build_file_name-$os-$arc
    zip -r "$build_file_name-$TRAVIS_TAG-$os-$arc.zip" $build_file_name-$os-$arc
done