#!/bin/sh

# Copyright 2020-present Open Networking Foundation
#
# SPDX-License-Identifier: Apache-2.0

set -xe

cd /free5gc

cat config/webuicfg.conf

./webconsole/webconsole -webuicfg config/webuicfg.conf

