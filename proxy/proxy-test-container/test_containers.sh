#!/bin/bash 

starttime=$(date +%s)

helpme()
{
	cat <<HELPMEHELPME

Syntax: ${0} -l <proxy_location> -n <network_id> -t <tenant_id> -k <test_kube?> -c <num_containers> -p <num_pods> 

* All flags technically optional; default values are set. 

Where:
	proxy_location =
	    local - test is targeting local instance of proxy (localhost:8087). With proxy running on a container, this should be the option in use. (default) 
	    dev-mon01 - test is targeting remote instance (https://containers-api-dev.stage1.ng.bluemix.net:9443)

	network_id = id of network element to be inspected (currently unsupported - default is empty)

	tenant_id = tenant_id to be tested. Currently, default is "test1".

	test_kube? = "true" or "false". Default is true. 

	num_containers = total number of containers to be created; default value is 5. 

	num_pods = total number of pods to be created. If argument not given, default value is 5. 

	  

Note: 
	when running local, specify whether autentication is done by CCSAPI:
		export DOCKER_CONFIG=certs/ccsapi 
	or local fileauth:
		export DOCKER_CONFIG=certs/fileauth

	*** If executing multiple-user fileauth, export DOCKER_CONFIG=certs/fileauth/<user_id> *** 

HELPMEHELPME
}


while test $# -gt 0; do
	case "$1" in 
		""|"-?"|"-h"|"--help"|"help")
			helpme
			exit 1
			;;
		-l)
			shift 
			if test $# -gt 0; then 
				PROXY_LOC=$(echo "$1" | tr '[:upper:]' '[:lower:]')
			else
				echo "Proxy location not specified"
				exit 1
			fi 
			shift 
			;;
		-n)
			shift
			if test $# -gt 0; then 
				try_net_id="$1"
			else 
				echo "No network id specified"
				exit 1
			fi 
			shift
			;; 
		-t)
			shift
			if test $# -gt 0; then 
				TENANT_ID="$1"
			else
				echo "No tenant_id specified"
			fi 
			shift 
			;;
		-k)
			shift 
			if test $# -gt 0; then
				TEST_KUBE=$(echo "$1" | tr '[:upper:]' '[:lower:]')
			else
				echo "Test kube flag not specified"
			fi 
			shift 
			;; 
		-c)
			shift
			if test $# -gt 0; then 
				NUM_CONTAINERS=$1
			else
				echo "Number of containers not specified"
			fi 
			shift
			;;
		-p)
			shift
			if test $# -gt 0; then 
				NUM_PODS=$1
			else 
				echo "Number of pods not specified"
			fi 
			shift 
			;;
	esac
done 


TEST_COUNT=0
SUCCESS_COUNT=0
date=$( date +%F )
time=$( date +%T )
timestamp="$date""_""$time"
RESULTS_PATH="logs/""$TENANT_ID""_test_swarm_results_""$timestamp"".log"
TEST_TYPE="swarm"

# create logs directory 
mkdir -p logs

# Setting default parameters
if [[ "$PROXY_LOC" == "" ]]; then
	PROXY_LOC="local"
fi
if [[ "$NUM_CONTAINERS" == "" ]]; then 
	NUM_CONTAINERS=5
fi 
if [[ "$TEST_KUBE" == "" ]]; then 
	TEST_KUBE=true
fi 
if [[ "$NUM_PODS" == "" ]]; then
	NUM_PODS=5
fi 
if [[ "$TENANT_ID" == "" ]]; then 
	TENANT_ID="test1"
fi 


NUM_NET_CONTAINERS=0 # No network for now
NUM_NO_NET_CONTAINERS=$NUM_CONTAINERS


source lib/common_test_func.sh 


