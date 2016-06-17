#!/bin/bash

function generic_docker_command() {
	local CMD=$1 
	eval "DOCKER_TLS_VERIFY=\"\" docker "$CMD
}

function check_result() {
	RESULT=$1
	TIME=$2
	TEST_COUNT=$3
	SUMMARY="$4" # Double quotes - literal value of all the characters? 
	RESULTS_PATH="$5"

	if [ $TIME -eq 0 ]; then 
		TIME=1 # Wat I don't know why
	fi 

	# I don't know what this lock stuff is, so ignoring...


	# WHAT AM I DOING

	# Process result
	if [ $RESULT -eq 0 ]; then
		# WHAT AM I DOING
		#log "Test SUCCESS: $SUMMARY in $TIME seconds."
		echo "SUCCESS Test $TEST_COUNT - $SUMMARY in $TIME seconds." >> $RESULTS_PATH
		echo "" >> $RESULTS_PATH

		# Pass back success
		echo 1

	else 
		# What is this complicated stuff in the original
		#log "Test FAILURE: $SUMMARY in $TIME seconds."
		echo "FAIL Test $TEST_COUNT - $SUMMARY in $TIME seconds." >> $RESULTS_PATH
		echo "" >> $RESULTS_PATH
		# Pass back failure
		echo 0
	fi 

}


# At end - summation of results
function sum_results() {
	RESULTS_PATH=$1
	TOTAL_COUNT=$2
	PASS_COUNT=$3

	echo "" >> $RESULTS_PATH
	printf '%*s\n' "${COLUMNS:-$(tput cols)}" '' | tr ' ' - >> $RESULTS_PATH
	echo "" >> $RESULTS_PATH

	FAIL_COUNT=$((TOTAL_COUNT - PASS_COUNT))
	echo "Test summary:" >> $RESULTS_PATH
	echo "Total = $TOTAL_COUNT tests" >> $RESULTS_PATH
	echo "Passed = $PASS_COUNT tests" >> $RESULTS_PATH
	echo "Failed = $FAIL_COUNT tests" >> $RESULTS_PATH

}
