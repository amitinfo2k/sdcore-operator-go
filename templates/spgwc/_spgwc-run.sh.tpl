#!/bin/bash
# Copyright 2019-present Open Networking Foundation
#
# SPDX-License-Identifier: Apache-2.0

APPLICATION=$1
set -xe

mkdir -p /opt/cp/config
cd /opt/cp/config
cp /etc/cp/config/{*.json,*.conf} .

echo "10.17.0.2 upf1" >> /etc/hosts

case $APPLICATION in
    "ngic_controlplane")
      echo "Starting ngic controlplane app"
      cat /opt/cp/config/cp.json
      cat /opt/cp/config/subscriber_mapping.json
      ngic_controlplane -f /etc/cp/config/
      ;;

    "gx-app")
      echo "Starting gx-app"
      SPGWC_IDENTITY="spgwc.omec.svc.cluster.local";
      DIAMETER_HOST=$(echo $SPGWC_IDENTITY| cut -d'.' -f1)
      DIAMETER_REALM="omec.svc.cluster.local";
      chmod +x /bin/make_certs.sh
      cp /bin/make_certs.sh /opt/cp/config
      /bin/make_certs.sh $DIAMETER_HOST $DIAMETER_REALM
      cd /opt/cp/
      gx_app
      ;;

    *)
      echo "invalid app $APPLICATION"
      ;;
esac
