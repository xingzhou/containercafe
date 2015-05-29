#!/bin/bash
set -v
cd -r src dockerize
cd dockerize
docker build -t hijack .
