package hss

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

func createDeployment(log logr.Logger, configMapVersion string, hssDeployment *v1alpha1.HSSDeployment) (*appsv1.Deployment, error) {
	namespace := hssDeployment.Namespace
	instanceName := hssDeployment.Name
	spec := hssDeployment.Spec

	previleged := true
	runAsUser := int64(0)
	mode := int32(493)
	configMode := int32(420)

	replicas, resourceRequirements, err := createResourceRequirements(spec)
	if err != nil {
		return nil, err
	}

	/*networkAttachmentDefinitionNetworks, err := createNetworkAttachmentDefinitionNetworks(hssDeployment.Name, &spec)
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
							Name:            "hss-bootstrap",
							Image:           controllers.HssDbImage,
							ImagePullPolicy: apiv1.PullIfNotPresent,
							SecurityContext: initSecurityContext,
							Command:         []string{"sh", "-xc"},
							Args:            []string{"sh /opt/c3po/hss/scripts/hss-bootstrap.sh"},
							VolumeMounts: []apiv1.VolumeMount{
								{
									MountPath: "/opt/c3po/hss/scripts",
									Name:      "scripts",
								},
							},
						},
					},
					Containers: []apiv1.Container{
						{
							Name:            "hss",
							Image:           controllers.HSSImage,
							ImagePullPolicy: apiv1.PullAlways,

							Command: []string{"bash", "-c", "/opt/c3po/hss/scripts/hss-run.sh; sleep 3600"},

							Env: []apiv1.EnvVar{
								{
									Name: "POD_IP",
									ValueFrom: &apiv1.EnvVarSource{
										FieldRef: &apiv1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									MountPath: "/opt/c3po/hss/scripts",
									Name:      "scripts",
								},
								{
									MountPath: "/etc/hss/conf",
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
										Name: "hss-scripts",
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
										Name: "hss-configs",
									},
									DefaultMode: &configMode,
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

func createService(hssDeployment *v1alpha1.HSSDeployment) *apiv1.Service {
	namespace := hssDeployment.Namespace
	instanceName := hssDeployment.Name

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
					Name:     "s6a",
					Protocol: apiv1.ProtocolTCP,
					Port:     3868,
					NodePort: 31868,
				},
				{
					Name:     "config-port",
					Protocol: apiv1.ProtocolTCP,
					Port:     8080,
					NodePort: 30081,
				},
				{
					Name:     "prometheus-exporter",
					Protocol: apiv1.ProtocolTCP,
					Port:     9089,
					NodePort: 30086,
				},
			},
			Type: apiv1.ServiceTypeNodePort,
		},
	}

	return service
}

func createConfigMap(log logr.Logger, hssDeployment *v1alpha1.HSSDeployment) (*apiv1.ConfigMap, error) {
	namespace := hssDeployment.Namespace
	//instanceName := hssDeployment.Name
	instanceName := "hss-configs"
	log.Info("createConfigMap++", "instanceName=", instanceName)

	/*n2ip, err := controllers.GetFirstInterfaceConfigIPv4(hssDeployment.Spec.Interfaces, "n2")
	if err != nil {
		log.Error(err, "Interface N2 not found in HSSDeployment Spec")
		return nil, err
	}*/

	templateValues := configurationTemplateValues{
		SVC_NAME:  instanceName,
		S1AP_PORT: 36412,
		S11_PORT:  2123,
	}

	configJson, err := renderConfigFiles(log, templateValues)
	if err != nil {
		log.Error(err, "Could not render HSS configuration template.")
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
			"acl.conf": configJson[0],
			"hss.conf": configJson[1],
			"hss.json": configJson[2],
			"oss.json": configJson[3],
		},
	}
	log.Info("createConfigMap--")
	return configMap, nil
}

func createScriptConfigMap(log logr.Logger, hssDeployment *v1alpha1.HSSDeployment) (*apiv1.ConfigMap, error) {
	namespace := hssDeployment.Namespace
	instanceName := "hss-scripts"
	log.Info("createScriptConfigMap++", "instanceName=", instanceName)

	/*n2ip, err := controllers.GetFirstInterfaceConfigIPv4(hssDeployment.Spec.Interfaces, "n2")
	if err != nil {
		log.Error(err, "Interface N2 not found in HSSDeployment Spec")
		return nil, err
	}*/

	templateValues := configurationTemplateValues{
		SVC_NAME:  instanceName,
		S1AP_PORT: 36412,
		S11_PORT:  2123,
	}

	hssScriptsConfig, err := renderConfigFiles(log, templateValues)
	if err != nil {
		log.Error(err, "Could not render HSS Scripts configuration template.")
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
			"hss-bootstrap.sh": hssScriptsConfig[0],
			"hss-run.sh":       hssScriptsConfig[1],
		},
	}
	log.Info("createScriptConfigMap--")
	return configMap, nil
}

func createResourceRequirements(hssDeploymentSpec v1alpha1.HSSDeploymentSpec) (int32, *apiv1.ResourceRequirements, error) {
	// TODO: Requirements should be calculated based on DL, UL
	// TODO: increase number of recpicas based on NFDeployment.Capacity.MaxSessions

	var replicas int32 = 1
	//downlink := resource.MustParse("5G")
	//uplink := resource.MustParse("1G")
	var cpuLimit string
	var cpuRequest string
	var memoryLimit string
	var memoryRequest string

	if hssDeploymentSpec.Capacity.MaxSubscribers > 1000 {
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

/*func createNetworkAttachmentDefinitionNetworks(templateName string, hssDeploymentSpec *v1alpha1.HSSDeploymentSpec) (string, error) {
	return controllers.CreateNetworkAttachmentDefinitionNetworks(templateName, map[string][]nephiov1alpha1.InterfaceConfig{
		"n2": controllers.GetInterfaceConfigs(hssDeploymentSpec.Interfaces, "n2"),
	})
}*/
