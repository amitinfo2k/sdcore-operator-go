package upf

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	schedulingv1 "k8s.io/api/scheduling/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
)

func deleteMeAfterDeletingUnusedImportedModules() {
	/*
		It is written to handle the error "Module Imported but not used",
		The user can delete the non-required modules from import and then delete this function also
	*/
	_ = time.Now()
	_ = &unstructured.Unstructured{}
	_ = corev1.Service{}
	_ = metav1.ObjectMeta{}
	_ = appsv1.Deployment{}
	_ = rbacv1.Role{}
	_ = schedulingv1.PriorityClass{}
	_ = intstr.FromInt(4)
	_, _ = resource.ParseQuantity("")
	_ = context.TODO()
	_ = fmt.Sprintf("")
	_ = ptr.To(32)
}

func int32Ptr(val int) *int32 {
	var a int32
	a = int32(val)
	return &a
}

func int64Ptr(val int) *int64 {
	var a int64
	a = int64(val)
	return &a
}

func intPtr(val int) *int {
	a := val
	return &a
}

func int16Ptr(val int) *int16 {
	var a int16
	a = int16(val)
	return &a
}

func boolPtr(val bool) *bool {
	a := val
	return &a
}

func stringPtr(val string) *string {
	a := val
	return &a
}

func getDataForSecret(encodedVal string) []byte {
	/*
		Concept: Based on my Understanding, corev1.Secret requires the actual data(not encoded) as secret-Data
		But in general terms, we put encoded values in secret-data, which make sense (why to write actual value in readable format)
		This function takes the encodedVal and decodes it and returns
	*/
	decodeVal, err := base64.StdEncoding.DecodeString(encodedVal)
	if err != nil {
		fmt.Println("Unable to decode the SecretVal ", encodedVal, " || This Secret Will Probably would give error during deployment| Kindly Check")
		return []byte(encodedVal)
	}
	return decodeVal
}

// Before Uncommenting the following function, Make sure the data-type of r is same as of your Reconciler,
// Replace "UPFDeploymentReconciler" with the type of your Reconciler
func (r *UPFDeploymentReconciler) CreateAll() {
	var err error
	namespaceProvided := "omec"

	for _, resource := range GetConfigMap() {
		if resource.ObjectMeta.Namespace == "" {
			resource.ObjectMeta.Namespace = namespaceProvided
		}
		err = r.Create(context.TODO(), resource)
		if err != nil {
			fmt.Println("Erorr During Creating resource of GetConfigMap()| Error --> |", err)
		}
	}

	for _, resource := range GetNetworkAttachmentDefinition() {
		if resource.GetNamespace() == "" {
			resource.SetNamespace(namespaceProvided)
		}
		err = r.Create(context.TODO(), resource)
		if err != nil {
			fmt.Println("Erorr During Creating resource of GetNetworkAttachmentDefinition()| Error --> |", err)
		}
	}

	for _, resource := range GetService() {
		if resource.ObjectMeta.Namespace == "" {
			resource.ObjectMeta.Namespace = namespaceProvided
		}
		err = r.Create(context.TODO(), resource)
		if err != nil {
			fmt.Println("Erorr During Creating resource of GetService()| Error --> |", err)
		}
	}

	for _, resource := range GetStatefulSet() {
		if resource.ObjectMeta.Namespace == "" {
			resource.ObjectMeta.Namespace = namespaceProvided
		}
		err = r.Create(context.TODO(), resource)
		if err != nil {
			fmt.Println("Erorr During Creating resource of GetStatefulSet()| Error --> |", err)
		}
	}

}

