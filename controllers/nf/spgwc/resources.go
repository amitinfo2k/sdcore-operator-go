package spgwc

import (
	"github.com/amitinfo2k/sdcore-operator-go/controllers"
	"github.com/go-logr/logr"
	nephiov1alpha1 "github.com/nephio-project/api/nf_deployments/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/util/intstr"
)

func createDeployment(log logr.Logger, configMapVersion string, spgwcDeployment *nephiov1alpha1.NFDeployment) (*appsv1.Deployment, error) {
	namespace := spgwcDeployment.Namespace
	instanceName := spgwcDeployment.Name
	spec := spgwcDeployment.Spec

	previleged := true
	runAsUser := int64(0)
	mode := int32(493)
	configMode := int32(420)

	replicas, resourceRequirements, err := createResourceRequirements(spec)
	if err != nil {
		return nil, err
	}

	networkAttachmentDefinitionNetworks, err := createNetworkAttachmentDefinitionNetworks(spgwcDeployment.Name, &spec)
	if err != nil {
		return nil, err
	}

	podAnnotations := make(map[string]string)
	podAnnotations[controllers.ConfigMapVersionAnnotation] = configMapVersion
	podAnnotations[controllers.NetworksAnnotation] = networkAttachmentDefinitionNetworks

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
	quntity, err := resource.ParseQuantity("1Mi")

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
					Containers: []apiv1.Container{
						{
							Name:            "spgwc",
							Image:           controllers.SPGWCImage,
							ImagePullPolicy: apiv1.PullAlways,
							SecurityContext: initSecurityContext,
							Command:         []string{"sh", "-xc"},
							Args:            []string{"bash /opt/cp/scripts/spgwc-run.sh ngic_controlplane"},

							Env: []apiv1.EnvVar{
								{
									Name: "POD_IP",
									ValueFrom: &apiv1.EnvVarSource{
										FieldRef: &apiv1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
								{
									Name: "MEM_LIMIT",
									ValueFrom: &apiv1.EnvVarSource{
										ResourceFieldRef: &apiv1.ResourceFieldSelector{
											ContainerName: "spgwc",
											Resource:      "limits.memory",
											Divisor:       quntity,
										},
									},
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									MountPath: "/opt/cp/scripts",
									Name:      "scripts",
								},
								{
									MountPath: "/etc/cp/config",
									Name:      "configs",
								},
							},
							Resources: *resourceRequirements,
						},
						{
							Name:            "gx-app",
							Image:           controllers.SPGWCImage,
							ImagePullPolicy: apiv1.PullAlways,
							SecurityContext: initSecurityContext,
							Command:         []string{"sh", "-xc"},
							Args:            []string{"bash /opt/cp/scripts/spgwc-run.sh gx-app"},
							Env: []apiv1.EnvVar{
								{
									Name: "POD_IP",
									ValueFrom: &apiv1.EnvVarSource{
										FieldRef: &apiv1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
								{
									Name:  "MANAGED_BY_CONFIG_POD",
									Value: "true",
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									MountPath: "/opt/cp/scripts",
									Name:      "scripts",
								},
								{
									MountPath: "/etc/cp/config",
									Name:      "configs",
								},
							},
							Resources: *resourceRequirements,
						},
						{
							Name:            "init-sync",
							Image:           controllers.CurlImage,
							ImagePullPolicy: apiv1.PullAlways,
							Stdin:           true,
							TTY:             true,
							Command:         []string{"sh", "-xc"},
							Args:            []string{"sh /opt/cp/scripts/spgwc-init.sh"},

							VolumeMounts: []apiv1.VolumeMount{
								{
									MountPath: "/opt/cp/scripts",
									Name:      "scripts",
								},
								{
									MountPath: "/etc/cp/config",
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
										Name: "spgwc-scripts",
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
										Name: "spgwc-configs",
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

func createService(spgwcDeployment *nephiov1alpha1.NFDeployment) *apiv1.Service {
	namespace := spgwcDeployment.Namespace
	//instanceName := spgwcDeployment.Name
	instanceName := "spgwc"

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
					Name:     "s11",
					Protocol: apiv1.ProtocolTCP,
					Port:     2123,
					NodePort: 32123,
				},
				{
					Name:     "pfcp",
					Protocol: apiv1.ProtocolUDP,
					Port:     8805,
					NodePort: 30021,
				},
				{
					Name:     "rest",
					Protocol: apiv1.ProtocolSCTP,
					Port:     8080,
					NodePort: 30080,
				}, {
					Name:     "prometheus-exporter",
					Protocol: apiv1.ProtocolTCP,
					Port:     3081,
					NodePort: 30084,
				},
			},
			Type: apiv1.ServiceTypeNodePort,
		},
	}

	return service
}

func createConfigMap(log logr.Logger, spgwcDeployment *nephiov1alpha1.NFDeployment) (*apiv1.ConfigMap, error) {
	namespace := spgwcDeployment.Namespace
	//instanceName := spgwcDeployment.Name
	instanceName := "spgwc-configs"
	log.Info("createConfigMap++", "instanceName=", instanceName)

	n4ip, err := controllers.GetFirstInterfaceConfigIPv4(spgwcDeployment.Spec.Interfaces, "n4")
	if err != nil {
		log.Error(err, "Interface PFCP not found in SPGWCDeployment Spec")
		return nil, err
	}

	templateValues := configurationTemplateValues{
		N4_IP:     n4ip,
		S1AP_PORT: 36412,
		S11_PORT:  2123,
	}

	configJson, err := renderConfigFiles(log, templateValues)
	if err != nil {
		log.Error(err, "Could not render SPGWC configuration template.")
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
			"config.json":             configJson[0],
			"cp.json":                 configJson[1],
			"gx.conf":                 configJson[2],
			"subscriber_mapping.json": configJson[3],
		},
	}
	log.Info("createConfigMap--")
	return configMap, nil
}

func createScriptConfigMap(log logr.Logger, spgwcDeployment *nephiov1alpha1.NFDeployment) (*apiv1.ConfigMap, error) {
	namespace := spgwcDeployment.Namespace
	instanceName := "spgwc-scripts"
	log.Info("createScriptConfigMap++", "instanceName=", instanceName)

	n4ip, err := controllers.GetFirstInterfaceConfigIPv4(spgwcDeployment.Spec.Interfaces, "n4")
	if err != nil {
		log.Error(err, "Interface PFCP not found in SPGWCDeployment Spec")
		return nil, err
	}

	templateValues := configurationTemplateValues{
		N4_IP:     n4ip,
		S1AP_PORT: 36412,
		S11_PORT:  2123,
	}

	spgwcScriptsConfig, err := renderScriptFiles(log, templateValues)
	if err != nil {
		log.Error(err, "Could not render SPGWC Scripts configuration template.")
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
			"spgwc-init.sh": spgwcScriptsConfig[0],
			"spgwc-run.sh":  spgwcScriptsConfig[1],
		},
	}
	log.Info("createScriptConfigMap--")
	return configMap, nil
}

func createResourceRequirements(spgwcDeploymentSpec nephiov1alpha1.NFDeploymentSpec) (int32, *apiv1.ResourceRequirements, error) {
	// TODO: Requirements should be calculated based on DL, UL
	// TODO: increase number of recpicas based on NFDeployment.Capacity.MaxSessions

	var replicas int32 = 1
	//downlink := resource.MustParse("5G")
	//uplink := resource.MustParse("1G")
	var cpuLimit string
	var cpuRequest string
	var memoryLimit string
	var memoryRequest string

	if spgwcDeploymentSpec.Capacity.MaxSubscribers > 1000 {
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

func createNetworkAttachmentDefinitionNetworks(templateName string, spgwcDeploymentSpec *nephiov1alpha1.NFDeploymentSpec) (string, error) {
	return controllers.CreateNetworkAttachmentDefinitionNetworks(templateName, map[string][]nephiov1alpha1.InterfaceConfig{
		"n4": controllers.GetInterfaceConfigs(spgwcDeploymentSpec.Interfaces, "n4"),
	})
}
