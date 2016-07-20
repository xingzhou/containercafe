#!/bin/bash
set -v
cd dockerize
#HOSTS=$(cat hosts)
#echo $HOSTS
# remove the previous container instance 
CMD="docker rm -f hjproxy"
echo $CMD
$(echo $CMD)

# start new container instance
CMD="docker run -v `pwd`:/opt/tls_certs -p 8087:8087 --name hjproxy proxy"

# to run container as a daemon use `-d` flag:
if [ "$1" == "-d" ] ; then
	CMD="docker run -d -v `pwd`:/opt/tls_certs -p 8087:8087 --name hjproxy proxy"
fi
echo $CMD
$(echo $CMD)
