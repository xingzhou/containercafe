#!/bin/bash
#set -v

# Help menu
helpme()
{
	cat <<HELPMEHELPME

Syntax: ${0} <env_name> <-d>
Where:
	env_name - name of the environment, e.g: dev-vbox
	-d - run api-proxy container in the background (optional)
	
HELPMEHELPME
}

if [[ "$1" == "" || "$1" == "-d" ]] ; then
	echo 'env_name must be set'
	helpme
	exit 1
fi
env_name=$1

# manage certificates - copy them to admin-certs
CERT_MASTER=../ansible/certs/dev-vbox-radiant01
CERTS=~/.openradiant/envs/$env_name
ACERTS=$CERTS/admin-certs
if [ ! -d "$CERT_MASTER" ]; then
  echo "missing $CERT_MASTER directory. Execute ansible scripts first"
  exit 99
fi

mkdir -p $ACERTS

cp -f $CERT_MASTER/ca* $ACERTS
cp -f $CERT_MASTER/admin-key.pem $ACERTS/kadmin.key
cp -f $CERT_MASTER/admin.pem $ACERTS/kadmin.pem

cp -f ../ansible/roles/keygen/files/api-proxy-openssl.cnf $ACERTS/api-proxy-openssl.cnf

# generate all proxy certs
./gen_server_certs.sh $ACERTS


#HOSTS=$(cat hosts)
#echo $HOSTS
# remove the previous container instance 
CMD="docker rm -f api-proxy"
echo $CMD
$(echo $CMD)

# start new container instance. Map the volume to CERTS
CMD="docker run -v $CERTS:/opt/tls_certs -p 8087:8087 -e env_name=$env_name --name api-proxy api-proxy"

# to run container as a daemon use `-d` flag:
if [ "$2" == "-d" ] ; then
	CMD="docker run -d -v $CERTS:/opt/tls_certs -p 8087:8087 -e env_name=$env_name --name api-proxy api-proxy tail -f /dev/null"
fi
echo $CMD
$(echo $CMD)
