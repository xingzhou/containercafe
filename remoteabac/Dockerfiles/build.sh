#!/bin/bash

IMAGE_NAME=${1:-radiant/remoteabac}
BUILD_ID=${2:-0}
BUILD_NUMBER=${3:-0}
GIT_TAG=${4:-''}

PROJECT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )"/.. && pwd )"

# build the binaries
docker build -t radiant/remoteabac_builder -f $PROJECT_DIR/Dockerfiles/Dockerfile.build $PROJECT_DIR

# copy out the binaries
docker run --rm -v $PROJECT_DIR/Dockerfiles:/tmp:rw radiant/remoteabac_builder cp /work/bin/remoteabac /tmp
docker run --rm -v $PROJECT_DIR/Dockerfiles:/tmp:rw radiant/remoteabac_builder cp /work/bin/ruser /tmp

# build the image
docker build \
--build-arg git_commit_id=${GITCOMMIT} \
--build-arg build_id=${BUILD_ID} \
--build-arg build_number=${BUILD_NUMBER} \
--build-arg build_date=${DATE} \
--build-arg git_tag=${GIT_TAG} \
--build-arg git_remote_url=${GITREMOTE} \
-f $PROJECT_DIR/Dockerfiles/Dockerfile \
-t $IMAGE_NAME --no-cache $PROJECT_DIR/Dockerfiles
