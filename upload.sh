#!/bin/sh
BINARY="superawesome"
DEPLOY="deploy.sh"
CONFIG="config.json"
USAGE="usage: $0 <remote-host>"
if [ "$#" -ne 1  ]; then
	echo "You must specify the remote host to which to deploy $BINARY"
	echo $USAGE
	exit 1
fi

if [ "$1" = "-h" ]; then
	echo $USAGE
	exit 1
fi
if [ "$1" = "superawesome.host" ]; then
	echo "You entered $1, did you mean to enter superawesome.home?"
	exit 1
fi

HOST=$1

echo "Building $BINARY"
env GOOS=linux GOARCH=arm GOARM=6 go build -o $BINARY 
echo "Uploading $BINARY and $DEPLOY"
scp $BINARY $HOST: 
scp $DEPLOY $HOST:
scp $CONFIG $HOST:
