#!/bin/bash

DIR=`dirname "$0"`
source "$DIR/common_test_func.sh"

function main() {
    local LOGS_PATH="$1"
    local TENANT_ID="$2"
    local NUM_PODS=$3
    local KUBE_NAME="${TENANT_ID}-kube-web-test"

    setup_env
    
    init_tests "$TENANT_ID" "kube" "$LOGS_PATH"

    test_get_pods 0

    local COUNTER=1
    while [  $COUNTER -le $NUM_PODS ]; do
        local kube_name="${KUBE_NAME}${COUNTER}"
        local yaml_path=`get_yaml_path "$TENANT_ID" $COUNTER`

        "$DIR/make_yaml.sh" $COUNTER "$TENANT_ID" "$yaml_path"

        delete_if_exists "$kube_name"
        test_create_pod "$yaml_path" 0
        test_describe_pod "$kube_name" 0
        test_get_pods 0

        test_create_pod "$yaml_path" 1
        test_get_pods 0

        let COUNTER=COUNTER+1
    done 

    COUNTER=1
    while [  $COUNTER -le $NUM_PODS ]; do
        local kube_name="${KUBE_NAME}${COUNTER}"
        local yaml_path=`get_yaml_path "$TENANT_ID" $COUNTER`

        test_delete_pod "$kube_name" 0
        test_get_pods 0

        let COUNTER=COUNTER+1
    done

    complete_tests "Kube"
}

function setup_env {
    local PLATFORM=`uname -s | tr '[:upper:]' '[:lower:]'`
    local KUBE_DIR="$DIR/kube"
    local KUBE_PATH="$KUBE_DIR/kubectl"
    export PATH="$KUBE_DIR:$PATH"

    if command -v kubectl >/dev/null 2>&1 ; then
        echo "Using system kubectl..."
    elif [[ ! -f "$KUBE_PATH" ]]; then
        mkdir -p "$KUBE_DIR"
        local LATEST_KUBECLT=`curl -L https://storage.googleapis.com/kubernetes-release/release/stable.txt`
        curl -sSLo "$KUBE_PATH" "https://storage.googleapis.com/kubernetes-release/release/$LATEST_KUBECLT/bin/$PLATFORM/amd64/kubectl"
        chmod +x "$KUBE_PATH"
    fi

    if [[ "$KUBECONFIG" == "" ]]; then
        export KUBECONFIG="$DIR/conf/kube/kubeconfig-radiant01.yaml"
    fi
    
    kubectl config view

    echo "" 
    echo "" 
}

function get_yaml_path {
    echo "$DIR/../conf/kube/${1}-web-test${2}.yaml"
}

function delete_if_exists {
    local NAME="$1"
    echo "*** Deleting pod $NAME if exists ***"
    kubectl delete pod "$NAME"
    echo ""
}

function test_get_pods {
    assert kubectl get pods \
           --equal $1 \
           --log "get pods"
}

function test_create_pod {
    assert kubectl create -f "$1" \
           --equal $2 \
           --log "create pod; pod = $1"
}

function test_describe_pod {
    assert kubectl describe pod "$1" \
           --equal $2 \
           --log "describe pod; pod = $1"
}

function test_delete_pod {
    assert kubectl delete pod "$1" \
           --equal $2 \
           --log "delete pod; pod = $1"
}

main "$@"
