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


# Create certificate




