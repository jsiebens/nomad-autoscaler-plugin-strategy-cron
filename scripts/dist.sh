#!/bin/bash

suffix=""
if [ $1 == "windows" ]; then
  suffix=".exe"
fi

CGO_ENABLED=0 GOOS=$1 GOARCH=$2 go build -ldflags "-s -w" -a -installsuffix cgo -o "dist/cron${suffix}"
zip -j dist/cron_$1_$2.zip "dist/cron${suffix}"
rm -rf "dist/cron${suffix}"