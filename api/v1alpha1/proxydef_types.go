/*
Copyright 2024.

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

// Important: Run "make" to regenerate code after modifying this file

// ProxyDefSpec defines the desired state of ProxyDef
type ProxyDefSpec struct {
	HTTPProxy  string `json:"httpProxy,omitempty"`
	HTTPSProxy string `json:"httpsProxy,omitempty"`
	NoProxy    string `json:"noProxy,omitempty"`

	// TODO: Not implemented yet
	NoProxyCIDRs  string `json:"noProxyCidrs,omitempty"`
	SocksProxy    string `json:"socksProxy,omitempty"`
	FTPProxy      string `json:"ftpProxy,omitempty"`
	ProxyUser     string `json:"proxyUser,omitempty"`
	ProxyPassword string `json:"proxyPassword,omitempty"`
	ProxyProtocol string `json:"proxyProtocol,omitempty"`
	ProxyPort     int    `json:"proxyPort,omitempty"`
	NonProxyHosts string `json:"nonProxyHosts,omitempty"`
	AutoDetect    bool   `json:"autoDetect,omitempty"`
}

// ProxyDefStatus defines the observed state of ProxyDef
type ProxyDefStatus struct {
	// Represents the observations of a ProxyDef's current state.
	// ProxyDef.status.conditions.type are: "Ready", "Syncing", and "Degraded"
	// ProxyDef.status.conditions.status are one of True, False, Unknown.
	// ProxyDef.status.conditions.reason the value should be a CamelCase string and producers of specific
	// condition types may define expected values and meanings for this field, and whether the values
	// are considered a guaranteed API.
	// ProxyDef.status.conditions.Message is a human readable message indicating details about the transition.
	// For further information see: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// Conditions store the status conditions of the ProxyDef instances
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ProxyDef is the Schema for the proxydefs API
type ProxyDef struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProxyDefSpec   `json:"spec,omitempty"`
	Status ProxyDefStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ProxyDefList contains a list of ProxyDef
type ProxyDefList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProxyDef `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProxyDef{}, &ProxyDefList{})
}
