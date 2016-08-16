#!/bin/bash
starttime=$(date +%s)

source ./common_test_func.sh 

helpme() 
{
    cat <<HELPMEHELPME

Syntax: ${0} 
Summary of results in ../logs/test_kube_pods_results.log

HELPMEHELPME
}


if [[ "$1" == "-?" || "$1" == "-h" || "$1" == "--help" || "$1" == "help" ]]; then
    helpme
    exit 1
fi

TENANT_ID="$2"
TEST_COUNT=0
SUCCESS_COUNT=0
if [[ -z "$LOGS_DIR" ]]; then
    LOGS_DIR="../logs"
fi
if [[ -z "$LOG_SUFFIX" ]]; then 
    LOG_SUFFIX="$( date +%F )_$( date +%H-%M-%S )"
fi
if [[ ! -z "$LOG_PREFIX" ]]; then 
    SEPARATOR="_"
else
    SEPARATOR=""
fi
RESULTS_PATH="${LOGS_DIR}/${LOG_PREFIX}${SEPARATOR}${TENANT_ID}_test_kube_pods_results_${LOG_SUFFIX}.log"

NUM_PODS=$1
TEST_TYPE="kube"

PROXY_LOC="$3"

space="f7f413cb-a678-412d-b024-8e17e28bcb88"
user="d7eae25d39f061dd40937d3839b96fc34d4401391823160f"
PLATFORM=`uname -s | tr '[:upper:]' '[:lower:]'`
KUBE_PATH="kube/kubectl"

if command -v kubectl >/dev/null 2>&1 ; then
    KUBE_PATH="kubectl"
elif [[ ! -f "$KUBE_PATH" ]]; then
    mkdir -p $(dirname "$KUBE_PATH")
    LATEST_KUBECLT=`curl -L https://storage.googleapis.com/kubernetes-release/release/stable.txt`
    curl -sSLo "$KUBE_PATH" https://storage.googleapis.com/kubernetes-release/release/$LATEST_KUBECLT/bin/$PLATFORM/amd64/kubectl
    chmod +x "$KUBE_PATH"
fi

YAML_PATH_A="../conf/kube/$TENANT_ID-web-test"
YAML_PATH_B=".yaml"
KUBE_NAME="$TENANT_ID-kube-web-test"

function setup_env() {
    if [[ $KUBECONFIG == "" ]]; then 
        eval "export KUBECONFIG=../conf/kube/kubeconfig-radiant01.yaml"
    fi 
    generic_kube_command $KUBE_PATH "config view"

    echo ""
    echo ""
    echo "" 
    echo "" 
}

function test_get_pods() {
     local EXPECTED=$1
     let "TEST_COUNT++"
     echo "Test $TEST_COUNT" && \
     echo "*** Testing get pods ***" 

     local output 
     local timestamp=$(date +"%Y%m%d.%H%M%S")

     local STARTTIME=$(date +%s)
     output=$(generic_kube_command $KUBE_PATH "get pods" 2>&1)
     local RESULT=$?
     local ENDTIME=$(date +%s)

     echo "$output"

     local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Get pods" $RESULTS_PATH $EXPECTED "$output" "$timestamp" "$TENANT_ID" "$TEST_TYPE")
     SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))

     echo ""
}

function delete_if_exists() {
    local NAME="$1"
    echo "*** Deleting pod $NAME if exists ***" && \
    generic_kube_command $KUBE_PATH "delete pod $NAME"
    echo ""
}

function test_create_pod() {
    local SPEC_FILE="$1"
    local EXPECTED=$2
    let "TEST_COUNT++"
    echo "Test $TEST_COUNT" && \
    echo "*** Testing create pod with $SPEC_FILE ***" 

    local output 
    local timestamp=$(date +"%Y%m%d.%H%M%S")

    local STARTTIME=$(date +%s)
    output=$(generic_kube_command $KUBE_PATH "create -f $SPEC_FILE" 2>&1)
    local RESULT=$?
    local ENDTIME=$(date +%s)

    echo "$output"

    local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Create pod with $SPEC_FILE" $RESULTS_PATH $EXPECTED "$output" "$timestamp" "$TENANT_ID" "$TEST_TYPE")
    SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))

    echo ""
}

function test_describe_pod() {
    local NAME="$1"
    local EXPECTED=$2
    let "TEST_COUNT++"
    echo "Test $TEST_COUNT" && \
    echo "*** Testing describe pod $NAME ***" 

    local output
    local timestamp=$(date +"%Y%m%d.%H%M%S")

    local STARTTIME=$(date +%s)
    output=$(generic_kube_command $KUBE_PATH "describe pod $NAME" 2>&1)
    local RESULT=$?
    local ENDTIME=$(date +%s)

    echo "$output"

    local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Describe pod $NAME" $RESULTS_PATH $EXPECTED "$output" "$timestamp" "$TENANT_ID" "$TEST_TYPE")
    SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))

    echo "" 
}

function test_delete_pod() {
    local NAME="$1"
    local EXPECTED=$2
    let "TEST_COUNT++"
    echo "Test $TEST_COUNT" && \
    echo "*** Testing delete pod $NAME ***" 

    local output
    local timestamp=$(date +"%Y%m%d.%H%M%S")

    local STARTTIME=$(date +%s)
    output=$(generic_kube_command $KUBE_PATH "delete pod $NAME" 2>&1)
    local RESULT=$?
    local ENDTIME=$(date +%s)

    echo "$output"

    local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Delete pod $NAME" $RESULTS_PATH $EXPECTED "$output" "$timestamp" "$TENANT_ID" "$TEST_TYPE")
    SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))
    
    echo ""
}




function main() {
    setup_env
    
    create_log_file $RESULTS_PATH

    proxy_location "$PROXY_LOC" $RESULTS_PATH

    test_get_pods 0 

    COUNTER=1
    while [  $COUNTER -le $NUM_PODS ]; do
        local kube_name="$KUBE_NAME""$COUNTER"
        local yaml_path="$YAML_PATH_A""$COUNTER""$YAML_PATH_B"

        ./make_yaml.sh $COUNTER "$TENANT_ID"

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
        local kube_name="$KUBE_NAME""$COUNTER"
        local yaml_path="$YAML_PATH_A""$COUNTER""$YAML_PATH_B"

        test_delete_pod "$kube_name" 0
        test_get_pods 0

        let COUNTER=COUNTER+1
    done 

    
    sum_results "Kube" $RESULTS_PATH $TEST_COUNT $SUCCESS_COUNT $starttime
}

main 

