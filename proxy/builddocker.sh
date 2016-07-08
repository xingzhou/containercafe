#!/bin/bash
set -v
cp -r src dockerize
cd dockerize
docker build -t hijack .
