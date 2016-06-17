#!/bin/bash 

# DUDE I HAVE NO IDEA WHAT'S GOING ON DRILUGHDSLIUGHDSLIUGHLISDUHGLIDSUHGLIUSDH

#??????? TIME TO DO RANDOM STUFF AND HOPE FOR THE BEST

TEST_COUNT=0
SUCCESS_COUNT=0
RESULTS_PATH="test_containers_results.txt"
try_net_id="$1"

source ./common_test_func_copy.sh # I HOPE THIS WORKS


# Docker ps
function test_ps() {
	local STARTTIME=$(date +%s)
	echo "*** Testing docker ps ***" && \
	generic_docker_command "ps"



	local RESULT=$?
	local ENDTIME=$(date +%s)
	let "TEST_COUNT++"

	local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker ps" $RESULTS_PATH)
	SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))

	# So test count increment, success count increment, write result to file. 
	echo ""
}

function test_ps_a() {
	local STARTTIME=$(date +%s)
	echo "*** Testing docker ps -a ***" && \
	generic_docker_command "ps -a"

	local RESULT=$?
	local ENDTIME=$(date +%s)
	let "TEST_COUNT++"

	local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker ps -a" $RESULTS_PATH)
	SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))
	echo ""
}



function test_create_nonet() {
	local NAME="nonet_""$1"
	local STARTTIME=$(date +%s)
	echo "*** Testing container creation w/ no net; name = $NAME; docker run ***" && \
	generic_docker_command "run -d --name $NAME --net none -m 128m mrsabath/web-ms"

	local RESULT=$?
	local ENDTIME=$(date +%s)
	let "TEST_COUNT++"

	local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker run; no net; name = $NAME" $RESULTS_PATH)
	SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))
	echo ""
}

function test_create_net() {
	local NAME="$1"
	local STARTTIME=$(date +%s)
	echo "*** Testing container creation w/net; name = $NAME; docker run ***" && \
	generic_docker_command "run -d --name $NAME -m 128m mrsabath/web-ms"

	local RESULT=$?
	local ENDTIME=$(date +%s)
	let "TEST_COUNT++"

	local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker run; w/net; name = $NAME" $RESULTS_PATH)
	SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))
	echo ""
}

function test_inspect() {
	local NAME="$1"
	local STARTTIME=$(date +%s)
	echo "*** Testing docker inspect; name = $NAME ***" && \
	generic_docker_command "inspect $NAME"

	local RESULT=$?
	local ENDTIME=$(date +%s)
	let "TEST_COUNT++"

	local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker inspect; name = $NAME" $RESULTS_PATH)
	SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))
	echo ""
}

function test_stop() {
	local NAME="$1"
	local STARTTIME=$(date +%s)
	echo "*** Testing docker stop; name = $NAME ***" && \
	generic_docker_command "stop $NAME"
	
	local RESULT=$?
	local ENDTIME=$(date +%s)
	let "TEST_COUNT++"

	local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker stop; name = $NAME" $RESULTS_PATH)
	SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))
	echo ""
}


function test_start() {
	local NAME="$1"
	local STARTTIME=$(date +%s)
	echo "*** Testing docker start; name = $NAME ***" && \
	generic_docker_command "start $NAME"

	local RESULT=$?
	local ENDTIME=$(date +%s)
	let "TEST_COUNT++"

	local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker start; name = $NAME" $RESULTS_PATH)
	SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))
	echo ""
}



function test_rm() {
	local NAME="$1"
	local STARTTIME=$(date +%s)
	echo "*** Testing docker rm; name = $NAME ***" && \
	generic_docker_command "rm $NAME"

	local RESULT=$?
	local ENDTIME=$(date +%s)
	let "TEST_COUNT++"

	local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker rm; name = $NAME" $RESULTS_PATH)
	SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))
	echo ""
}

function test_rm_f() {
	local NAME="$1"
	local STARTTIME=$(date +%s)
	echo "*** Testing docker rm -f; name = $NAME ***" && \
	generic_docker_command "rm -f $NAME"

	local RESULT=$?
	local ENDTIME=$(date +%s)
	let "TEST_COUNT++"

	local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker rm -f; name = $NAME" $RESULTS_PATH)
	SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))
	echo ""

}

function test_network_ls() {
	local STARTTIME=$(date +%s)
	echo "*** Testing docker network ls ***" && \
	generic_docker_command "network ls"

	local RESULT=$?
	local ENDTIME=$(date +%s)
	let "TEST_COUNT++"

	local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker network ls" $RESULTS_PATH)
	SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))
	echo ""
}

function test_network_inspect() {
	local NAME="$1"
	local STARTTIME=$(date +%s)
	echo "*** Testing docker network inspect; name = $NAME ***" && \
	generic_docker_command "network inspect $NAME"

	local RESULT=$?
	local ENDTIME=$(date +%s)
	let "TEST_COUNT++"

	local IS_SUCCESS=$(check_result $RESULT $(($ENDTIME - $STARTTIME)) $TEST_COUNT "Docker network inspect; name = $NAME" $RESULTS_PATH)
	SUCCESS_COUNT=$(($SUCCESS_COUNT + $IS_SUCCESS))
	echo ""
}






function main() {
	if [ -f $RESULTS_PATH ]; then
		rm $RESULTS_PATH
	fi
	touch $RESULTS_PATH

	# Docker ps
	test_ps 

	# Docker inspect - fail when given id DNE
	echo "Docker inspect test1 should fail" >> $RESULTS_PATH
	test_inspect "test1"
	
	# Docker run (create containers)
	test_create_nonet "test1"
	test_ps

	test_create_net "test1"
	test_ps

	# Docker inspect
	test_inspect "test1"

	# Docker stop 
	test_stop "nonet_test1"
	test_ps_a
	test_ps
	
	# Docker rm 
	test_rm "nonet_test1"
	test_ps_a
	test_ps 


	# Start - if resource dne, should error, if already up, nothing, and should start it up if stopped. 
	echo "Docker start nonet_test1 should fail" >> $RESULTS_PATH
	test_start "nonet_test1"

	test_start "test1"
	test_ps_a

	test_stop "test1"
	test_ps_a
	test_start "test1"
	test_ps_a

	# Docker rm - should fail with a running container
	echo "Docker rm test1 should fail; running container" >> $RESULTS_PATH
	test_rm "test1"


	test_stop "test1"
	test_rm "test1"
	test_ps_a


	# Docker rm -f
	test_create_net "test2"
	test_ps_a
	test_rm_f "test2"
	test_ps_a




	# Whoot okay so we can list, create, stop, and remove for now. 

	test_network_ls

	test_network_inspect $try_net_id
	test_network_inspect "default"


	# SUM THE RESULTS
	sum_results $RESULTS_PATH $TEST_COUNT $SUCCESS_COUNT

}

main


# DON'T FORGET TO CLEAN UP
