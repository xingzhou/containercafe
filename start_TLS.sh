#!/bin/bash

# Help menu
helpme()
{
	cat <<HELPMEHELPME

Syntax: ${0} <Apikey> <Space_id> 
Where:
	Apikey = Apikey for this user
	Space_id = Id of the desired space

HELPMEHELPME
}

# Check args
if [[ "$1" == "-?" || "$1" == "-h" || "$1" == "--help" || "$1" == "help" ]] ; then
  helpme
  exit 1
elif [[ "$1" == "" || "$2" == "" ]] ; then
	echo "Incorrect Arguments"
	helpme
	exit 1
fi

API_KEY=$1
SPACE_ID=$2


# Create certificate
# WHAT ABOUT LOCATION OF CA / PASSWORD rn just my configuration
echo "Creating certificates"
./mk_user_cert.sh "$API_KEY"

# Right now, doing only 1 space per user
if [ $? -eq 1 ]; then
	echo "This user already has credentials"
	exit 1
fi


echo "Writing to creds.json"
# Now add to creds.json file
if [ ! -f creds.json ]; then
	echo "Creating creds.json file"
	touch creds.json
fi
CREDS_STRING_1='{"Status":200, "Node":"10.140.171.205:443", "Docker_id":"", "Container":"", "Swarm_shard":true, "Tls_override":true, "Space_id":"myspace", "Reg_namespace":"swarm", "Apikey":'
QUOTE='"'
CREDS_STRING_2=', "Orguuid":"orgname", "Userid":"userid", "Endpoint_type":"radiant", "TLS_path": '
TLS_path_endpt="/user_certificates/"
TLS_path=$PWD$TLS_path_endpt$API_KEY
CLOSE="}"

ECHO_STRING=$CREDS_STRING_1$QUOTE$API_KEY$QUOTE$CREDS_STRING_2$QUOTE$TLS_path$QUOTE$CLOSE
echo $ECHO_STRING >> creds.json  
