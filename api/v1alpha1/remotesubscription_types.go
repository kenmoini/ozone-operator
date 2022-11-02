/*
Copyright 2022 Ken Moini.

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

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RemoteSubscriptionSpec defines the desired state of RemoteSubscription
type RemoteSubscriptionSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// RemoteCluster defines the remote cluster and the connection to it
	RemoteCluster RemoteCluster `json:"remoteCluster"`

	// Operator defines the local operator to be deployed to the remote cluster
	Operator Operator `json:"operator"`
}

// RemoteCluster defines the remote cluster and the connection to it
type RemoteCluster struct {
	// Name is the name of the remote cluster as found in RHACM
	Name string `json:"name"`
	// RefreshInterval is the interval at which the remote cluster will be refreshed and synced
	RefreshInterval int `json:"refreshInterval,omitempty"`
}

// Operator defines the local operator to be deployed to the remote cluster
type Operator struct {
	// PackageName is the name of the operator PackageManifest
	PackageName string `json:"packageName"`
	// PackageNamespace is the namespace of the operator PackageManifest
	PackageNamespace string `json:"packageNamespace"`
	// Channel is the name of the operator channel - if this is not specified, the Operator default channel will be used
	Channel string `json:"channel,omitempty"`
	// StartingCSV is the name of the operator CSV to start with - if this is not specified, the latest CSV in the channel will be used
	StartingCSV string `json:"startingCSV,omitempty"`
	// InstallPlanApproval is the approval strategy for the operator install plan
	InstallPlanApproval string `json:"installPlanApproval,omitempty"`
	// Source is the name of the CatalogSource that contains the operator
	Source string `json:"source"`
	// SourceNamespace is the namespace of the CatalogSource that contains the operator
	SourceNamespace string `json:"sourceNamespace"`
	// InstallMode is the install mode for the operator, options include `all-namespaces`, `single-namespace`, `multi-namespace`, or `own-namespace` - if this is not specified, the Operator's default Installation Mode will be used.
	InstallMode string `json:"installMode,omitempty"`
	// TargetNamespace is the namespace in the remote cluster where the operator will be deployed
	TargetNamespace string `json:"targetNamespace,omitempty"`
}

// RemoteSubscriptionStatus defines the observed state of RemoteSubscription
type RemoteSubscriptionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// RemoteClusterStatus is the status of the remote cluster and the connection to it
	RemoteClusterStatus RemoteClusterStatus `json:"remoteClusterStatus"`

	// OperatorStatus is the status of the operator deployment to the remote cluster
	OperatorStatus OperatorStatus `json:"operatorStatus"`

	// Conditions is a list of conditions and their status
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

// RemoteClusterStatus is the status of the remote cluster and the connection to it
type RemoteClusterStatus struct {
	// Status is the status of the remote cluster
	Status string `json:"status"`
	// Message is a message about the status of the remote cluster
	Message string `json:"message"`
}

// OperatorStatus is the status of the operator deployment to the remote cluster
type OperatorStatus struct {
	// Status is the status of the operator deployment
	Status string `json:"status"`
	// Message is a message about the status of the operator deployment
	Message string `json:"message"`
}

// Condition is a condition and its status
type Condition struct {
	// Type is the type of the condition
	Type string `json:"type"`
	// Status is the status of the condition
	Status string `json:"status"`
	// Message is a message about the status of the condition
	Message string `json:"message"`
	// LastTransitionTime is the last time the condition transitioned from one status to another
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`
	// Reason is the reason for the condition's last transition
	Reason string `json:"reason"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RemoteSubscription is the Schema for the remotesubscriptions API
type RemoteSubscription struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RemoteSubscriptionSpec   `json:"spec,omitempty"`
	Status RemoteSubscriptionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RemoteSubscriptionList contains a list of RemoteSubscription
type RemoteSubscriptionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RemoteSubscription `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RemoteSubscription{}, &RemoteSubscriptionList{})
}
