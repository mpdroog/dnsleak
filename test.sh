#!/bin/bash
#set -e
set -x
set -u

ping random1.dnstest.spyoff.com -c 2 -t 2
curl -H "Content-Type: application/json" -X POST -d '{"domain":["random1.dnstest.spyoff.com","random2.dnstest.spyoff.com"]}' http://ns-dnstest.spyoff.com
