#!/bin/bash
# Copyright 2019-present Open Networking Foundation
#
# SPDX-License-Identifier: Apache-2.0

APPLICATION=$1
set -xe

mkdir -p /opt/cp/config
cd /opt/cp/config
cp /etc/cp/config/{*.json,*.conf} .

#echo "172.29.4.68 upf" >> /etc/hosts
#echo "172.29.4.70 upf2" >> /etc/hosts

#Single Server
#echo "172.29.4.17 upf3" >> /etc/hosts

#Tamu
#echo "172.29.8.3 upf-campus01-1" >> /etc/hosts

case $APPLICATION in
    "ngic_controlplane")
      echo "Starting ngic controlplane app"
      cat /opt/cp/config/cp.json
      cat /opt/cp/config/subscriber_mapping.json
      cp /bin/ngic_controlplane /tmp/coredump/

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
      cp /bin/gx_app /tmp/coredump/
      cd /opt/cp/
      gx_app
      ;;

    *)
      echo "invalid app $APPLICATION"
      ;;
esac