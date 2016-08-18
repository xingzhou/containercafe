#!/bin/bash
set -v

export GOPATH="$(pwd):$GOPATH"

go build -o bin/hijack src/hijack/hijack.go

ls -l bin
