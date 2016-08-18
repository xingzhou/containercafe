#!/bin/bash
# set -v
# copy src into dockerize directory so it can be build by Dockerfile
cp -r src dockerize

# copy creds file and scripts to dockerize dir.
touch dockerize/creds.json
cp make_TLS_certs.sh dockerize/
cp mk_user_cert.sh dockerize/
cp mk_kubeconfig.sh dockerize/

cd dockerize
docker build -t api-proxy .
