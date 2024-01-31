/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"reflect"

	nephiov1alpha1 "github.com/nephio-project/api/nf_deployments/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// Config4GDeploymentSpec defines the desired state of Config4GDeployment
type Config4GDeploymentSpec struct {
	nephiov1alpha1.NFDeploymentSpec `json:",inline" yaml:",inline"`
}

// Config4GDeploymentStatus defines the observed state of Config4GDeployment
type Config4GDeploymentStatus struct {
	nephiov1alpha1.NFDeploymentStatus `json:",inline" yaml:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Config4GDeployment is the Schema for the pcrfdeployments API
type Config4GDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   Config4GDeploymentSpec   `json:"spec,omitempty"`
	Status Config4GDeploymentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// Config4GDeploymentList contains a list of Config4GDeployment
type Config4GDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Config4GDeployment `json:"items"`
}

// Implement NFDeployment interface

func (d *Config4GDeployment) GetNFDeploymentSpec() *nephiov1alpha1.NFDeploymentSpec {
	return d.Spec.NFDeploymentSpec.DeepCopy()
}
func (d *Config4GDeployment) GetNFDeploymentStatus() *nephiov1alpha1.NFDeploymentStatus {
	return d.Status.NFDeploymentStatus.DeepCopy()
}
func (d *Config4GDeployment) SetNFDeploymentSpec(s *nephiov1alpha1.NFDeploymentSpec) {
	s.DeepCopyInto(&d.Spec.NFDeploymentSpec)
}
func (d *Config4GDeployment) SetNFDeploymentStatus(s *nephiov1alpha1.NFDeploymentStatus) {
	s.DeepCopyInto(&d.Status.NFDeploymentStatus)
}

// Interface type metadata.
var (
	Config4GDeploymentKind             = reflect.TypeOf(Config4GDeployment{}).Name()
	Config4GDeploymentGroupKind        = schema.GroupKind{Group: nephiov1alpha1.Group, Kind: Config4GDeploymentKind}.String()
	Config4GDeploymentKindAPIVersion   = Config4GDeploymentKind + "." + nephiov1alpha1.GroupVersion.String()
	Config4GDeploymentGroupVersionKind = nephiov1alpha1.GroupVersion.WithKind(Config4GDeploymentKind)
)

func init() {
	nephiov1alpha1.SchemeBuilder.Register(&Config4GDeployment{}, &Config4GDeploymentList{})
}
