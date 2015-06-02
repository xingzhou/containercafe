#!/bin/bash
set -v
#export GOPATH=/c/src/containers-hijack-proxy
export GOPATH=~/containers-hijack-proxy
go build -o bin/hijack src/hijack/hijack.go
ls -l bin
