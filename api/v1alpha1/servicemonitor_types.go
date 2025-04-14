/*
Copyright 2025.

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

// ServiceMonitorSpec defines the desired state of ServiceMonitor
type ServiceMonitorSpec struct {
	// TargetDomain is the ExternalName to monitor in services
	// +kubebuilder:validation:Required
	TargetDomain string `json:"targetDomain"`
}

// ServiceMonitorStatus defines the observed state of ServiceMonitor
type ServiceMonitorStatus struct {
	// Count is the number of services with matching ExternalName
	Count int32 `json:"count"`

	// LastCheckTime is the last time the services were checked
	LastCheckTime *metav1.Time `json:"lastCheckTime,omitempty"`

	// Conditions represent the latest available observations of the ServiceMonitor's state
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Count",type="integer",JSONPath=".status.count"
//+kubebuilder:printcolumn:name="LastCheck",type="date",JSONPath=".status.lastCheckTime"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// ServiceMonitor is the Schema for the servicemonitors API
type ServiceMonitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServiceMonitorSpec   `json:"spec,omitempty"`
	Status ServiceMonitorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ServiceMonitorList contains a list of ServiceMonitor
type ServiceMonitorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServiceMonitor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ServiceMonitor{}, &ServiceMonitorList{})
}
