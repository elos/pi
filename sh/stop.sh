#!/bin/bash

set -e

if [ -e /tmp/grovepi/pid ]; then
	echo "Killing old run"
	sudo kill $(cat /tmp/grovepi/pid);
	rm running.pid;
else
	echo "None running"
fi
