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

set -e

SCRIPTDIR=$(cd $(dirname "$0") && pwd)
ROOTDIR="$SCRIPTDIR/../.."

cd $TRAVIS_BUILD_DIR
./tools/travis/scancode.sh
make lint
make test
BUILD_VERSION="latest"
if [ ! -z "$TRAVIS_TAG" ] ; then
    BUILD_VERSION=$TRAVIS_TAG
fi
make build VERSION=$BUILD_VERSION
export PATH=$PATH:$TRAVIS_BUILD_DIR

HOMEDIR="$(dirname "$TRAVIS_BUILD_DIR")"
cd $HOMEDIR

# OpenWhisk clone to fixed directory location
git clone --depth 3 https://github.com/apache/openwhisk.git openwhisk

# Build script for Travis-CI.
WHISKDIR="$HOMEDIR/openwhisk"

cd $WHISKDIR
./tools/travis/setup.sh

ANSIBLE_CMD="ansible-playbook -i ${ROOTDIR}/ansible/environments/local -e docker_image_prefix=openwhisk -e docker_image_tag=nightly"

cd $WHISKDIR/ansible
$ANSIBLE_CMD setup.yml
$ANSIBLE_CMD prereq.yml
$ANSIBLE_CMD couchdb.yml
$ANSIBLE_CMD initdb.yml
$ANSIBLE_CMD wipe.yml
$ANSIBLE_CMD openwhisk.yml -e '{"openwhisk_cli":{"installation_mode":"remote","remote":{"name":"OpenWhisk_CLI","dest_name":"OpenWhisk_CLI","location":"https://github.com/apache/openwhisk-cli/releases/download/latest"}}}'
$ANSIBLE_CMD properties.yml
$ANSIBLE_CMD apigateway.yml
$ANSIBLE_CMD routemgmt.yml

export OPENWHISK_HOME="$(dirname "$TRAVIS_BUILD_DIR")/openwhisk"

cd $TRAVIS_BUILD_DIR
make integration_test
