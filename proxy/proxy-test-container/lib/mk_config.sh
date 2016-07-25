#!/bin/bash 

API_KEY="$1"
FILE_PATH="$2""/config.json"

if [ -f $FILE_PATH ]; then
	rm $FILE_PATH
fi
touch $FILE_PATH


echo "{" >> $FILE_PATH
echo "	\"HttpHeaders\": {" >> $FILE_PATH
echo "		\"X-Tls-Client-Dn\": \"/CN=$API_KEY\"" >> $FILE_PATH
echo "	}" >> $FILE_PATH
echo "}" >> $FILE_PATH



# {
#     "HttpHeaders": {
#           "X-Tls-Client-Dn": "/CN=CWch4JDD9GcTmog2w8zxUd0UyYMb6L2zEevN5RAowJnzyC4c"
#     }
# }