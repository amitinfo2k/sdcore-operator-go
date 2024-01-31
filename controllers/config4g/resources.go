package config4g

import (
	"github.com/amitinfo2k/sdcore-operator-go/api/v1alpha1"
	"github.com/amitinfo2k/sdcore-operator-go/controllers"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/util/intstr"
)

func createDeployment(log logr.Logger, configMapVersion string, config4gDeployment *v1alpha1.PCRFDeployment) (*appsv1.Deployment, error) {
	namespace := config4gDeployment.Namespace
	instanceName := config4gDeployment.Name
	spec := config4gDeployment.Spec

	previleged := true
	runAsUser := int64(0)
	mode := int32(493)

	replicas, resourceRequirements, err := createResourceRequirements(spec)
	if err != nil {
		return nil, err
	}

	/*networkAttachmentDefinitionNetworks, err := createNetworkAttachmentDefinitionNetworks(config4gDeployment.Name, &spec)
	if err != nil {
		return nil, err
	}*/

	podAnnotations := make(map[string]string)
	podAnnotations[controllers.ConfigMapVersionAnnotation] = configMapVersion
	//podAnnotations[controllers.NetworksAnnotation] = networkAttachmentDefinitionNetworks

	initSecurityContext := &apiv1.SecurityContext{
		Privileged: &previleged,
		RunAsUser:  &runAsUser,
	}

	/*securityContext := &apiv1.SecurityContext{
		Capabilities: &apiv1.Capabilities{
			Add: []apiv1.Capability{"NET_ADMIN"},
		},
	}

	*/

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instanceName,
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": instanceName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: podAnnotations,
					Labels: map[string]string{
						"name": instanceName,
					},
				},
				Spec: apiv1.PodSpec{
					InitContainers: []apiv1.Container{
						{
							Name:            "config4g-bootstrap",
							Image:           controllers.PCRFDbImage,
							ImagePullPolicy: apiv1.PullIfNotPresent,
							SecurityContext: initSecurityContext,
							Command:         []string{"sh", "-xc"},
							Args:            []string{"sh /opt/c3po/config4g/config4g-bootstrap.sh"},
							VolumeMounts: []apiv1.VolumeMount{
								{
									MountPath: "/opt/c3po/config4g",
									Name:      "scripts",
								},
							},
						},
					},
					Containers: []apiv1.Container{
						{
							Name:            "config4g",
							Image:           controllers.PCRFImage,
							ImagePullPolicy: apiv1.PullAlways,
							SecurityContext: initSecurityContext,
							Command:         []string{"bash", "-c", "/opt/c3po/config4g/scripts/config4g-run.sh"},

							VolumeMounts: []apiv1.VolumeMount{
								{
									MountPath: "/opt/c3po/config4g/scripts",
									Name:      "scripts",
								},
								{
									MountPath: "/etc/config4g/conf",
									Name:      "configs",
								},
							},
							Resources: *resourceRequirements,
						},
					}, // Containers
					DNSPolicy:     apiv1.DNSClusterFirst,
					RestartPolicy: apiv1.RestartPolicyAlways,
					Volumes: []apiv1.Volume{
						{
							Name: "scripts",
							VolumeSource: apiv1.VolumeSource{
								ConfigMap: &apiv1.ConfigMapVolumeSource{
									LocalObjectReference: apiv1.LocalObjectReference{
										Name: "config4g-scripts",
									},
									DefaultMode: &mode,
								},
							},
						},
						{
							Name: "configs",
							VolumeSource: apiv1.VolumeSource{
								ConfigMap: &apiv1.ConfigMapVolumeSource{
									LocalObjectReference: apiv1.LocalObjectReference{
										Name: "config4g-configs",
									},
									DefaultMode: &mode,
								},
							},
						},
					}, // Volumes
				}, // PodSpec
			}, // PodTemplateSpec
		}, // PodTemplateSpec
	}

	return deployment, nil
}

func createService(config4gDeployment *v1alpha1.PCRFDeployment) *apiv1.Service {
	namespace := config4gDeployment.Namespace
	instanceName := config4gDeployment.Name

	labels := map[string]string{
		"name": instanceName,
	}

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instanceName,
			Namespace: namespace,
		},
		Spec: apiv1.ServiceSpec{
			Selector: labels,
			Ports: []apiv1.ServicePort{
				{
					Name:     "gx",
					Protocol: apiv1.ProtocolTCP,
					Port:     3868,
					NodePort: 31868,
				},
				{
					Name:     "prometheus-exporter",
					Protocol: apiv1.ProtocolTCP,
					Port:     9089,
					NodePort: 30087,
				},
				{
					Name:     "config-port",
					Protocol: apiv1.ProtocolTCP,
					Port:     8080,
					NodePort: 30082,
				},
			},
			Type: apiv1.ServiceTypeNodePort,
		},
	}

	return service
}