case "$PROXY_LOC" in 
	"local")
		eval "export DOCKER_HOST=localhost:8087"
		eval "export DOCKER_TLS_VERIFY=1"
		CMD_PREFIX="\"\""
		;;
	"dev-mon01")
		eval "export DOCKER_HOST=tcp://containers-api-dev.stage1.ng.bluemix.net:9443"
		eval "export DOCKER_TLS_VERIFY=1"
		eval "export DOCKER_CERT_PATH=certs/dev-mon01"
		eval "export DOCKER_CONFIG=certs/dev-mon01"
		CMD_PREFIX="\"\""
		;;
esac 


# Docker ps
function test_ps() {
	local EXPECTED=$1
	let "TEST_COUNT++"
	echo "Test $TEST_COUNT" && \
	echo "*** Testing docker ps ***" 

	local output
	local timestamp=$(date +"%Y%m%d.%H%M%S")

	local STARTTIME=$(date +%s) 
	output=$(generic_docker_command $CMD_PREFIX "ps" 2>&1)
	local RESULT=$?
	local ENDTIME=$(date +%s)
	
	echo "$output"


	local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker ps" $RESULTS_PATH $EXPECTED "$output" "$timestamp" "$TENANT_ID" "$TEST_TYPE")
	SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))

	echo ""
}

function test_ps_a() {
	local EXPECTED=$1
	let "TEST_COUNT++"
	echo "Test $TEST_COUNT" && \
	echo "*** Testing docker ps -a ***" 

	local output 
	local timestamp=$(date +"%Y%m%d.%H%M%S")

	local STARTTIME=$(date +%s)
	output=$(generic_docker_command $CMD_PREFIX "ps -a" 2>&1)
	local RESULT=$?
	local ENDTIME=$(date +%s)

	echo "$output"
	
	local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker ps -a" $RESULTS_PATH $EXPECTED "$output" "$timestamp" "$TENANT_ID" "$TEST_TYPE")
	SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))
	echo ""
}



function test_create_nonet() {
	local NAME="$1"
	local EXPECTED=$2
	let "TEST_COUNT++"
	echo "Test $TEST_COUNT" && \
	echo "*** Testing container creation w/ no net; name = $NAME; docker run ***" 

	local output 
	local timestamp=$(date +"%Y%m%d.%H%M%S")

	local STARTTIME=$(date +%s)
	output=$(generic_docker_command $CMD_PREFIX "run -d --name $NAME --net none -m 128m mrsabath/web-ms" 2>&1)
	local RESULT=$?
	local ENDTIME=$(date +%s)

	echo "$output"
	
	local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker run $NAME w/o net" $RESULTS_PATH $EXPECTED "$output" "$timestamp" "$TENANT_ID" "$TEST_TYPE")
	SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))
	echo ""
}

function test_create_net() {
	local NAME="$1"
	local EXPECTED=$2
	let "TEST_COUNT++"
	echo "Test $TEST_COUNT" && \
	echo "*** Testing container creation w/net; name = $NAME; docker run ***" 
	
	local output
	local timestamp=$(date +"%Y%m%d.%H%M%S")

	local STARTTIME=$(date +%s)
	output=$(generic_docker_command $CMD_PREFIX "run -d --name $NAME -m 128m mrsabath/web-ms" 2>&1)
	local RESULT=$?
	local ENDTIME=$(date +%s)

	echo "$output"
	

	local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker run $NAME w/net" $RESULTS_PATH $EXPECTED "$output" "$timestamp" "$TENANT_ID" "$TEST_TYPE")
	SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))
	echo ""
}

function test_inspect() {
	local NAME="$1"
	local EXPECTED=$2
	let "TEST_COUNT++"
	echo "Test $TEST_COUNT" && \
	echo "*** Testing docker inspect; name = $NAME ***" 
	

	local output
	local timestamp=$(date +"%Y%m%d.%H%M%S")

	local STARTTIME=$(date +%s)
	output=$(generic_docker_command $CMD_PREFIX "inspect $NAME" 2>&1)
	local RESULT=$?
	local ENDTIME=$(date +%s)

	echo "$output"
	
	local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker inspect $NAME" $RESULTS_PATH $EXPECTED "$output" "$timestamp" "$TENANT_ID" "$TEST_TYPE")
	SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))
	echo ""
}

