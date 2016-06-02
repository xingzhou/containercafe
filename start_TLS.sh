#!/bin/bash

# Help menu
helpme()
{
	cat <<HELPMEHELPME

<<<<<<< 7d176e9d2cf10ab0225e38a5974550a586dd378c
Syntax: ${0} <Space_id> 
Where:
=======
Syntax: ${0} <Apikey> <Space_id> 
Where:
	Apikey = Apikey for this user
>>>>>>> Clean go code, and written TLS script
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


API_KEY_LEN=48
SPACE_ID=$1


# Generate API key
generate_api_key() {
	echo "Generating API key"
	API_KEY=$(cat /dev/urandom | env LC_CTYPE=C tr -dc 'a-zA-Z0-9' | fold -w $API_KEY_LEN | head -n 1)

}
generate_api_key
echo "Generated API key: "$API_KEY


# Create certificate
# WHAT ABOUT LOCATION OF CA / PASSWORD rn just my configuration
echo "Creating certificates"
./mk_user_cert.sh "$API_KEY"

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


echo "{\"Status\":200, \"Node\":\"10.140.171.205:443\", \"Docker_id\":\"\", \"Container\": \"\", \"Swarm_shard\":true, \"Tls_override\":true, \"Space_id\":\"$SPACE_ID\", \"Reg_namespace\":\"swarm\", \"Apikey\":\"$API_KEY\", \"Orguuid\":\"orgname\", \"Userid\":\"userid\", \"Endpoint_type\":\"radiant\", \"TLS_path\":\"$PWD/user_certificates/$API_KEY\"}" >> creds.json

echo "Certificates created for Apikey "$API_KEY
echo "Located at "$PWD"/user_certificates/"$API_KEY
