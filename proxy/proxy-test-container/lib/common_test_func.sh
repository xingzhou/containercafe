#!/bin/bash

function create_log_file() {
    RESULTS_PATH="$1"

	if [ -f $RESULTS_PATH ]; then
		rm $RESULTS_PATH
	fi
	touch $RESULTS_PATH

	echo "DOCKER_HOST=$DOCKER_HOST" >> $RESULTS_PATH
    echo "DOCKER_TLS_VERIFY=$DOCKER_TLS_VERIFY" >> $RESULTS_PATH
    echo "DOCKER_CERT_PATH=$DOCKER_CERT_PATH" >> $RESULTS_PATH
    echo "KUBECONFIG=$KUBECONFIG" >> $RESULTS_PATH
    echo "" >> $RESULTS_PATH
}

function proxy_location() {
	PROXY_LOC="$1"
	RESULTS_PATH="$2"
	echo "$PROXY_LOC" >> $RESULTS_PATH
	echo "" >> $RESULTS_PATH
}

function generic_docker_command() {
	local CMD_PREFIX=$1
	local CMD=$2
	if [[ "$CMD_PREFIX" == "\"\"" ]]; then
		CMD_PREFIX="" # Flag for no prefix
	fi 
	eval $CMD_PREFIX" docker "$CMD
}

function generic_kube_command() {
	local KUBE_PATH=$1
	local CMD=$2
	eval $KUBE_PATH" "$CMD
}

function check_result() {
	RESULT=$1
	TIME=$2
	TEST_COUNT=$3
	SUMMARY="$4"  
	RESULTS_PATH="$5"
	EXPECTED=$6 # 0 = success
	CMD_OUTPUT="$7"
	TIMESTAMP="$8"

	TENANT_ID="$9"
	TYPE="${10}"


	# Process result
	if [ $RESULT -eq $EXPECTED ]; then
		echo "$TIMESTAMP,$TIME,$TENANT_ID,$TYPE,test$TEST_COUNT,$SUMMARY,PASS,OK" >> $RESULTS_PATH
		echo "" >> $RESULTS_PATH
		# Pass back success
		echo 1
	else 
		CMD_OUTPUT=$(echo $CMD_OUTPUT | tr "\n" " ")
		echo "$TIMESTAMP,$TIME,$TENANT_ID,$TYPE,test$TEST_COUNT,$SUMMARY,FAIL,$CMD_OUTPUT" >> $RESULTS_PATH
		echo "" >> $RESULTS_PATH
		# Pass back failure
		echo 0
	fi 
}


# At end - summation of results
function sum_results() {
	TEST_TYPE="$1"
	RESULTS_PATH=$2
	TOTAL_COUNT=$3
	PASS_COUNT=$4
	starttime=$5

	echo "" >> $RESULTS_PATH
	printf '%*s\n' "${COLUMNS:-$(tput cols)}" '' | tr ' ' - >> $RESULTS_PATH
	echo "" >> $RESULTS_PATH

	FAIL_COUNT=$((TOTAL_COUNT - PASS_COUNT))
	echo "$TEST_TYPE Test summary:" >> $RESULTS_PATH
	echo "Total = $TOTAL_COUNT tests" >> $RESULTS_PATH
	echo "Passed = $PASS_COUNT tests" >> $RESULTS_PATH
	echo "Failed = $FAIL_COUNT tests" >> $RESULTS_PATH
	echo "Total time = $(($(date +%s) - $starttime)) sec" >> $RESULTS_PATH

}
