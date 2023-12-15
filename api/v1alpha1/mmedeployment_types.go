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
	nephiov1alpha1 "github.com/nephio-project/api/nf_deployments/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"reflect"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MMEDeploymentSpec defines the desired state of MMEDeployment
type MMEDeploymentSpec struct {
	nephiov1alpha1.NFDeploymentSpec `json:",inline" yaml:",inline"`
}

// MMEDeploymentStatus defines the observed state of MMEDeployment
type MMEDeploymentStatus struct {
	nephiov1alpha1.NFDeploymentStatus `json:",inline" yaml:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MMEDeployment is the Schema for the mmedeployments API
type MMEDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MMEDeploymentSpec   `json:"spec,omitempty"`
	Status MMEDeploymentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MMEDeploymentList contains a list of MMEDeployment
type MMEDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MMEDeployment `json:"items"`
}

// Implement NFDeployment interface

func (d *MMEDeployment) GetNFDeploymentSpec() *nephiov1alpha1.NFDeploymentSpec {
	return d.Spec.NFDeploymentSpec.DeepCopy()
}
func (d *MMEDeployment) GetNFDeploymentStatus() *nephiov1alpha1.NFDeploymentStatus {
	return d.Status.NFDeploymentStatus.DeepCopy()
}
func (d *MMEDeployment) SetNFDeploymentSpec(s *nephiov1alpha1.NFDeploymentSpec) {
	s.DeepCopyInto(&d.Spec.NFDeploymentSpec)
}
func (d *MMEDeployment) SetNFDeploymentStatus(s *nephiov1alpha1.NFDeploymentStatus) {
	s.DeepCopyInto(&d.Status.NFDeploymentStatus)
}

// Interface type metadata.
var (
	MMEDeploymentKind             = reflect.TypeOf(MMEDeployment{}).Name()
	MMEDeploymentGroupKind        = schema.GroupKind{Group: nephiov1alpha1.Group, Kind: MMEDeploymentKind}.String()
	MMEDeploymentKindAPIVersion   = MMEDeploymentKind + "." + nephiov1alpha1.GroupVersion.String()
	MMEDeploymentGroupVersionKind = nephiov1alpha1.GroupVersion.WithKind(MMEDeploymentKind)
)

func init() {
	nephiov1alpha1.SchemeBuilder.Register(&MMEDeployment{}, &MMEDeploymentList{})
}