// Before Uncommenting the following function, Make sure the data-type of r is same as of your Reconciler,
// Replace "UPFDeploymentReconciler" with the type of your Reconciler
func (r *UPFDeploymentReconciler) DeleteAll() {
	var err error
	namespaceProvided := "omec"

	for _, resource := range GetConfigMap() {
		if resource.ObjectMeta.Namespace == "" {
			resource.ObjectMeta.Namespace = namespaceProvided
		}
		err = r.Delete(context.TODO(), resource)
		if err != nil {
			fmt.Println("Erorr During Deleting resource of GetConfigMap()| Error --> |", err)
		}
	}

	for _, resource := range GetNetworkAttachmentDefinition() {
		if resource.GetNamespace() == "" {
			resource.SetNamespace(namespaceProvided)
		}
		err = r.Delete(context.TODO(), resource)
		if err != nil {
			fmt.Println("Erorr During Deleting resource of GetNetworkAttachmentDefinition()| Error --> |", err)
		}
	}

	for _, resource := range GetService() {
		if resource.ObjectMeta.Namespace == "" {
			resource.ObjectMeta.Namespace = namespaceProvided
		}
		err = r.Delete(context.TODO(), resource)
		if err != nil {
			fmt.Println("Erorr During Deleting resource of GetService()| Error --> |", err)
		}
	}

	for _, resource := range GetStatefulSet() {
		if resource.ObjectMeta.Namespace == "" {
			resource.ObjectMeta.Namespace = namespaceProvided
		}
		err = r.Delete(context.TODO(), resource)
		if err != nil {
			fmt.Println("Erorr During Deleting resource of GetStatefulSet()| Error --> |", err)
		}
	}

}

func GetConfigMap() []*corev1.ConfigMap {

	configMap1 := &corev1.ConfigMap{
		Data: map[string]string{
			"upf.json": "{\"access\":{\"ifname\":\"access\"},\"core\":{\"ifname\":\"core\"},\"cpiface\":{\"dnn\":\"internet\",\"hostname\":\"upf\",\"http_port\":\"8080\"},\"enable_notify_bess\":true,\"max_sessions\":50000,\"measure_flow\":false,\"measure_upf\":true,\"mode\":\"af_packet\",\"notify_sockaddr\":\"/pod-share/notifycp\",\"qci_qos_config\":[{\"burst_duration_ms\":10,\"cbs\":50000,\"ebs\":50000,\"pbs\":50000,\"priority\":7,\"qci\":0}],\"slice_rate_limit_config\":{\"n3_bps\":1000000000,\"n3_burst_bytes\":12500000,\"n6_bps\":1000000000,\"n6_burst_bytes\":12500000},\"table_sizes\":{\"appQERLookup\":200000,\"farLookup\":150000,\"pdrLookup\":50000,\"sessionQERLookup\":100000},\"workers\":1}",
			"bessd-poststart.sh": "#!/bin/bash\n" +
				"\n" +
				"# Copyright 2020-present Open Networking Foundation\n" +
				"#\n" +
				"# SPDX-License-Identifier: Apache-2.0\n" +
				"\n" +
				"set -ex\n" +
				"\n" +
				"until bessctl run /opt/bess/bessctl/conf/up4; do\n" +
				"    sleep 2;\n" +
				"done;\n" +
				"",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "upf",
			Labels: map[string]string{
				"release": "release-name",
				"app":     "upf",
			},
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
	}

	return []*corev1.ConfigMap{configMap1}
}

func GetNetworkAttachmentDefinition() []*unstructured.Unstructured {

	networkAttachmentDefinition1 := &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "k8s.cni.cncf.io/v1",
			"kind":       "NetworkAttachmentDefinition",
			"metadata": map[string]any{
				"name": "access-net",
			},
			"spec": map[string]any{
				"config": "{ \"cniVersion\": \"0.3.1\", \"type\": \"macvlan\", \"master\": \"data\", \"ipam\": { \"type\": \"static\" }, \"capabilities\": { \"mac\": true} }",
			},
		},
	}

	networkAttachmentDefinition2 := &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "k8s.cni.cncf.io/v1",
			"kind":       "NetworkAttachmentDefinition",
			"metadata": map[string]any{
				"name": "core-net",
			},
			"spec": map[string]any{
				"config": "{ \"cniVersion\": \"0.3.1\", \"type\": \"macvlan\", \"master\": \"data\", \"ipam\": { \"type\": \"static\" }, \"capabilities\": { \"mac\": true} }",
			},
		},
	}

	return []*unstructured.Unstructured{networkAttachmentDefinition1, networkAttachmentDefinition2}
}

