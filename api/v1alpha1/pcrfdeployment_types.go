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

// PCRFDeploymentSpec defines the desired state of PCRFDeployment
type PCRFDeploymentSpec struct {
	nephiov1alpha1.NFDeploymentSpec `json:",inline" yaml:",inline"`
}

// PCRFDeploymentStatus defines the observed state of PCRFDeployment
type PCRFDeploymentStatus struct {
	nephiov1alpha1.NFDeploymentStatus `json:",inline" yaml:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PCRFDeployment is the Schema for the pcrfdeployments API
type PCRFDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PCRFDeploymentSpec   `json:"spec,omitempty"`
	Status PCRFDeploymentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PCRFDeploymentList contains a list of PCRFDeployment
type PCRFDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PCRFDeployment `json:"items"`
}

// Implement NFDeployment interface

func (d *PCRFDeployment) GetNFDeploymentSpec() *nephiov1alpha1.NFDeploymentSpec {
	return d.Spec.NFDeploymentSpec.DeepCopy()
}
func (d *PCRFDeployment) GetNFDeploymentStatus() *nephiov1alpha1.NFDeploymentStatus {
	return d.Status.NFDeploymentStatus.DeepCopy()
}
func (d *PCRFDeployment) SetNFDeploymentSpec(s *nephiov1alpha1.NFDeploymentSpec) {
	s.DeepCopyInto(&d.Spec.NFDeploymentSpec)
}
func (d *PCRFDeployment) SetNFDeploymentStatus(s *nephiov1alpha1.NFDeploymentStatus) {
	s.DeepCopyInto(&d.Status.NFDeploymentStatus)
}

// Interface type metadata.
var (
	PCRFDeploymentKind             = reflect.TypeOf(PCRFDeployment{}).Name()
	PCRFDeploymentGroupKind        = schema.GroupKind{Group: nephiov1alpha1.Group, Kind: PCRFDeploymentKind}.String()
	PCRFDeploymentKindAPIVersion   = PCRFDeploymentKind + "." + nephiov1alpha1.GroupVersion.String()
	PCRFDeploymentGroupVersionKind = nephiov1alpha1.GroupVersion.WithKind(PCRFDeploymentKind)
)

func init() {
	nephiov1alpha1.SchemeBuilder.Register(&PCRFDeployment{}, &PCRFDeploymentList{})
}
