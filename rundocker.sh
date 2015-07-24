#!/bin/bash
set -v
cd dockerize
HOSTS=$(cat hosts)
echo $HOSTS
CMD="docker run $HOSTS hijack"
echo $CMD
$(echo $CMD)
