#!/usr/bin/env bash
#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

declare -a os_list=("linux" "darwin" "windows")
declare -a arc_list=("amd64" "386")
build_file_name=${1:-"wskdeploy"}
build_version=${2:-"$TRAVIS_TAG"}

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
        GOOS=$os GOARCH=$arc go build -ldflags "-X main.Version=$build_version" -o build/$os/$wskdeploy
        cd build/$os
        if [[ "$os" == "linux" ]]; then
            tar -czvf "$TRAVIS_BUILD_DIR/$build_file_name-$TRAVIS_TAG-$os_name-$arc.tgz" $wskdeploy
        else
            zip -r "$TRAVIS_BUILD_DIR/$build_file_name-$TRAVIS_TAG-$os_name-$arc.zip" $wskdeploy
        fi
    done
done
