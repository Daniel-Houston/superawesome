#!/bin/sh
NAME="superawesome"
USAGE="usage: $0 <remote-host>"
if [ "$#" -ne 1  ]; then
	echo "You must specify the remote host to which to deploy $NAME"
	echo $USAGE
	exit 1
fi

if [ "$1" = "-h" ]; then
	echo $USAGE
	exit 1
fi

HOST=$1

echo "Building $NAME"
env GOOS=linux GOARCH=arm GOARM=6 go build -o $NAME 
echo "Deploying $NAME"
scp $NAME $HOST:
