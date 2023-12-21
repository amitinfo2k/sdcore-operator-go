package mme

import (
	"bytes"
	"text/template"

	"github.com/go-logr/logr"
)

const configJsonTemplateSource = `
  {
    "mme": {
      "apnlist": {
        "default": "spgwc",
        "internet": "spgwc"
      },
      "code": 1,
      "feature_list": {
        "dcnr_support": "disabled"
      },
      "group_id": 1,
      "logging": "debug",
      "name": "vmmestandalone",
      "plmnlist": {
        "plmn1": "mcc=312,mnc=440",
        "plmn2": "mcc=208,mnc=01",
        "plmn3": "mcc=315,mnc=010"
      },
      "prom_port": 3081,
      "security": {
        "int_alg_list": "[EIA1, EIA2, EIA0]",
        "sec_alg_list": "[EEA0, EEA1, EEA2]"
      }
    },
    "s11": {
      "egtp_default_port": 2123
    },
    "s1ap": {
      "sctp_port": 36412,
      "sctp_port_external": 32413
    },
    "s6a": {
      "host": "hss.omec.svc.cluster.local",
      "host_type": "freediameter",
      "realm": "omec.svc.cluster.local"
    }
  }
`

const s6aFdJsonTemplateSource = `
  Identity = "mme.omec.svc.cluster.local";
  Realm = "omec.svc.cluster.local";
  TLS_Cred = "conf/mme.cert.pem",
            "conf/mme.key.pem";
  TLS_CA = "conf/cacert.pem";
  AppServThreads = 40;
  SCTP_streams = 3;
  NoRelay;
  No_IPv6;
  #Port = 3868;
  #SecPort = 3869;

  ConnectPeer = "hss.omec.svc.cluster.local" { No_TLS; port = 3868; };

  LoadExtension = "/usr/local/lib/freeDiameter/dict_3gpp2_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_draftload_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_etsi283034_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_rfc4004_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_rfc4006bis_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_rfc4072_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_rfc4590_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_rfc5447_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_rfc5580_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_rfc5777_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_rfc5778_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_rfc6734_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_rfc6942_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_rfc7155_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_rfc7683_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_rfc7944_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_ts29061_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_ts29128_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_ts29154_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_ts29173_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_ts29212_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_ts29214_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_ts29215_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_ts29217_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_ts29229_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_ts29272_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_ts29273_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_ts29329_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_ts29336_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_ts29337_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_ts29338_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_ts29343_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_ts29344_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_ts29345_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_ts29368_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_ts29468_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_ts32299_avps.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_S6as6d.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_S6c.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_S6t.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_SGd.fdx";
  LoadExtension = "/usr/local/lib/freeDiameter/dict_T6aT6bT7.fdx";
`
const initScriptTemplateSource = `
  #!/bin/sh

  # Copyright 2019-present Open Networking Foundation
  #
  # SPDX-License-Identifier: Apache-2.0

  set -ex

  cp /opt/mme/config/config.json /opt/mme/config/shared/config.json
  cd /opt/mme/config/shared

  # Set local IP address for s1ap and s11 networks to the config
  jq --arg MME_LOCAL_IP "$POD_IP" '.mme.ip_addr=$MME_LOCAL_IP' config.json > config.tmp && mv config.tmp config.json
  jq --arg MME_LOCAL_IP "$POD_IP" '.s1ap.s1ap_local_addr=$MME_LOCAL_IP' config.json > config.tmp && mv config.tmp config.json
  jq --arg MME_LOCAL_IP "$POD_IP" '.s11.egtp_local_addr=$MME_LOCAL_IP' config.json > config.tmp && mv config.tmp config.json

  # Set SPGWC address to the config
  # We need to convert service domain name to actual IP address
  # because mme apps does not take domain address - should be fixed in openmme
  # TODO 
  SPGWC_ADDR=10.15.10.11
  jq --arg SPGWC_ADDR "$SPGWC_ADDR" '.s11.sgw_addr //= $SPGWC_ADDR' config.json > config.tmp && mv config.tmp config.json
  jq --arg SPGWC_ADDR "$SPGWC_ADDR" '.s11.pgw_addr //= $SPGWC_ADDR' config.json > config.tmp && mv config.tmp config.json

  # Add additional redundant keys - should be fixed in openmme
  HSS_TYPE=$(jq -r '.s6a.host_type' config.json)
  HSS_HOST=$(jq -r '.s6a.host' config.json)
  jq --arg HSS_TYPE "$HSS_TYPE" '.s6a.hss_type=$HSS_TYPE' config.json > config.tmp && mv config.tmp config.json
  jq --arg HSS_HOST "$HSS_HOST" '.s6a.host_name=$HSS_HOST' config.json > config.tmp && mv config.tmp config.json

  # Copy the final configs for each applications
  cp /opt/mme/config/shared/config.json /opt/mme/config/shared/mme.json
  cp /opt/mme/config/shared/config.json /opt/mme/config/shared/s11.json
  cp /opt/mme/config/shared/config.json /opt/mme/config/shared/s1ap.json
  cp /opt/mme/config/shared/config.json /opt/mme/config/shared/s6a.json
  cp /opt/mme/config/s6a_fd.conf /opt/mme/config/shared/s6a_fd.conf

  #This multiple copies of config needs some cleanup. For now I want 
  #that after running mme_init config to be present in the target directory
  cp /opt/mme/config/shared/* /openmme/target/conf/

  # Generate certs
  MME_IDENTITY=mme.omec.svc.cluster.local
  DIAMETER_HOST=$(echo $MME_IDENTITY | cut -d'.' -f1)
  DIAMETER_REALM=omec.svc.cluster.local

  cp /openmme/target/conf/make_certs.sh /opt/mme/config/shared/make_certs.sh
  cd /opt/mme/config/shared
  ./make_certs.sh $DIAMETER_HOST $DIAMETER_REALM
`

