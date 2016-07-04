#!/bin/bash

set -e

./stop.sh
go run ./grovepi/cmd/grovepi/main.go -config /tmp/grovepi.config > /var/grovepi/stdout.txt 2> /var/grovepi/stderr.txt &
echo $! > /var/grovepi/pid
echo "Started"
