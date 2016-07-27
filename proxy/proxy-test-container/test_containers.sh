#!/bin/bash 

starttime=$(date +%s)

helpme()
{
	cat <<HELPMEHELPME

Syntax: ${0} <proxy_location> <network_id> <tenant_id> <test_kube?> <num_containers> <num_pods> 
Where:
	proxy_location =
	    local - test is targeting local instance of proxy (localhost:8087) (default if argument not given)
	    dev-mon01 - test is targeting remote instance (https://containers-api-dev.stage1.ng.bluemix.net:9443)

	network_id = id of network element to be inspected (using default as default)

	tenant_id = user_id of current user; should match Apikey being used in the config file. Currently, default is "test1".

	test_kube? = "true" or "false". In case of multi-user testing, only 1 user should run the kube tests ("true"); all others should be "false". 

	num_containers = total number of containers to be created; RIGHT NOW ALL W/O NETWORK. 

	num_pods = total number of pods to be created. if argument not given, default value is 5. 

	  

Note: 
	when running local, specify whether autentication is done by CCSAPI:
		export DOCKER_CONFIG=certs/ccsapi 
	or local fileauth:
		export DOCKER_CONFIG=certs/fileauth

	*** If executing multiple-user fileauth, export DOCKER_CONFIG=certs/fileauth/<user_id> *** 

HELPMEHELPME
}

# helpme()
# {
# 	cat <<HELPMEHELPME

# Syntax: ${0} <proxy_location> <network_id> <tenant_id> <test_kube?> <num_containers> <num_pods> 
# Where:
# 	proxy_location =
# 	    local - test is targeting local instance of proxy (localhost:8087) (default if argument not given)
# 	    dev-mon01 - test is targeting remote instance (https://containers-api-dev.stage1.ng.bluemix.net:9443)

# 	network_id = id of network element to be inspected (using default as default)

# 	tenant_id = user_id of current user; should match Apikey being used in the config file. Currently, default is "test1".

# 	test_kube? = "true" or "false". In case of multi-user testing, only 1 user should run the kube tests ("true"); all others should be "false". 

# 	num_containers = total number of containers to be created; equal split between containers with net, and containers without. 
# 		If argument not given, default value is 5. 

# 	num_pods = total number of pods to be created. if argument not given, default value is 5. 

	  

# Note: 
# 	when running local, specify whether autentication is done by CCSAPI:
# 		export DOCKER_CONFIG=certs/ccsapi 
# 	or local fileauth:
# 		export DOCKER_CONFIG=certs/fileauth

# 	*** If executing multiple-user fileauth, export DOCKER_CONFIG=certs/fileauth/<user_id> *** 

# HELPMEHELPME
# }


# Help check
if [[ "$1" == "" || "$1" == "-?" || "$1" == "-h" || "$1" == "--help" || "$1" == "help" ]] ; then
	helpme
	exit 1
fi


TEST_COUNT=0
SUCCESS_COUNT=0
date=$( date +%F )
time=$( date +%T )
timestamp="$date""_""$time"

TENANT_ID="$3"
RESULTS_PATH="logs/""$TENANT_ID""_test_swarm_results_""$timestamp"".log"
PROXY_LOC=$(echo "$1" | tr '[:upper:]' '[:lower:]')
try_net_id="$2"


TEST_KUBE=$(echo "$4" | tr '[:upper:]' '[:lower:]')
NUM_CONTAINERS=$5
NUM_PODS=$6



TEST_TYPE="swarm"

# create logs directory 
mkdir -p logs

# Setting default parameters
if [[ "$PROXY_LOC" == "" ]]; then
	PROXY_LOC="local"
fi
if [[ "$try_net_id" == "" ]]; then
	try_net_id="default"
fi 
if [[ "$NUM_CONTAINERS" == "" ]]; then 
	NUM_CONTAINERS=5
fi 
if [[ "$NUM_PODS" == "" ]]; then
	NUM_PODS=5
fi 
if [[ "$TENANT_ID" == "" ]]; then 
	TENANT_ID="test1"
fi 


# NUM_NET_CONTAINERS=$((NUM_CONTAINERS / 2))
# NUM_NO_NET_CONTAINERS=$((NUM_CONTAINERS - NUM_NET_CONTAINERS))

NUM_NET_CONTAINERS=0 # No network for now
NUM_NO_NET_CONTAINERS=$NUM_CONTAINERS



source lib/common_test_func.sh # I HOPE THIS WORKS


# Setting the proxy location stuff
case "$PROXY_LOC" in 
	"local")
		eval "export DOCKER_HOST=localhost:8087"
		eval "export DOCKER_TLS_VERIFY=1"
		#CMD_PREFIX="DOCKER_TLS_VERIFY=\"\""
		CMD_PREFIX="\"\""
		;;
	"dev-mon01")
		eval "export DOCKER_HOST=tcp://containers-api-dev.stage1.ng.bluemix.net:9443"
		eval "export DOCKER_TLS_VERIFY=1"
		eval "export DOCKER_CERT_PATH=certs/dev-mon01"
		eval "export DOCKER_CONFIG=certs/dev-mon01"
		# eval "title=\"dev-mon01 auth\""
		# eval "echo -n -e "\033]0;$title\007""
		CMD_PREFIX="\"\""
		;;
