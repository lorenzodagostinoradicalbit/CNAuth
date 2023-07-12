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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// JWTKeySpec defines the desired state of JWTKey
type JWTKeySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of JWTKey. Edit jwtkey_types.go to remove/update
	Key string `json:"key,omitempty"`
}

// JWTKeyStatus defines the observed state of JWTKey
type JWTKeyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Key string `json:"key,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// JWTKey is the Schema for the jwtkeys API
type JWTKey struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JWTKeySpec   `json:"spec,omitempty"`
	Status JWTKeyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// JWTKeyList contains a list of JWTKey
type JWTKeyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JWTKey `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JWTKey{}, &JWTKeyList{})
}
