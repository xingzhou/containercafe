#!/bin/bash

function init_tests {
    TENANT_ID="$1"
    TEST_TYPE="$2"
    RESULTS_PATH="$3"
    TEST_COUNT=0
    SUCCESS_COUNT=0
    TEST_START_TIME=$(date +%s)
    
    create_log_file "$RESULTS_PATH"
}

function complete_tests {
    sum_results "$1" "$RESULTS_PATH" $TEST_COUNT $SUCCESS_COUNT $TEST_START_TIME
}

function assert {
    local TEST_CMD=()
    local EXPECTED=  
    local TEST_NAME= 
    
    while test $# -gt 0; do
        case "$1" in 
            "--equal")
                EXPECTED=$2
                shift 2
                ;;
            "--log")
                TEST_NAME="$2"
                shift 2
                ;;
            *)
                TEST_CMD+=("$1")
                shift
                ;;
        esac
    done
    
    let "TEST_COUNT++"
    echo "Test $TEST_COUNT" && \
    echo "*** Testing $TEST_NAME ***" 

    local output 
    local timestamp=$(date +"%Y%m%d.%H%M%S")

    local STARTTIME=$(date +%s)
    output=$("${TEST_CMD[@]}" 2>&1)
    local RESULT=$?
    local ENDTIME=$(date +%s)

    echo "$output"
    
    local log_command=$(printf " %s" "${TEST_CMD[@]}")
    log_command=${log_command:1}
    
    local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "$log_command" "$RESULTS_PATH" $EXPECTED "$output" "$timestamp" "$TENANT_ID" "$TEST_TYPE")
    SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))
    echo ""
}

function create_log_file {
    local LOG_FILE="$1"
    
    mkdir -p "$(dirname "$LOG_FILE")"

    if [ -f "$LOG_FILE" ]; then
        rm "$LOG_FILE"
    fi
    touch "$LOG_FILE"

    echo "DOCKER_HOST=$DOCKER_HOST" >> "$LOG_FILE"
    echo "DOCKER_TLS_VERIFY=$DOCKER_TLS_VERIFY" >> "$LOG_FILE"
    echo "DOCKER_CERT_PATH=$DOCKER_CERT_PATH" >> "$LOG_FILE"
    echo "KUBECONFIG=$KUBECONFIG" >> "$LOG_FILE"
    echo "" >> "$LOG_FILE"
}

function check_result {
    local RESULT=$1
    local TIME=$2
    local TEST_COUNT=$3
    local SUMMARY="$4"
    local LOG_FILE="$5"
    local EXPECTED=$6 # 0 = success
    local CMD_OUTPUT="$7"
    local TIMESTAMP="$8"

    local TENANT_ID="$9"
    local TYPE="${10}"

    # Process result
    if [ $RESULT -eq $EXPECTED ]; then
        echo "$TIMESTAMP,$TIME,$TENANT_ID,$TYPE,test$TEST_COUNT,$SUMMARY,PASS,OK" >> "$LOG_FILE"
        echo "" >> "$LOG_FILE"
        # Pass back success
        echo 1
    else
        CMD_OUTPUT=$(echo "$CMD_OUTPUT" | tr "\n" " ")
        echo "$TIMESTAMP,$TIME,$TENANT_ID,$TYPE,test$TEST_COUNT,$SUMMARY,FAIL,$CMD_OUTPUT" >> "$LOG_FILE"
        echo "" >> "$LOG_FILE"
        # Pass back failure
        echo 0
    fi
}

# At end - summation of results
function sum_results {
    local TEST_TYPE="$1"
    local LOG_FILE=$2
    local TOTAL_COUNT=$3
    local PASS_COUNT=$4
    local starttime=$5

    echo "" >> "$LOG_FILE"
    printf '%*s\n' "${COLUMNS:-$(tput cols)}" '' | tr ' ' - >> "$LOG_FILE"
    echo "" >> "$LOG_FILE"

    FAIL_COUNT=$((TOTAL_COUNT - PASS_COUNT))
    echo "$TEST_TYPE Test summary:" >> "$LOG_FILE"
    echo "Total = $TOTAL_COUNT tests" >> "$LOG_FILE"
    echo "Passed = $PASS_COUNT tests" >> "$LOG_FILE"
    echo "Failed = $FAIL_COUNT tests" >> "$LOG_FILE"
    echo "Total time = $(($(date +%s) - $starttime)) sec" >> "$LOG_FILE"
}
