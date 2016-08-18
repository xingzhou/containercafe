#!/bin/bash

function main {
    local tenants=()
    if [[ `echo "$CI" | tr '[:upper:]' '[:lower:]'` == true && `check_for_tenants "$@"` == false ]]; then
        IFS=$'\n' tenants=(`docker exec api-proxy cat /api-proxy/creds.json \
            | grep -oP "((?<=\"TLS_path\":\")|(?<=\"Space_id\":\"))[^\"]+" \
            | rev \
            | sed 'N;s/\(.*\)\n\([^/]*\)\(\/\)\([^/]*\)\(.*\)/t-\n\2\:\1\:\4/' \
            | rev`)
    fi

    local CERTS_DIR=`docker_image_env CERTS_DIR`
    local LOGS_DIR=`docker_image_env LOGS_DIR`

    docker run -v "$HOME/.openradiant/envs":"$CERTS_DIR":ro -v api-proxy-tests-logs:"$LOGS_DIR" --net="host" api-proxy-tests "${tenants[@]}" "$@"
}

function docker_image_env {
    docker inspect -f "{{ .Config.Env }}" api-proxy-tests | grep -oP "(?<=$1\=)([^\s\]]+)"
}

function check_for_tenants {
    local have_tenant=false
    
    while [[ $# -gt 0 && $have_tenant == false ]]; do
        case "$1" in 
        -t)
            shift
            have_tenant=true
            ;;
        *)
            shift
            ;;
        esac
    done
    
    echo "$have_tenant"
}

main "$@"
