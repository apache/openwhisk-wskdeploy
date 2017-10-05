#!/bin/bash
set -eu

SCRIPTDIR=$(cd $(dirname "$0") && pwd)

dockerhub_image_prefix="$1"
dockerhub_image_name="$2"
dockerhub_image_tag="$3"
dockerhub_image="${dockerhub_image_prefix}/${dockerhub_image_name}:${dockerhub_image_tag}"

docker login -u "${DOCKER_USER}" -p "${DOCKER_PASSWORD}"

cp $TRAVIS_BUILD_DIR/build/linux/wskdeploy $SCRIPTDIR/

echo docker build $SCRIPTDIR --tag ${dockerhub_image}
docker build $SCRIPTDIR --tag ${dockerhub_image}

echo docker push ${dockerhub_image}
docker push ${dockerhub_image}
