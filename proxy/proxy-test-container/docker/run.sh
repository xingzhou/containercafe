#!/bin/bash

tenants=()
if [[ `echo "$CI" | tr '[:upper:]' '[:lower:]'` == "true" ]]; then
    IFS=$'\n' tenants=(`docker exec api-proxy cat /api-proxy/creds.json \
        | tail -n +2 \
        | grep -oP "((?<=\"TLS_path\":\")|(?<=\"Space_id\":\"))[^\"]+" \
        | rev \
        | sed 'N;s/\(.*\)\n\([^/]*\)\(\/\)\([^/]*\)\(.*\)/t-\n\2\:\1\:\4/' \
        | rev`)
fi

docker run -v ~/.openradiant/envs:/tests/certs:ro -v api-proxy-tests-logs:/tests/logs --net="host" api-proxy-tests "${tenants[@]}" "$@"
