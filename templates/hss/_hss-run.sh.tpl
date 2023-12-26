#!/bin/bash
# Copyright 2019-present Open Networking Foundation
#
# SPDX-License-Identifier: Apache-2.0

set -ex
cp /bin/hss /tmp/coredump/

CONF_DIR="/opt/c3po/hss/conf"
LOGS_DIR="/opt/c3po/hss/logs"
mkdir -p $CONF_DIR $LOGS_DIR

cp /etc/hss/conf/{acl.conf,hss.json,hss.conf,oss.json} $CONF_DIR
cat $CONF_DIR/{hss.json,hss.conf}

cd $CONF_DIR
make_certs.sh hss omec.svc.cluster.local

cd ..
hss -j $CONF_DIR/hss.json