#!/bin/bash
set -v
# copy src into dockerize directory so it can be build by Dockerfile
cp -r src dockerize

# manage certificates - copy them to admin-certs
CERTS=../ansible/certs/dev-vbox-radiant01
ACERTS=admin-certs
mkdir -p $ACERTS
if [ ! -d "$CERTS" ]; then
  echo "missing $CERTS directory"
  exit 99
fi

cp -f $CERTS/ca* $ACERTS
cp -f $CERTS/admin-key.pem $ACERTS/kadmin.key
cp -f $CERTS/admin.pem $ACERTS/kadmin.pem

./gen_server_certs.sh

# copy TLS certs, keys and scripts to dockerize dir.
cp $ACERTS/* dockerize
cp creds.json dockerize/
cp make_TLS_certs.sh dockerize/
cp mk_user_cert.sh dockerize/
cp mk_kubeconfig.sh dockerize/

cd dockerize
docker build -t proxy .
