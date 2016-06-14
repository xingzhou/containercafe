#!/bin/bash

# Help menu
helpme()
{
	cat <<HELPMEHELPME

Syntax: ${0} <Space_id> 
Where:
	Space_id = Id of the desired space

HELPMEHELPME
}

# Check args
if [[ "$1" == "-?" || "$1" == "-h" || "$1" == "--help" || "$1" == "help" ]] ; then
  helpme
  exit 1
elif [[ "$1" == "" ]] ; then
	echo "Incorrect Arguments"
	helpme
	exit 1
fi


# Variables
API_KEY_LEN=48
SPACE_ID=$1

STATUS=200
TARGET_SERVER="10.140.146.7"
DOCKER_ID=""
CONTAINER=""
SWARM_SHARD=true
TLS_OVERRIDE=true
REG_NAMESPACE="swarm"
ORGUUID="orgname"
USERID="userid"
ENDPOINT_TYPE="radiant"


# Generate API key
generate_api_key() {
	echo "Generating API key of length "$API_KEY_LEN
	API_KEY=$(cat /dev/urandom | env LC_CTYPE=C tr -dc 'a-zA-Z0-9' | fold -w $API_KEY_LEN | head -n 1)

}
generate_api_key
echo "Generated API key: "$API_KEY


# Create certificate
# WHAT ABOUT LOCATION OF CA / PASSWORD rn just my configuration
echo "Creating certificates"
TLS_dir=$(./mk_user_cert.sh "$API_KEY")

# Right now, doing only 1 space per user
if [ $? -eq 1 ]; then
	echo "This Apikey already has credentials. Process terminating."
	exit 1
fi


echo "Writing to creds.json"
# Now add to creds.json file
if [ ! -f creds.json ]; then
	echo "Creating creds.json file"
	touch creds.json
fi



echo "{\"Status\":$STATUS, \"Node\":\"$TARGET_SERVER\", \"Docker_id\":\"$DOCKER_ID\", \"Container\":\"$CONTAINER\", \"Swarm_shard\":$SWARM_SHARD, \"Tls_override\":$TLS_OVERRIDE, \"Space_id\":\"$SPACE_ID\", \"Reg_namespace\":\"$REG_NAMESPACE\", \"Apikey\":\"$API_KEY\", \"Orguuid\":\"$ORGUUID\", \"Userid\":\"$USERID\", \"Endpoint_type\":\"$ENDPOINT_TYPE\", \"TLS_path\":\"$TLS_dir\"}" >> creds.json

echo "Certificates created for Apikey "$API_KEY
echo "Located at "$TLS_dir
