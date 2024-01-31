#!/bin/sh
# Copyright 2021-present Open Networking Foundation
#
# SPDX-License-Identifier: Apache-2.0

while ! curl -f --connect-timeout 5 http://spgwc.omec.svc:8080/startup
do
  echo Waiting for SPGWC to be ready
  nslookup spgwc.omec.svc
  sleep 5
done
echo SPGWC is ready

echo Posting to sync URL 
while ! curl --connect-timeout 5 -f -X POST 
do
  echo Failed posting to sync URL
  sleep 5
done
echo

echo Sleeping forever
while true
do
  sleep 86400
done
