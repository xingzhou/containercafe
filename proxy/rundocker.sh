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

function main {
    if [[ "$1" == "" || "$1" == "-d" ]] ; then
        echo 'env_name must be set'
        helpme
        exit 1
    fi
    local env_name="$1"

    # manage certificates - copy them to admin-certs
    local CERT_MASTER="../ansible/certs/dev-vbox-radiant01"
    local CERTS="$HOME/.openradiant/envs/$env_name"
    local ACERTS="$CERTS/admin-certs"
    if [ ! -d "$CERT_MASTER" ]; then
      echo "missing $CERT_MASTER directory. Execute ansible scripts first"
      exit 99
    fi

    if [[ ! -d "$ACERTS" || `verify_certs_matching "$CERT_MASTER" "$ACERTS"` = false ]]; then
        [[ -d "$ACERTS" ]] && rm -rf "$ACERTS"
        mkdir -p "$ACERTS"
        
        copy_certs "$CERT_MASTER" "$ACERTS"
    else
        echo "WARNING: using the existing certs in $ACERTS"
        echo "To recreate the certs, delete this directory"
    fi

    # create an empty creds.json if necessary
    touch "$CERTS/creds.json"

    # to run container as a daemon use `-d` flag:
    local EXTRA_FLAGS=()
    if [ "$2" == "-d" ] ; then
        EXTRA_FLAGS+=("-d")
    fi

    set -x

    # remove the previous container instance 
    docker rm -f api-proxy

    # start new container instance. Map the volume to CERTS
    docker run "${EXTRA_FLAGS[@]}" -v "$CERTS":/opt/tls_certs -p 8087:8087 -e "env_name=$env_name" --name api-proxy api-proxy
}

function copy_certs {
    local CERT_MASTER="$1"
    local ACERTS="$2"

    # copy the master certs
    cp -f "$CERT_MASTER/ca"* "$ACERTS"
    cp -f "$CERT_MASTER/admin-key.pem" "$ACERTS/kadmin.key"
    cp -f "$CERT_MASTER/admin.pem" "$ACERTS/kadmin.pem"

    cp -f ../ansible/roles/keygen-shard/files/api-proxy-openssl.cnf "$ACERTS/api-proxy-openssl.cnf"

    # generate all proxy certs
    ./gen_server_certs.sh "$ACERTS"
}

function verify_certs_matching {
    local CERT_MASTER="$1"
    local ACERTS="$2"

    local matching=true
    compare_certs "$CERT_MASTER/ca.pem" "$ACERTS/ca.pem" || matching=false
    compare_certs "$CERT_MASTER/admin.pem" "$ACERTS/kadmin.pem" || matching=false

    echo "$matching"
}

function compare_certs {
    local CERT1="$1"
    local CERT2="$2"
    
    if [[ -f "$CERT1" && -f "$CERT2" && `diff --brief <(cert_to_text "$CERT1") <(cert_to_text "$CERT2")` = "" ]]; then
        return 0
    else
        return 1
    fi
}

function cert_to_text {
    local CERT_PATH="$1"
    
    openssl x509 -in "$CERT_PATH" -text -noout 2>&1
}

# Main

main "$@"
