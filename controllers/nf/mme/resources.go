package mme

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

func createDeployment(log logr.Logger, configMapVersion string, mmeDeployment *nephiov1alpha1.NFDeployment) (*appsv1.Deployment, error) {
	namespace := mmeDeployment.Namespace
	instanceName := mmeDeployment.Name
	spec := mmeDeployment.Spec

	previleged := true
	runAsUser := int64(0)
	mode := int32(493)

	replicas, resourceRequirements, err := createResourceRequirements(spec)
	if err != nil {
		return nil, err
	}

	networkAttachmentDefinitionNetworks, err := createNetworkAttachmentDefinitionNetworks(mmeDeployment.Name, &spec)
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
							Name:            "mme-load-sctp-module",
							Image:           controllers.MMEInitImage,
							ImagePullPolicy: apiv1.PullIfNotPresent,
							SecurityContext: initSecurityContext,
							Command:         []string{"sh", "-xc"},
							Args:            []string{"if chroot /mnt/host-rootfs modinfo nf_conntrack_proto_sctp > /dev/null 2>&1; then chroot /mnt/host-rootfs modprobe nf_conntrack_proto_sctp; fi; chroot /mnt/host-rootfs modprobe tipc"},
							VolumeMounts: []apiv1.VolumeMount{
								{
									MountPath: "/mnt/host-rootfs",
									Name:      "host-rootfs",
								},
							},
						},
						{
							Name:            "mme-init",
							Image:           controllers.MMEImage,
							ImagePullPolicy: apiv1.PullIfNotPresent,
							Command:         []string{"sh", "-xc"},
							Args:            []string{"sh /opt/mme/scripts/mme-init.sh"},
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
									MountPath: "/opt/mme/scripts",
									Name:      "scripts",
								},
								{
									MountPath: "/opt/mme/config",
									Name:      "configs",
								},
								{
									MountPath: "/opt/mme/config/shared",
									Name:      "shared-data",
								},
							},
							Resources: *resourceRequirements,
						},
					},
					Containers: []apiv1.Container{
						{
							Name:            "mme-app",
							Image:           controllers.MMEImage,
							ImagePullPolicy: apiv1.PullAlways,

							Command: []string{"sh", "-xc"},
							Args:    []string{"sh /opt/mme/scripts/mme-run.sh mme-app"},

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
									MountPath: "/opt/mme/scripts",
									Name:      "scripts",
								},
								{
									MountPath: "/opt/mme/config",
									Name:      "configs",
								},
								{
									MountPath: "/opt/mme/config/shared",
									Name:      "shared-data",
								},
								{
									MountPath: "/tmp",
									Name:      "shared-app",
								},
							},
							Resources: *resourceRequirements,
						},
						{
							Name:            "s1ap-app",
							Image:           controllers.MMEImage,
							ImagePullPolicy: apiv1.PullAlways,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "s1ap",
									Protocol:      apiv1.ProtocolSCTP,
									ContainerPort: 36412,
								},
							},
							Command: []string{"sh", "-xc"},
							Args:    []string{"sh /opt/mme/scripts/mme-run.sh s1ap-app"},
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
									Name:  "MMERUNENV",
									Value: "container",
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									MountPath: "/opt/mme/scripts",
									Name:      "scripts",
								},
								{
									MountPath: "/opt/mme/config",
									Name:      "configs",
								},
								{
									MountPath: "/opt/mme/config/shared",
									Name:      "shared-data",
								},
								{
									MountPath: "/tmp",
									Name:      "shared-app",
								},
							},
							Resources: *resourceRequirements,
						},
						{
							Name:            "s11-app",
							Image:           controllers.MMEImage,
							ImagePullPolicy: apiv1.PullAlways,
							Command:         []string{"sh", "-xc"},
							Args:            []string{"sh /opt/mme/scripts/mme-run.sh s11-app"},
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
									Name:  "MMERUNENV",
									Value: "container",
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									MountPath: "/opt/mme/scripts",
									Name:      "scripts",
								},
								{
									MountPath: "/opt/mme/config/shared",
									Name:      "shared-data",
								},
								{
									MountPath: "/tmp",
									Name:      "shared-app",
								},
							},
							Resources: *resourceRequirements,
						},
						{
							Name:            "s6a-app",
							Image:           controllers.MMEImage,
							ImagePullPolicy: apiv1.PullAlways,
							Command:         []string{"sh", "-xc"},
							Args:            []string{"sh /opt/mme/scripts/mme-run.sh s6a-app"},
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
									Name:  "MMERUNENV",
									Value: "container",
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									MountPath: "/opt/mme/scripts",
									Name:      "scripts",
								},
								{
									MountPath: "/opt/mme/config/shared",
									Name:      "shared-data",
								},
								{
									MountPath: "/tmp",
									Name:      "shared-app",
								},
							},
							Resources: *resourceRequirements,
						},
					}, // Containers
					DNSPolicy:     apiv1.DNSClusterFirst,
					RestartPolicy: apiv1.RestartPolicyAlways,
					Volumes: []apiv1.Volume{
						{
							Name: "mme-volume",
							VolumeSource: apiv1.VolumeSource{
								Projected: &apiv1.ProjectedVolumeSource{
									Sources: []apiv1.VolumeProjection{
										{
											ConfigMap: &apiv1.ConfigMapProjection{
												LocalObjectReference: apiv1.LocalObjectReference{
													Name: instanceName,
												},
												Items: []apiv1.KeyToPath{
													{
														Key:  "config.json",
														Path: "config.json",
													},
												},
											},
										},
									},
								},
							},
						},
						{
							Name: "scripts",
							VolumeSource: apiv1.VolumeSource{
								ConfigMap: &apiv1.ConfigMapVolumeSource{
									LocalObjectReference: apiv1.LocalObjectReference{
										Name: "mme-scripts",
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
										Name: "mme-configs",
									},
									DefaultMode: &mode,
								},
							},
						},
						{
							Name: "shared-data",
							VolumeSource: apiv1.VolumeSource{
								EmptyDir: &apiv1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "shared-app",
							VolumeSource: apiv1.VolumeSource{
								EmptyDir: &apiv1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "host-rootfs",
							VolumeSource: apiv1.VolumeSource{
								HostPath: &apiv1.HostPathVolumeSource{
									Path: "/",
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

func createService(mmeDeployment *nephiov1alpha1.NFDeployment) *apiv1.Service {
	namespace := mmeDeployment.Namespace
	//instanceName := mmeDeployment.Name
	instanceName := "mme"

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
					Protocol: apiv1.ProtocolUDP,
					Port:     2123,
					NodePort: 32124,
				},
				{
					Name:     "s6a",
					Protocol: apiv1.ProtocolTCP,
					Port:     3868,
					NodePort: 32269,
				},
				{
					Name:     "s1ap",
					Protocol: apiv1.ProtocolSCTP,
					Port:     80,
					NodePort: 32412,
				}, {
					Name:     "prometheus-exporter",
					Protocol: apiv1.ProtocolTCP,
					Port:     3081,
					NodePort: 30085,
				},
				{
					Name:     "mme-app-config",
					Protocol: apiv1.ProtocolTCP,
					Port:     8080,
				},
				{
					Name:     "mme-s1ap-config",
					Protocol: apiv1.ProtocolTCP,
					Port:     8081,
				},
			},
			Type: apiv1.ServiceTypeNodePort,
		},
	}

	return service
}

func createConfigMap(log logr.Logger, mmeDeployment *nephiov1alpha1.NFDeployment) (*apiv1.ConfigMap, error) {
	namespace := mmeDeployment.Namespace
	//instanceName := mmeDeployment.Name
	instanceName := "mme-configs"
	log.Info("createConfigMap++", "instanceName=", instanceName)

	_, err := controllers.GetFirstInterfaceConfigIPv4(mmeDeployment.Spec.Interfaces, "n2")
	if err != nil {
		log.Error(err, "Interface S1AP not found in MMEDeployment Spec")
		return nil, err
	}

	templateValues := configurationTemplateValues{
		SVC_NAME:  instanceName,
		S1AP_PORT: 36412,
		S11_PORT:  2123,
	}

	configJson, err := renderConfigFiles(log, templateValues)
	if err != nil {
		log.Error(err, "Could not render MME configuration template.")
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
			"config.json": configJson[0],
			"s6a_fd.conf": configJson[1],
		},
	}
	log.Info("createConfigMap--")
	return configMap, nil
}

func createScriptConfigMap(log logr.Logger, mmeDeployment *nephiov1alpha1.NFDeployment) (*apiv1.ConfigMap, error) {
	namespace := mmeDeployment.Namespace
	instanceName := "mme-scripts"
	log.Info("createScriptConfigMap++", "instanceName=", instanceName)

	_, err := controllers.GetFirstInterfaceConfigIPv4(mmeDeployment.Spec.Interfaces, "n2")
	if err != nil {
		log.Error(err, "Interface s1ap not found in MMEDeployment Spec")
		return nil, err
	}

	templateValues := configurationTemplateValues{
		SVC_NAME:  instanceName,
		S1AP_PORT: 36412,
		S11_PORT:  2123,
	}

	mmeScriptsConfig, err := renderScriptFiles(log, templateValues)
	if err != nil {
		log.Error(err, "Could not render MME Scripts configuration template.")
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
			"mme-init.sh": mmeScriptsConfig[0],
			"mme-run.sh":  mmeScriptsConfig[1],
		},
	}
	log.Info("createScriptConfigMap--")
	return configMap, nil
}

func createResourceRequirements(mmeDeploymentSpec nephiov1alpha1.NFDeploymentSpec) (int32, *apiv1.ResourceRequirements, error) {
	// TODO: Requirements should be calculated based on DL, UL
	// TODO: increase number of recpicas based on NFDeployment.Capacity.MaxSessions

	var replicas int32 = 1
	//downlink := resource.MustParse("5G")
	//uplink := resource.MustParse("1G")
	var cpuLimit string
	var cpuRequest string
	var memoryLimit string
	var memoryRequest string

	if mmeDeploymentSpec.Capacity.MaxSubscribers > 1000 {
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

func createNetworkAttachmentDefinitionNetworks(templateName string, mmeDeploymentSpec *nephiov1alpha1.NFDeploymentSpec) (string, error) {
	return controllers.CreateNetworkAttachmentDefinitionNetworks(templateName, map[string][]nephiov1alpha1.InterfaceConfig{
		"n2": controllers.GetInterfaceConfigs(mmeDeploymentSpec.Interfaces, "n2"),
	})
}
