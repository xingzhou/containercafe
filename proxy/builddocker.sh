#!/bin/bash
# set -v

function build_proxy_image {
    local name="$1"
    local context="$2"
    
    # copy src into Dockerfile directory
    [[ -d "$context/src" ]] && rm -rf "$context/src"
    cp -r src "$context"

    # copy scripts into Dockerfile directory
    cp make_TLS_certs.sh mk_user_cert.sh mk_kubeconfig.sh "$context"
    
    # build the docker image
    docker build -t "$name" "$context"
}

build_proxy_image "api-proxy" "dockerize/"