esac 



function location() {
	case "$PROXY_LOC" in 
		"local")
			if [[ $DOCKER_CONFIG == *"fileauth"* ]]; then 
				PROXY_LOC="local fileauth"
				proxy_location "local fileauth" $RESULTS_PATH
			elif [[ $DOCKER_CONFIG == *"ccsapi"* ]]; then 
				PROXY_LOC="local ccsapi"
				proxy_location "local ccsapi" $RESULTS_PATH
			else 
				helpme
				exit 1
			fi 
			;;
		"dev-mon01")
			PROXY_LOC="dev-mon01"
			proxy_location "dev-mon01" $RESULTS_PATH
			;;
	esac 

}


# Docker ps
function test_ps() {
	local EXPECTED=$1
	let "TEST_COUNT++"
	echo "Test $TEST_COUNT" && \
	echo "*** Testing docker ps ***" 

	local output
	local timestamp=$(date +"%Y%m%d.%H%M%S")

	# local STARTTIME=$(date +%s)
	local STARTTIME=$(date +%s) # milliseconds
	output=$(generic_docker_command $CMD_PREFIX "ps" 2>&1)
	local RESULT=$?
	local ENDTIME=$(date +%s)
	
	echo "$output"


	local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker ps" $RESULTS_PATH $EXPECTED "$output" "$timestamp" "$TENANT_ID" "$TEST_TYPE")
	SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))

	# So test count increment, success count increment, write result to file. 
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
	

	# local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker run; no net; name = $NAME" $RESULTS_PATH $EXPECTED "$output" "$timestamp" "$TENANT_ID" "$TEST_TYPE")
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
	

	# local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker run; w/net; name = $NAME" $RESULTS_PATH $EXPECTED "$output" "$timestamp" "$TENANT_ID" "$TEST_TYPE")
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
	
	# local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker inspect; name = $NAME" $RESULTS_PATH $EXPECTED "$output" "$timestamp" "$TENANT_ID" "$TEST_TYPE")
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
	

	# local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker stop; name = $NAME" $RESULTS_PATH $EXPECTED "$output" "$timestamp" "$TENANT_ID" "$TEST_TYPE")
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

	# local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker start; name = $NAME" $RESULTS_PATH $EXPECTED "$output" "$timestamp" "$TENANT_ID" "$TEST_TYPE")
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

	# local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker rm; name = $NAME" $RESULTS_PATH $EXPECTED "$output" "$timestamp" "$TENANT_ID" "$TEST_TYPE")
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

	# local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker rm -f; name = $NAME" $RESULTS_PATH $EXPECTED "$output" "$timestamp" "$TENANT_ID" "$TEST_TYPE")
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

	# local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker network inspect; name = $NAME" $RESULTS_PATH $EXPECTED "$output" "$timestamp" "$TENANT_ID" "$TEST_TYPE")
	local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker network inspect $NAME" $RESULTS_PATH $EXPECTED "$output" "$timestamp" "$TENANT_ID" "$TEST_TYPE")
	SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))
	echo ""
}