function test_stop() {
	local NAME="$1"
	local EXPECTED=$2
	let "TEST_COUNT++"
	echo "Test $TEST_COUNT" && \
	echo "*** Testing docker stop; name = $NAME ***" 
	
	local output 
	local timestamp=$(date +"%Y%m%d.%H%M%S")

	local STARTTIME=$(date +%s)
	output=$(generic_docker_command $CMD_PREFIX "stop $NAME" 2>&1)
	local RESULT=$?
	local ENDTIME=$(date +%s)

	echo "$output"
	

	local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker stop $NAME" $RESULTS_PATH $EXPECTED "$output" "$timestamp" "$TENANT_ID" "$TEST_TYPE")
	SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))
	echo ""
}


function test_start() {
	local NAME="$1"
	local EXPECTED=$2
	let "TEST_COUNT++"
	echo "Test $TEST_COUNT" && \
	echo "*** Testing docker start; name = $NAME ***" 


	local output
	local timestamp=$(date +"%Y%m%d.%H%M%S")

	local STARTTIME=$(date +%s)
	output=$(generic_docker_command $CMD_PREFIX "start $NAME" 2>&1)
	local RESULT=$?
	local ENDTIME=$(date +%s)

	echo "$output"

	local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker start $NAME" $RESULTS_PATH $EXPECTED "$output" "$timestamp" "$TENANT_ID" "$TEST_TYPE")
	SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))
	echo ""
}



function test_rm() {
	local NAME="$1"
	local EXPECTED=$2
	let "TEST_COUNT++"
	echo "Test $TEST_COUNT" && \
	echo "*** Testing docker rm; name = $NAME ***" 

	local output 
	local timestamp=$(date +"%Y%m%d.%H%M%S")

	local STARTTIME=$(date +%s)
	output=$(generic_docker_command $CMD_PREFIX "rm $NAME" 2>&1)
	local RESULT=$?
	local ENDTIME=$(date +%s)
	
	echo "$output"

	local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker rm $NAME" $RESULTS_PATH $EXPECTED "$output" "$timestamp" "$TENANT_ID" "$TEST_TYPE")
	SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))
	echo ""
}

function test_rm_f() {
	local NAME="$1"
	local EXPECTED=$2
	let "TEST_COUNT++"
	echo "Test $TEST_COUNT" && \
	echo "*** Testing docker rm -f; name = $NAME ***" 
	
	local output 
	local timestamp=$(date +"%Y%m%d.%H%M%S")

	local STARTTIME=$(date +%s)
	output=$(generic_docker_command $CMD_PREFIX "rm -f $NAME" 2>&1)
	local RESULT=$?
	local ENDTIME=$(date +%s)
	
	echo "$output"

	local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker rm -f $NAME" $RESULTS_PATH $EXPECTED "$output" "$timestamp" "$TENANT_ID" "$TEST_TYPE")
	SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))
	echo ""

}

function test_network_ls() {
	local EXPECTED=$1
	let "TEST_COUNT++"
	echo "Test $TEST_COUNT" && \
	echo "*** Testing docker network ls ***" 
	
	local output 
	local timestamp=$(date +"%Y%m%d.%H%M%S")

	local STARTTIME=$(date +%s)
	output=$(generic_docker_command $CMD_PREFIX "network ls" 2>&1)
	local RESULT=$?
	local ENDTIME=$(date +%s)
	
	echo "$output"

	local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker network ls" $RESULTS_PATH $EXPECTED "$output" "$timestamp" "$TENANT_ID" "$TEST_TYPE")
	SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))
	echo ""
}

