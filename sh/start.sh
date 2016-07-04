#!/bin/bash

set -e

./stop.sh
go run sense.go > stdout.txt 2> stderr.txt &
echo $! > running.pid
echo "Started"
