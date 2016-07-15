#!/bin/bash

LOCAL_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source $LOCAL_DIR/env.sh

PROJECT_DIR=$GOPATH/src/github.ibm.com/alchemy-containers/remoteabac

docker run --rm -it -v /var/run/docker.sock:/var/run/docker.sock \
  -v /usr/local/bin/docker:/bin/docker -v $(pwd)/../:$PROJECT_DIR --workdir $PROJECT_DIR \
  ubuntu:trusty Dockerfiles/build.sh $1 $2 $3 $4

