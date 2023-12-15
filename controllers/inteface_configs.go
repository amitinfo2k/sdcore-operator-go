package controllers

import (
	"fmt"
	"net"

	nephiov1alpha1 "github.com/nephio-project/api/nf_deployments/v1alpha1"
)

func GetInterfaceConfigs(interfaceConfigs []nephiov1alpha1.InterfaceConfig, interfaceName string) []nephiov1alpha1.InterfaceConfig {
	var selectedInterfaceConfigs []nephiov1alpha1.InterfaceConfig

	for _, interfaceConfig := range interfaceConfigs {
		if interfaceConfig.Name == interfaceName {
			selectedInterfaceConfigs = append(selectedInterfaceConfigs, interfaceConfig)
		}
	}

	return selectedInterfaceConfigs
}

func GetFirstInterfaceConfig(interfaceConfigs []nephiov1alpha1.InterfaceConfig, interfaceName string) (*nephiov1alpha1.InterfaceConfig, error) {
	for _, interfaceConfig := range interfaceConfigs {
		fmt.Println("GetFirstInterfaceConfig::", interfaceConfig.Name)
		if interfaceConfig.Name == interfaceName {
			return &interfaceConfig, nil
		}
	}

	return nil, fmt.Errorf("Interface %q not found", interfaceName)
}

func GetFirstInterfaceConfigIPv4(interfaceConfigs []nephiov1alpha1.InterfaceConfig, interfaceName string) (string, error) {
	interfaceConfig, err := GetFirstInterfaceConfig(interfaceConfigs, interfaceName)
	if err != nil {
		return "", err
	}

	ip, _, err := net.ParseCIDR(interfaceConfig.IPv4.Address)
	if err != nil {
		return "", err
	}

	return ip.String(), nil
}
