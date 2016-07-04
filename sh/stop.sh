#!/bin/bash

set -e

if [ -e /var/grovepi/pid ]; then
	echo "Killing old run"
	sudo kill $(cat /var/grovepi/pid);
	rm running.pid;
else
	echo "None running"
fi
