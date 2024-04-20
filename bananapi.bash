#!/bin/bash
echo 'Compiling for bananapi'
GOOS=linux GOARCH=arm  go build
if [ $? -ne 0 ]; then
	echo 'An error has occurred! Aborting the script execution...'
	exit 1
fi
echo 'Copy  to banana'
scp  rgpio root@192.168.88.20:/home/rura
