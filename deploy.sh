#!/bin/bash
CONFIG="config.json"
if [ ! -f $1 ]; then
	echo "Usage: $0 <executable to deploy>"
	exit 1
fi
sudo mv $1 /srv/
sudo mv $CONFIG /srv/
sudo pushd /srv/
./$1 & 
disown
