#!/bin/bash
# Copyright 2020-present Open Networking Foundation
#
# SPDX-License-Identifier: Apache-2.0

set -ex

until cqlsh --file /opt/c3po/pcrfdb/pcrf_cassandra.cql cassandra;
    do echo "Provisioning PCRFDB";
    sleep 2;
done