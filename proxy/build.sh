#!/bin/bash
set -v
PWD=`pwd`
export GOPATH=$PWD:$GOPATH
#export GOPATH=~/containers-hijack-proxy
go build -o bin/hijack src/hijack/hijack.go
ls -l bin
