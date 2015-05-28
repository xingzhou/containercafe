#!/bin/bash
set -v
export GOPATH=~/containers-hijack-proxy
go build -o bin/hijack src/hijack.go
ls -l bin