#!/bin/bash
cmd=$1
http_server=$2
http_port=$3
running="import BaseHTTPServer as bhs, SimpleHTTPServer as shs; bhs.HTTPServer((\"$http_server\", $http_port ), shs.SimpleHTTPRequestHandler).serve_forever()"
pid=""
startHTTP ()
{
  python -c "$running" &
  if [ $? -eq 0 ]
  then
    pid=$!
    echo $pid > $(pwd)/pid
    echo "Successfully started the HTTP server"
    exit 0
  else
    echo "Could not start the HTTP server" >&2
    exit 1
  fi
}

stopHTTP ()
{
  if [ [-s $(pwd)/pid] ]
  then
    echo "PID was not found, exiting..."
    exit 1
  fi
  kill -9 $(cat $(pwd)/pid)
  if [ $? -eq 0 ]
  then
    echo "Successfully stoped the HTTP server"
    echo "" > $(pwd)/pid
    exit 0
  elif [ $(pidof python 2>/dev/null) ]
  then
    kill -9 $(pidof python 2>/dev/null)
    echo "Could not stop gacefully the HTTP server" >&2
    exit 1
  else
    echo "Could not stop the HTTP server" >&2
    exit 1
  fi
}

monitorHTTP ()
{
  response=$(curl --write-out %{http_code} --silent --output /dev/null http://$http_server:$http_port)
  if [ $response -eq 200 ]
  then
    echo "Successfully monitored the HTTP server"
    exit 0
  else
    echo "Did not get a 200 response form the HTTP server" >&2
    exit 1
  fi

}

case $cmd in

  start)    startHTTP;;

  stop)     stopHTTP;;

  monitor)  monitorHTTP;;

  *)        echo "no command is specified, please use: start | stop | monitor";;
esac
