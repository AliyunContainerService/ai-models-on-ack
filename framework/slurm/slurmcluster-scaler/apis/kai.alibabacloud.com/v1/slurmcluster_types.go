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

package v1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SlurmClusterSpec defines the desired state of SlurmCluster
type SlurmClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// HeadGroupSpecs are the spec for the slurmctld pod
	Slurmctld  SlurmctldSpec   `json:"slurmctld"`
	Slurmdbd   *SlurmdbdSpec   `json:"slurmdbd,omitempty"`
	Slurmrestd *SlurmrestdSpec `json:"slurmrestd,omitempty"`
	// WorkerGroupSpecs are the specs for the slurmd pods
	WorkerGroupSpecs []WorkerGroupSpec `json:"workerGroupSpecs,omitempty"`
	//
	MungeKeyPath  string `json:"mungeConfPath,omitempty"`
	SlurmConfPath string `json:"slurmConfPath,omitempty"`
}

type SlurmdbdSpec struct {
	// Template is a pod template for the worker
	// The fitst container will be considered as the slurmrestd
	// The operator will overwrite the command of the slurmrestd container
	Template v1.PodTemplateSpec `json:"template"`
}

type SlurmrestdSpec struct {
	Replicas *int32 `json:"replicas,omitempty"`
	// Template is a pod template for the SlurmRestd
	// The fitst container will be considered as the slurmrestd
	// The operator will overwrite the command of the slurmrestd container
	Template v1.PodTemplateSpec `json:"template"`
}

// Slurmctld are the spec for the head pod
type SlurmctldSpec struct {
	// Template is a pod template for the worker
	// The fitst container will be considered as the slurmctld
	// The operator will overwrite the command of the slurmctld container
	// Do not set hostname in the template, hostname will be set automatically
	Template v1.PodTemplateSpec `json:"template"`
	// EnablePrometheus indicates whether to enable prometheus metrics
	// If true, please make sure your image is built with exporter exported
	EnablePrometheus bool `json:"enablePrometheus,omitempty"`
	// Replicas indicates how many slurmctld pod should be created
	// see https://slurm.schedmd.com/quickstart_admin.html#HA to find more
	// introduction
	Replicas *int32 `json:"replicas,omitempty"`
	// InitBash will be executed before any command in template.
	// Template can be used in init bash with SlurmClusterSpec to be the param.
	// For example, you can use {{.MungeKeyPath}} to get munge key path in slurm
	// cluster spec.
	// The command for pod will be: {{InitBash}}\n{{Command Args}}\n{{Sleep Inf}}
	InitBash string `json:"initBash,omitempty"`
}

// WorkerGroupSpec are the specs for the worker pods
type WorkerGroupSpec struct {
	// we can have multiple worker groups, we distinguish them by name
	// +kubebuilder:validation:Required
	GroupName string `json:"groupName"`
	// Replicas is the number of desired Pods for this worker group.
	// +kubebuilder:default:=0
	Replicas              int32    `json:"replicas,omitempty"`
	PodsToBeDeleted       []string `json:"podsToBeDeleted,omitempty"`
	PodIndexesToBeCreated []int32  `json:"podsToBeCreated,omitempty"`
	// Template is a pod template for the worker
	// The fitst container will be considered as the slurmd
	// The operator will overwrite the command of the slurmd container
	Template v1.PodTemplateSpec `json:"template"`
	// InitBash will be executed before any command in template.
	// Template can be used in init bash with SlurmClusterSpec to be the param.
	// For example, you can use {{.MungeKeyPath}} to get munge key path in slurm
	// cluster spec.
	// The command for pod will be: {{InitBash}}\n{{Command Args}}\n{{Sleep Inf}}
	InitBash string `json:"initBash,omitempty"`
}

// The overall state of the Slurm cluster.
type ClusterState string

const (
	Ready       ClusterState = "ready"
	Unhealthy   ClusterState = "unhealthy"
	Failed      ClusterState = "failed"
	Progressing ClusterState = "progressing"
)

// SlurmClusterStatus defines the observed state of SlurmCluster
type SlurmClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Status reflects the status of the cluster
	State ClusterState `json:"state,omitempty"`
	// AvailableWorkerReplicas indicates how many replicas are available in the cluster
	AvailableWorkerReplicas int32    `json:"availableWorkerReplicas,omitempty"`
	AvailableWorkers        []string `json:"availableWorkers,omitempty"`
	// LastUpdateTime indicates last update timestamp for this cluster status.
	// +nullable
	LastUpdateTime *metav1.Time `json:"lastUpdateTime,omitempty"`
	// Reason provides more information about current State
	Reason string `json:"reason,omitempty"`
	// observedGeneration is the most recent generation observed for this SlurmCluster. It corresponds to the
	// SlurmCluster's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="available workers",type=integer,JSONPath=".status.availableWorkerReplicas",priority=0
// +kubebuilder:printcolumn:name="status",type="string",JSONPath=".status.state",priority=0
// +kubebuilder:printcolumn:name="age",type="date",JSONPath=".metadata.creationTimestamp",priority=0
// +groupName=kai.alibabacloud.com
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// SlurmCluster is the Schema for the slurmclusters API
type SlurmCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SlurmClusterSpec   `json:"spec,omitempty"`
	Status SlurmClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// SlurmClusterList contains a list of SlurmCluster
type SlurmClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SlurmCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SlurmCluster{}, &SlurmClusterList{})
}
