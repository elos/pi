#!/bin/bash

set -e

./stop.sh
go run sense.go > /var/grovepi/stdout.txt 2> /var/grovepi/stderr.txt &
echo $! > /var/grovepi/pid
echo "Started"
