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

// SPGWCDeploymentSpec defines the desired state of SPGWCDeployment
type SPGWCDeploymentSpec struct {
	nephiov1alpha1.NFDeploymentSpec `json:",inline" yaml:",inline"`
}

// SPGWCDeploymentStatus defines the observed state of SPGWCDeployment
type SPGWCDeploymentStatus struct {
	nephiov1alpha1.NFDeploymentStatus `json:",inline" yaml:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SPGWCDeployment is the Schema for the spgwcdeployments API
type SPGWCDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SPGWCDeploymentSpec   `json:"spec,omitempty"`
	Status SPGWCDeploymentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SPGWCDeploymentList contains a list of SPGWCDeployment
type SPGWCDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SPGWCDeployment `json:"items"`
}

// Implement NFDeployment interface

func (d *SPGWCDeployment) GetNFDeploymentSpec() *nephiov1alpha1.NFDeploymentSpec {
	return d.Spec.NFDeploymentSpec.DeepCopy()
}
func (d *SPGWCDeployment) GetNFDeploymentStatus() *nephiov1alpha1.NFDeploymentStatus {
	return d.Status.NFDeploymentStatus.DeepCopy()
}
func (d *SPGWCDeployment) SetNFDeploymentSpec(s *nephiov1alpha1.NFDeploymentSpec) {
	s.DeepCopyInto(&d.Spec.NFDeploymentSpec)
}
func (d *SPGWCDeployment) SetNFDeploymentStatus(s *nephiov1alpha1.NFDeploymentStatus) {
	s.DeepCopyInto(&d.Status.NFDeploymentStatus)
}

// Interface type metadata.
var (
	SPGWCDeploymentKind             = reflect.TypeOf(SPGWCDeployment{}).Name()
	SPGWCDeploymentGroupKind        = schema.GroupKind{Group: nephiov1alpha1.Group, Kind: SPGWCDeploymentKind}.String()
	SPGWCDeploymentKindAPIVersion   = SPGWCDeploymentKind + "." + nephiov1alpha1.GroupVersion.String()
	SPGWCDeploymentGroupVersionKind = nephiov1alpha1.GroupVersion.WithKind(SPGWCDeploymentKind)
)

func init() {
	nephiov1alpha1.SchemeBuilder.Register(&SPGWCDeployment{}, &SPGWCDeploymentList{})
}