function test_network_inspect() {
	local NAME="$1"
	local EXPECTED=$2
	let "TEST_COUNT++"
	echo "Test $TEST_COUNT" && \
	echo "*** Testing docker network inspect; name = $NAME ***" 
	
	local output 
	local timestamp=$(date +"%Y%m%d.%H%M%S")

	local STARTTIME=$(date +%s)
	output=$(generic_docker_command $CMD_PREFIX "network inspect $NAME" 2>&1)
	local RESULT=$?
	local ENDTIME=$(date +%s)
	
	echo "$output"

	local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker network inspect $NAME" $RESULTS_PATH $EXPECTED "$output" "$timestamp" "$TENANT_ID" "$TEST_TYPE")
	SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))
	echo ""
}



function main() {
	if [ -f $RESULTS_PATH ]; then
		rm $RESULTS_PATH
	fi
	touch $RESULTS_PATH

	 
	NONET_COUNTER=1
	while [  $NONET_COUNTER -le $NUM_NO_NET_CONTAINERS ]; do
		test_ps 0 

		test_inspect "$TENANT_ID""_test1" 1 

		test_create_nonet "$TENANT_ID""_nonet_test""$NONET_COUNTER" 0 
		test_ps 0 

		test_inspect "$TENANT_ID""_nonet_test""$NONET_COUNTER" 0 

		test_stop "$TENANT_ID""_test1" 1 # Should fail because doesn't exist
		test_stop "$TENANT_ID""_nonet_test""$NONET_COUNTER" 0 
		test_ps_a 0 
		test_ps 0 

		test_start "$TENANT_ID""_test1" 1 # Doesn't exist
		test_start "$TENANT_ID""_nonet_test""$NONET_COUNTER" 0 

		test_ps_a 0 


		# Test can't remove without stopping first
		test_rm "$TENANT_ID""_nonet_test""$NONET_COUNTER" 1 

		let NONET_COUNTER=NONET_COUNTER+1
	done 


	NET_COUNTER=1
	while [  $NET_COUNTER -le $NUM_NET_CONTAINERS ]; do
		test_ps 0 

		test_create_net "$TENANT_ID""_net_test""$NET_COUNTER" 0 
		test_ps 0 

		test_inspect "$TENANT_ID""_net_test""$NET_COUNTER" 0 

		test_stop "$TENANT_ID""_net_test""$NET_COUNTER" 0 
		test_ps_a 0 
		test_ps 0 

		test_start "$TENANT_ID""_net_test""$NET_COUNTER" 0 

		test_ps_a 0 

		# Test can't remove without stopping first
		test_rm "$TENANT_ID""_net_test""$NET_COUNTER" 1 

		let NET_COUNTER=NET_COUNTER+1
	done 


	NONET_COUNTER=1
	while [  $NONET_COUNTER -le $NUM_NO_NET_CONTAINERS ]; do
		if [ $(( $NONET_COUNTER % 2 )) -eq 0 ]; then 
			test_stop "$TENANT_ID""_nonet_test""$NONET_COUNTER" 0
			test_rm "$TENANT_ID""_nonet_test""$NONET_COUNTER" 0
		else
			test_rm_f "$TENANT_ID""_nonet_test""$NONET_COUNTER" 0 
		fi 

		let NONET_COUNTER=NONET_COUNTER+1
	done 


	test_network_ls 0

	if [[ "$try_net_id" != "" ]]; then 
		test_network_inspect "$try_net_id" 0
	fi 

	# Make sure everything is clean at the end 
	test_ps_a 0
	test_ps 0 

	sum_results "Docker Containers" $RESULTS_PATH $TEST_COUNT $SUCCESS_COUNT $starttime

}

main

echo "$TEST_KUBE"
if [ "$TEST_KUBE" = true ]; then 
	cd lib
	./test_kube_pods.sh $NUM_PODS "$TENANT_ID" "$PROXY_LOC"
fi 

