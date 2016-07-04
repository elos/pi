#!/bin/bash

set -e

mkdir -p /tmp/grovepi

./stop.sh
go run ../grovepi/cmd/grovepi/main.go -config /tmp/grovepi/config > /tmp/grovepi/stdout.txt 2> /tmp/grovepi/stderr.txt &
echo $! > /tmp/grovepi/pid
echo "Started"
cat /tmp/grovepi/stderr.txt
tail -f /tmp/grovepi/stdout.txt
