#!/bin/bash
if [ ! -f $1 ]; then
	echo "Usage: $0 <executable to deploy>"
	exit 1
fi
sudo mv $1 /srv/
sudo pushd /srv/