func GetService() []*corev1.Service {

	service1 := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app":     "upf",
				"release": "release-name",
			},
			Name: "upf",
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{

				corev1.ServicePort{
					Name:     "pfcp",
					Port:     8805,
					Protocol: corev1.Protocol("UDP"),
				},
			},
			PublishNotReadyAddresses: false,
			Selector: map[string]string{
				"app":     "upf",
				"release": "release-name",
			},
			Type: corev1.ServiceType("ClusterIP"),
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
	}

	service2 := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app":     "upf",
				"release": "release-name",
			},
			Name: "upf-http",
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{

				corev1.ServicePort{
					Name:     "bess-web",
					Port:     8000,
					Protocol: corev1.Protocol("TCP"),
				},
				corev1.ServicePort{
					Name:     "prometheus-exporter",
					Port:     8080,
					Protocol: corev1.Protocol("TCP"),
				},
			},
			PublishNotReadyAddresses: false,
			Selector: map[string]string{
				"app":     "upf",
				"release": "release-name",
			},
			Type: corev1.ServiceType("ClusterIP"),
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
	}

	return []*corev1.Service{service1, service2}
}

func GetStatefulSet() []*appsv1.StatefulSet {

	statefulSet1 := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app":     "upf",
				"release": "release-name",
			},
			Name: "upf",
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: "upf-headless",
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					HostNetwork: false,
					HostPID:     false,
					ImagePullSecrets: []corev1.LocalObjectReference{

						corev1.LocalObjectReference{
							Name: "aether.registry",
						},
					},
					InitContainers: []corev1.Container{

						corev1.Container{
							StdinOnce: false,
							Command: []string{

								"sh",
								"-xec",
							},
							ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
							Name:            "bess-init",
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									"cpu":    resource.MustParse("128m"),
									"memory": resource.MustParse("64Mi"),
								},
								Requests: corev1.ResourceList{
									"cpu":    resource.MustParse("128m"),
									"memory": resource.MustParse("64Mi"),
								},
							},
							Stdin: false,
							Args: []string{

								"ip route replace 192.168.251.0/24 via ; ip route replace default via  metric 110; iptables -I OUTPUT -p icmp --icmp-type port-unreachable -j DROP;",
							},
							Image: "omecproject/upf-epc-bess:master-5786085",
							SecurityContext: &corev1.SecurityContext{
								Capabilities: &corev1.Capabilities{
									Add: []corev1.Capability{

										corev1.Capability("NET_ADMIN"),
									},
								},
							},
							TTY: false,
						},
					},
					ShareProcessNamespace: boolPtr(true),
					Volumes: []corev1.Volume{

						corev1.Volume{
							Name: "configs",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "upf",
									},
									DefaultMode: int32Ptr(493),
								},
							},
						},
						corev1.Volume{
							Name: "shared-app",
						},
					},
					Containers: []corev1.Container{

						corev1.Container{
							ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
							Lifecycle: &corev1.Lifecycle{
								PostStart: &corev1.LifecycleHandler{
									Exec: &corev1.ExecAction{
										Command: []string{

											"/etc/bess/conf/bessd-poststart.sh",
										},
									},
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									"cpu":    resource.MustParse("2"),
									"memory": resource.MustParse("1Gi"),
								},
								Limits: corev1.ResourceList{
									"cpu":    resource.MustParse("2"),
									"memory": resource.MustParse("1Gi"),
								},
							},
							Env: []corev1.EnvVar{

								corev1.EnvVar{
									Name:  "CONF_FILE",
									Value: "/etc/bess/conf/upf.json",
								},
							},
							LivenessProbe: &corev1.Probe{
								InitialDelaySeconds: 15,
								PeriodSeconds:       20,
								ProbeHandler: corev1.ProbeHandler{
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.IntOrString{
											IntVal: 10514,
										},
									},
								},
							},
							Stdin:     true,
							StdinOnce: false,
							Image:     "omecproject/upf-epc-bess:master-5786085",
							Name:      "bessd",
							SecurityContext: &corev1.SecurityContext{
								Capabilities: &corev1.Capabilities{
									Add: []corev1.Capability{

										corev1.Capability("IPC_LOCK"),
									},
								},
							},
							TTY: true,
							VolumeMounts: []corev1.VolumeMount{

								corev1.VolumeMount{
									Name:      "shared-app",
									ReadOnly:  false,
									MountPath: "/pod-share",
								},
								corev1.VolumeMount{
									MountPath: "/etc/bess/conf",
									Name:      "configs",
									ReadOnly:  false,
								},
							},
							Command: []string{

								"/bin/bash",
								"-xc",
							},
							Args: []string{

								"bessd -m 0 -f -grpc-url=0.0.0.0:10514",
							},
						},
						corev1.Container{
							Command: []string{

								"/opt/bess/bessctl/conf/route_control.py",
							},
							ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
							Name:            "routectl",
							Stdin:           false,
							Args: []string{

								"-i",
								"access",
								"core",
							},
							Image: "omecproject/upf-epc-bess:master-5786085",
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									"memory": resource.MustParse("128Mi"),
									"cpu":    resource.MustParse("256m"),
								},
								Requests: corev1.ResourceList{
									"cpu":    resource.MustParse("256m"),
									"memory": resource.MustParse("128Mi"),
								},
							},
							StdinOnce: false,
							TTY:       false,
							Env: []corev1.EnvVar{

								corev1.EnvVar{
									Name:  "PYTHONUNBUFFERED",
									Value: "1",
								},
							},
						},
						corev1.Container{
							Image:           "omecproject/upf-epc-bess:master-5786085",
							ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
							Name:            "web",
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									"cpu":    resource.MustParse("256m"),
									"memory": resource.MustParse("128Mi"),
								},
								Requests: corev1.ResourceList{
									"memory": resource.MustParse("128Mi"),
									"cpu":    resource.MustParse("256m"),
								},
							},
							Stdin:     false,
							StdinOnce: false,
							TTY:       false,
							Command: []string{

								"/bin/bash",
								"-xc",
								"bessctl http 0.0.0.0 8000",
							},
						},
						corev1.Container{
							Name: "pfcp-agent",
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									"cpu":    resource.MustParse("256m"),
									"memory": resource.MustParse("128Mi"),
								},
								Requests: corev1.ResourceList{
									"cpu":    resource.MustParse("256m"),
									"memory": resource.MustParse("128Mi"),
								},
							},
							StdinOnce: false,
							TTY:       false,
							Stdin:     false,
							VolumeMounts: []corev1.VolumeMount{

								corev1.VolumeMount{
									ReadOnly:  false,
									MountPath: "/pod-share",
									Name:      "shared-app",
								},
								corev1.VolumeMount{
									MountPath: "/tmp/conf",
									Name:      "configs",
									ReadOnly:  false,
								},
							},
							Args: []string{

								"-config",
								"/tmp/conf/upf.json",
							},
							Command: []string{

								"pfcpiface",
							},
							Image:           "omecproject/upf-epc-pfcpiface:master-5786085",
							ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
						},
						corev1.Container{
							ImagePullPolicy: corev1.PullPolicy("IfNotPresent"),
							TTY:             false,
							Args: []string{

								"while true; do\n" +
									"  # arping does not work - BESS graph is still disconnected\n" +
									"  #arping -c 2 -I access \n" +
									"  #arping -c 2 -I core \n" +
									"  ping -c 2 \n" +
									"  ping -c 2 \n" +
									"  sleep 10\n" +
									"done\n" +
									"",
							},
							Command: []string{

								"sh",
								"-xc",
							},
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									"cpu":    resource.MustParse("128m"),
									"memory": resource.MustParse("64Mi"),
								},
								Requests: corev1.ResourceList{
									"cpu":    resource.MustParse("128m"),
									"memory": resource.MustParse("64Mi"),
								},
							},
							Stdin:     false,
							StdinOnce: false,
							Image:     "registry.aetherproject.org/tools/busybox:stable",
							Name:      "arping",
						},
					},
					HostIPC: false,
				},
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"k8s.v1.cni.cncf.io/networks": "[ { \"name\": \"access-net\", \"interface\": \"access\", \"ips\": [] }, { \"name\": \"core-net\", \"interface\": \"core\", \"ips\": [] } ]",
					},
					Labels: map[string]string{
						"app":     "upf",
						"release": "release-name",
					},
				},
			},
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":     "upf",
					"release": "release-name",
				},
			},
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "StatefulSet",
		},
	}

	return []*appsv1.StatefulSet{statefulSet1}
}
