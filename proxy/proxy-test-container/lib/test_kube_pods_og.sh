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

# Help check
if [[ "$1" == "-?" || "$1" == "-h" || "$1" == "--help" || "$1" == "help" ]]; then
	helpme
	exit 1
fi

TEST_COUNT=0
SUCCESS_COUNT=0
date=$( date +%F )
time=$( date +%T )
timestamp="$date""_""$time"

TENANT_ID="$2"
RESULTS_PATH="../logs/""$TENANT_ID""_test_kube_pods_results_""$timestamp"".log"



NUM_PODS=$1
NUM_ONE_PODS=$((NUM_PODS / 2))
NUM_TWO_PODS=$((NUM_PODS - NUM_ONE_PODS)) 


TEST_TYPE="kube"

PROXY_LOC="$3"

# Okay, I know how to create, inspect, delete.

DOCKER_CERT_PATH="../certs/kube"
space="f7f413cb-a678-412d-b024-8e17e28bcb88"
user="d7eae25d39f061dd40937d3839b96fc34d4401391823160f"
KUBE_PATH="kube/kubectl"

ONE_YAML_PATH="../conf/kube/web-test1.yaml"
ONE_NAME="kube-web-test1"
TWO_YAML_PATH="../conf/kube/web-test2.yaml"
TWO_NAME="kube-web-test2"

# Probably gotta set up environment first
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
	 # WHAT EVEN
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
	if [ -f $RESULTS_PATH ]; then 
		rm $RESULTS_PATH
	fi
	touch $RESULTS_PATH

	setup_env

	proxy_location "$PROXY_LOC" $RESULTS_PATH


	test_get_pods 0 

	ONE_COUNTER=0
	while [  $ONE_COUNTER -lt $NUM_ONE_PODS ]; do

		delete_if_exists "$ONE_NAME"
		test_describe_pod "$ONE_NAME" 1 # Fail b/c doesn't exist 
		test_create_pod "$ONE_YAML_PATH" 0 
		test_get_pods 0
		test_describe_pod "$ONE_NAME" 0

		test_create_pod "$ONE_YAML_PATH" 1 # Should fail; already have 1 up.
		test_get_pods 0

		let ONE_COUNTER=ONE_COUNTER+1
	done 

	TWO_COUNTER=0
	while [  $TWO_COUNTER -lt $NUM_TWO_PODS ]; do

		delete_if_exists "$TWO_NAME"
		test_create_pod "$TWO_YAML_PATH" 0 
		test_describe_pod "$TWO_NAME" 0
		test_get_pods 0

		test_create_pod "$TWO_YAML_PATH" 1
		test_get_pods 0

		let TWO_COUNTER=TWO_COUNTER+1
	done 


	# Deleting 
	test_delete_pod "$ONE_NAME" 0
	test_get_pods 0
	test_describe_pod "$ONE_NAME" 1 # Fail b/c doesn't exist 
	test_delete_pod "$ONE_NAME" 1 # Fail b/c doesn't exist

	test_delete_pod "$TWO_NAME" 0 

	test_get_pods 0

	# Report summary of tests
	sum_results "Kube" $RESULTS_PATH $TEST_COUNT $SUCCESS_COUNT $starttime
}

main 


# REMEMBER TO RESTRUCTURE THIS LATER
# CHECK IF DEV-MON01 OTHERWISE ONLY RUN SWARM
