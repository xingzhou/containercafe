#!/bin/bash
set -v
cp -r src dockerize
cp -f ../ansible/certs/dev-vbox-radiant01/ca* admin-certs
cp -f ../ansible/certs/dev-vbox-radiant01/admin-key.pem admin-certs/kadmin.key
cp -f ../ansible/certs/dev-vbox-radiant01/admin.pem admin-certs/kadmin.pem
./gen_server_certs.sh

# copy TLS scripts to dockerize dir.
cp admin-certs/* dockerize
cp creds.json dockerize/
cp make_TLS_certs.sh dockerize/
cp mk_user_cert.sh dockerize/
cp mk_kubeconfig.sh dockerize/

cd dockerize
docker build -t proxy .
