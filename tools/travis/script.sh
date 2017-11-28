#!/usr/bin/env bash

set -e

cd $TRAVIS_BUILD_DIR
./tools/travis/scancode.sh
make lint
make build
export PATH=$PATH:$TRAVIS_BUILD_DIR

HOMEDIR="$(dirname "$TRAVIS_BUILD_DIR")"
cd $HOMEDIR

# OpenWhisk clone to fixed directory location
git clone https://github.com/apache/incubator-openwhisk.git openwhisk
cd openwhisk
git reset --hard 0f679adcf8cf1f29fde8f99c16f66a98e3ec021f

# Build script for Travis-CI.
WHISKDIR="$HOMEDIR/openwhisk"

cd $WHISKDIR
./tools/travis/setup.sh

ANSIBLE_CMD="ansible-playbook -i environments/local -e docker_image_prefix=testing"
TERM=dumb ./gradlew distDocker -PdockerImagePrefix=testing

cd $WHISKDIR/ansible
$ANSIBLE_CMD setup.yml
$ANSIBLE_CMD prereq.yml
$ANSIBLE_CMD couchdb.yml
$ANSIBLE_CMD initdb.yml

$ANSIBLE_CMD wipe.yml
$ANSIBLE_CMD openwhisk.yml -e '{"openwhisk_cli":{"installation_mode":"remote","remote":{"name":"OpenWhisk_CLI","dest_name":"OpenWhisk_CLI","location":"https://github.com/apache/incubator-openwhisk-cli/releases/download/latest"}}}'


export OPENWHISK_HOME="$(dirname "$TRAVIS_BUILD_DIR")/openwhisk"

cd $TRAVIS_BUILD_DIR
make integration_test
