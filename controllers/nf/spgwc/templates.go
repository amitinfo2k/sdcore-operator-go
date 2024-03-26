package spgwc

import (
	"bytes"
	"text/template"

	"github.com/go-logr/logr"
)

type configurationTemplateValues struct {
	N4_IP     string
	S1AP_PORT int
	S11_PORT  int
}

func renderConfigFiles(log logr.Logger, values configurationTemplateValues) ([]string, error) {
	var buffer bytes.Buffer
	var theArray []string

	xfiles := [4]string{"templates/spgwc/config.json", "templates/spgwc/cp.json", "templates/spgwc/gx.conf", "templates/spgwc/subscriber_mapping.json"}
	for i, v := range xfiles {
		log.Info("renderConfigFiles++", "i=", i, "v=", v)
		configTemplate, err := template.ParseFiles(v)
		if err == nil {
			if err := configTemplate.Execute(&buffer, values); err == nil {
				theArray = append(theArray, buffer.String())
				buffer.Reset()
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
}

func renderScriptFiles(log logr.Logger, values configurationTemplateValues) ([]string, error) {
	var buffer bytes.Buffer
	var theArray []string
	xfiles := [2]string{"templates/spgwc/_spgwc-init.sh.tpl", "templates/spgwc/_spgwc-run.sh.tpl"}

	for i, v := range xfiles {
		log.Info("renderScriptFiles++", "i=", i, "v=", v)
		configTemplate, err := template.ParseFiles(v)
		if err == nil {
			if err := configTemplate.Execute(&buffer, values); err == nil {
				theArray = append(theArray, buffer.String())
				buffer.Reset()
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
}