func createConfigMap(log logr.Logger, config4gDeployment *v1alpha1.PCRFDeployment) (*apiv1.ConfigMap, error) {
	namespace := config4gDeployment.Namespace
	//instanceName := config4gDeployment.Name
	instanceName := "config4g-configs"
	log.Info("createConfigMap++", "instanceName=", instanceName)

	/*n2ip, err := controllers.GetFirstInterfaceConfigIPv4(config4gDeployment.Spec.Interfaces, "n2")
	if err != nil {
		log.Error(err, "Interface N2 not found in PCRFDeployment Spec")
		return nil, err
	}*/

	templateValues := configurationTemplateValues{
		SVC_NAME:  instanceName,
		S1AP_PORT: 36412,
		S11_PORT:  2123,
	}

	configJson, err := renderConfigFiles(log, templateValues)
	if err != nil {
		log.Error(err, "Could not render PCRF configuration template.")
		return nil, err
	}

	configMap := &apiv1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      instanceName,
		},
		Data: map[string]string{
			"acl.conf":                configJson[0],
			"oss.json":                configJson[1],
			"config4g.json":               configJson[2],
			"config4g.conf":               configJson[3],
			"subscriber_mapping.json": configJson[4],
		},
	}
	log.Info("createConfigMap--")
	return configMap, nil
}

func createScriptConfigMap(log logr.Logger, config4gDeployment *v1alpha1.PCRFDeployment) (*apiv1.ConfigMap, error) {
	namespace := config4gDeployment.Namespace
	instanceName := "config4g-scripts"
	log.Info("createScriptConfigMap++", "instanceName=", instanceName)

	/*n2ip, err := controllers.GetFirstInterfaceConfigIPv4(config4gDeployment.Spec.Interfaces, "n2")
	if err != nil {
		log.Error(err, "Interface N2 not found in PCRFDeployment Spec")
		return nil, err
	}*/

	templateValues := configurationTemplateValues{
		SVC_NAME:  instanceName,
		S1AP_PORT: 36412,
		S11_PORT:  2123,
	}

	config4gScriptsConfig, err := renderScriptFiles(log, templateValues)
	if err != nil {
		log.Error(err, "Could not render PCRF Scripts configuration template.")
		return nil, err
	}

	configMap := &apiv1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      instanceName,
		},
		Data: map[string]string{
			"config4g-bootstrap.sh": config4gScriptsConfig[0],
			"config4g-run.sh":       config4gScriptsConfig[1],
		},
	}
	log.Info("createScriptConfigMap--")
	return configMap, nil
}

func createResourceRequirements(config4gDeploymentSpec v1alpha1.PCRFDeploymentSpec) (int32, *apiv1.ResourceRequirements, error) {
	// TODO: Requirements should be calculated based on DL, UL
	// TODO: increase number of recpicas based on NFDeployment.Capacity.MaxSessions

	var replicas int32 = 1
	//downlink := resource.MustParse("5G")
	//uplink := resource.MustParse("1G")
	var cpuLimit string
	var cpuRequest string
	var memoryLimit string
	var memoryRequest string

	if config4gDeploymentSpec.Capacity.MaxSubscribers > 1000 {
		cpuLimit = "300m"
		memoryLimit = "256Mi"
		cpuRequest = "300m"
		memoryRequest = "256Mi"
	} else {
		cpuLimit = "150m"
		memoryLimit = "128Mi"
		cpuRequest = "150m"
		memoryRequest = "128Mi"
	}

	resources := apiv1.ResourceRequirements{
		Limits: apiv1.ResourceList{
			apiv1.ResourceCPU:    resource.MustParse(cpuLimit),
			apiv1.ResourceMemory: resource.MustParse(memoryLimit),
		},
		Requests: apiv1.ResourceList{
			apiv1.ResourceCPU:    resource.MustParse(cpuRequest),
			apiv1.ResourceMemory: resource.MustParse(memoryRequest),
		},
	}

	return replicas, &resources, nil
}

/*func createNetworkAttachmentDefinitionNetworks(templateName string, config4gDeploymentSpec *v1alpha1.PCRFDeploymentSpec) (string, error) {
	return controllers.CreateNetworkAttachmentDefinitionNetworks(templateName, map[string][]nephiov1alpha1.InterfaceConfig{
		"n2": controllers.GetInterfaceConfigs(config4gDeploymentSpec.Interfaces, "n2"),
	})
}*/