const runScriptTemplateSource = `
  #!/bin/bash

  # Copyright 2019-present Open Networking Foundation
  #
  # SPDX-License-Identifier: Apache-2.0
  runScriptTemplateSource
  APPLICATION=$1

  # copy config files to openmme target directly
  cp /opt/mme/config/shared/* /openmme/target/conf/

  cd /openmme/target
  export LD_LIBRARY_PATH=/usr/local/lib:./lib

  case $APPLICATION in
      "mme-app")
        echo "Starting mme-app"
        echo "conf/mme.json"
        cat conf/mme.json
        ./bin/mme-app
        ;;
      "s1ap-app")
        echo "Starting s1ap-app"
        echo "conf/s1ap.json"
        cat conf/s1ap.json
        ./bin/s1ap-app
        ;;
      "s6a-app")
        echo "Starting s6a-app"
        echo "conf/s6a.json"
        cat conf/s6a.json
        echo "conf/s6a_fd.conf"
        cat conf/s6a_fd.conf
        ./bin/s6a-app
        ;;
      "s11-app")
        echo "Starting s11-app"
        echo "conf/s11.json"
        cat conf/s11.json
        ./bin/s11-app
        ;;
      *)
        echo "invalid app $APPLICATION"
        ;;
  esac
`

var intScriptTemplate = template.Must(template.New("MMEScripts").Parse(initScriptTemplateSource))
var runScriptTemplate = template.Must(template.New("MMERunScript").Parse(runScriptTemplateSource))
var configJsonTemplate = template.Must(template.New("MMEConfigJson").Parse(configJsonTemplateSource))
var s6aFdJsonTemplate = template.Must(template.New("MMES6A").Parse(s6aFdJsonTemplateSource))

type configurationTemplateValues struct {
	SVC_NAME  string
	S1AP_PORT int
	S11_PORT  int
}

func renderConfigJsonTemplate(values configurationTemplateValues) (string, error) {
	var buffer bytes.Buffer
	if err := configJsonTemplate.Execute(&buffer, values); err == nil {
		return buffer.String(), nil
	} else {
		return "", err
	}
} //

func renderS6AFDJsonTemplate(values configurationTemplateValues) (string, error) {
	var buffer bytes.Buffer
	if err := s6aFdJsonTemplate.Execute(&buffer, values); err == nil {
		return buffer.String(), nil
	} else {
		return "", err
	}
}

/*func renderConfigFiles(log logr.Logger, values configurationTemplateValues) ([]string, error) {
	var buffer bytes.Buffer
	var theArray []string
	xfiles := [2]string{"templates/_mme-init.sh.tpl", "templates/_mme-run.sh.tpl"}

	for i, v := range xfiles {
		log.Info("renderConfigFiles++", "i=", i, "v=", v)
		configTemplate, err := template.ParseFiles(v)
		if err == nil {
			if err := configTemplate.Execute(&buffer, values); err == nil {
				theArray[i] = buffer.String()
			} else {
				log.Error(err, "Error while rendering template")
				return nil, err
			}
		} else {
			log.Error(err, "Error while reading template")
			return nil, err
		}
	}
	return theArray, nil
}*/

func renderConfigFiles(log logr.Logger, values configurationTemplateValues) ([]string, error) {
	var buffer bytes.Buffer
	theArray := []string{}

	if err := intScriptTemplate.Execute(&buffer, values); err == nil {
		//theArray[0] = buffer.String()
		theArray = append(theArray, buffer.String())
	} else {
		log.Error(err, "Error while rendering template")
		return nil, err
	}
	buffer.Reset()
	if err := runScriptTemplate.Execute(&buffer, values); err == nil {
		theArray = append(theArray, buffer.String())
	} else {
		log.Error(err, "Error while rendering template")
		return nil, err
	}

	return theArray, nil
}