function main() {
	if [ -f $RESULTS_PATH ]; then
		rm $RESULTS_PATH
	fi
	touch $RESULTS_PATH


	#location 

	# Try the no net ones first. 
	NONET_COUNTER=1
	while [  $NONET_COUNTER -le $NUM_NO_NET_CONTAINERS ]; do
		test_ps 0 # Test 1

		test_inspect "$TENANT_ID""_test1" 1 # Test 2

		test_create_nonet "$TENANT_ID""_nonet_test""$NONET_COUNTER" 0 # Test 3
		test_ps 0 # Test 4

		# Test 5
		# Test 6

		test_inspect "$TENANT_ID""_nonet_test""$NONET_COUNTER" 0 # Test 7 

		# Docker stop
		test_stop "$TENANT_ID""_test1" 1 # Should fail because doesn't exist
		test_stop "$TENANT_ID""_nonet_test""$NONET_COUNTER" 0 # Test 8 
		test_ps_a 0 # Test 9
		test_ps 0 # Test 10

		# Test start 

		# Start - if resource dne, should error, if already up, nothing, and should start it up if stopped. 
		#echo "Docker start nonet_test1 should fail" >> $RESULTS_PATH
		test_start "$TENANT_ID""_test1" 1 # Fails because doesn't exist
		test_start "$TENANT_ID""_nonet_test""$NONET_COUNTER" 0 # Test 14

		test_ps_a 0 # Test 16


		# Test can't remove without stopping first
		test_rm "$TENANT_ID""_nonet_test""$NONET_COUNTER" 1 

		let NONET_COUNTER=NONET_COUNTER+1
	done 


	# Now the net ones. 
	NET_COUNTER=1
	while [  $NET_COUNTER -le $NUM_NET_CONTAINERS ]; do
		test_ps 0 # Test 1

		test_create_net "$TENANT_ID""_net_test""$NET_COUNTER" 0 # Test 3
		test_ps 0 # Test 4

		# Test 5
		# Test 6

		test_inspect "$TENANT_ID""_net_test""$NET_COUNTER" 0 # Test 7 

		# Docker stop
		test_stop "$TENANT_ID""_net_test""$NET_COUNTER" 0 # Test 8 
		test_ps_a 0 # Test 9
		test_ps 0 # Test 10

		# Test start 

		# Start - if resource dne, should error, if already up, nothing, and should start it up if stopped. 
		#echo "Docker start nonet_test1 should fail" >> $RESULTS_PATH
		test_start "$TENANT_ID""_net_test""$NET_COUNTER" 0 # Test 14

		test_ps_a 0 # Test 16


		# Test can't remove without stopping first
		test_rm "$TENANT_ID""_net_test""$NET_COUNTER" 1 

		let NET_COUNTER=NET_COUNTER+1
	done 



	# REMOVE THEM ALL AT THE ENNNDDD
	# Test 11
	# Test 12
	# Test 13

	# Trying removing some with rm -f too 
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


	# NET_COUNTER=1
	# while [  $NET_COUNTER -le $NUM_NET_CONTAINERS ]; do
	# 	test_rm_f "$TENANT_ID""_net_test""$NET_COUNTER" 0 

	# 	let NET_COUNTER=NET_COUNTER+1
	# done 



	test_network_ls 0

	# REMOVING NETWORK INSPECT FOR NOW 

	# test_network_inspect $try_net_id 0

	# test_network_inspect "default" 0 


	# Make sure everything is clean at the end 
	test_ps_a 0
	test_ps 0 



	# SUM THE RESULTS
	sum_results "Docker Containers" $RESULTS_PATH $TEST_COUNT $SUCCESS_COUNT $starttime

}

main


if [[ "$TEST_KUBE" == "true" ]]; then 
	cd lib
	./test_kube_pods.sh $NUM_PODS "$TENANT_ID" "$PROXY_LOC"
fi 

