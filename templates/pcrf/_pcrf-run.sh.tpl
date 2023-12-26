#!/bin/bash
# Copyright 2020-present Open Networking Foundation
#
# SPDX-License-Identifier: Apache-2.0

set -ex
cp /bin/pcrf /tmp/coredump/

CONF_DIR="/opt/c3po/pcrf/conf"
LOGS_DIR="/opt/c3po/pcrf/logs"
#TODO - Need to remove logs directory
mkdir -p $CONF_DIR $LOGS_DIR

cp /etc/pcrf/conf/{acl.conf,pcrf.json,pcrf.conf,oss.json,subscriber_mapping.json} $CONF_DIR
cat $CONF_DIR/{pcrf.json,pcrf.conf}

cd $CONF_DIR
make_certs.sh pcrf omec.svc.cluster.local

cd ..
pcrf -j $CONF_DIR/pcrf.json