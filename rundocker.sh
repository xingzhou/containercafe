#!/bin/bash
set -v
cd dockerize
#HOSTS=$(cat hosts)
#echo $HOSTS
CMD="docker run -v `pwd`:/opt/tls_certs -p 8087:8087 hijack"
echo $CMD
$(echo $CMD)
