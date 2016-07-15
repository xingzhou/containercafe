#!/bin/bash

# Make sure i) can sudo as root and ii) Docker is running

IMAGE_NAME=haih/remoteabac

# For cross compiling
env GOOS=linux GOARCH=amd64 godep go install ../cmd/remoteabac/remoteabac.go
env GOOS=linux GOARCH=amd64 godep go install ../cmd/ruser/ruser.go

cp $GOBIN/remoteabac .
cp $GOBIN/ruser .
touch empty

docker build -f haih.Dockerfile -t $IMAGE_NAME --no-cache .
