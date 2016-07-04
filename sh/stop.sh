#!/bin/bash

set -e

if [ -e running.pid ]; then
	echo "Killing old run"
	sudo kill $(cat running.pid);
	rm running.pid;
else
	echo "None running"
fi
