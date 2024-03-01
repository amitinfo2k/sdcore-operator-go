package mme

import (
	"bytes"
	"text/template"

	"github.com/go-logr/logr"
)

type configurationTemplateValues struct {
	SVC_NAME  string
	S1AP_PORT int
	S11_PORT  int
}

func renderConfigFiles(log logr.Logger, values configurationTemplateValues) ([]string, error) {
	var buffer bytes.Buffer
	var theArray []string

	xfiles := [2]string{"templates/mme/config.json", "templates/mme/s6a_fd.conf"}
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
	xfiles := [2]string{"templates/mme/_mme-init.sh.tpl", "templates/mme/_mme-run.sh.tpl"}

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
