#!/bin/bash
set -v
cp -r src dockerize
# copy TLS scripts to dockerize dir.
cp admin-certs/* dockerize
cp creds.json dockerize/
cp make_TLS_certs.sh dockerize/
cp mk_user_cert.sh dockerize/

cd dockerize
docker build -t proxy .
