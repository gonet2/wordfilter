#!/bin/bash
set -e
NSQD_HOST="http://172.17.42.1:4151"
case $1 in 
	production)
		NSQD_HOST="http://172.17.42.1:4151"
		;;
	testing)
		NSQD_HOST="http://172.17.42.1:4151"
		;;
esac
export NSQD_HOST=$NSQD_HOST
echo "NSQD_HOST set to:" $NSQD_HOST
$GOBIN/wordfilter
