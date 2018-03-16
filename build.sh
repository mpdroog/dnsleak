#!/bin/bash
set -e
set -x
set -u

env GOOS=linux GOARCH=amd64 go build
echo -n $'\003' | dd bs=1 count=1 seek=7 conv=notrunc of=./dnsleak

scp dnsleak root@ns-dnstest.spyoff.com:~
