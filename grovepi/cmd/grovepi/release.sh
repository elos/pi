#!/bin/bash

set -e

make pi
git add ./
git commit -m "releasing build/pi [release.sh]"
