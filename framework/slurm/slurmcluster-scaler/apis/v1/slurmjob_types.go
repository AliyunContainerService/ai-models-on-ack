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

// SlurmJobSpec defines the desired state of SlurmJob
type SlurmJobSpec struct {
	// ShutdownAfterJobFinishes will determine whether to delete the cluster once SlurmJob succeed or failed.
	// +kubebuilder:default:=false
	ShutdownAfterJobFinishes bool `json:"shutdownAfterJobFinishes,omitempty"`
	// TTLSecondsAfterFinished is the TTL to clean up SlurmCluster.
	// It's only working when ShutdownAfterJobFinishes set to true.
	TTLSecondsAfterFinished *int32 `json:"ttlSecondsAfterFinished,omitempty"`
	// SlurmClusterSpec is the cluster template to run the job
	SlurmClusterSpec *SlurmClusterSpec `json:"slurmClusterSpec,omitempty"`
	//
	SlurmCluster *v1.ObjectReference `json:"slurmCluster,omitempty"`
	// default is root
	RunAsUser string `json:"runAsUser,omitempty"`
	//
	StandOutputPath string `json:"standOutputPath,omitempty"`
	StandErrorPath  string `json:"standErrorPath,omitempty"`
	// suspend specifies whether the SlurmJob controller should create a SlurmCluster instance
	// If a job is applied with the suspend field set to true,
	// the SlurmCluster will not be created and will wait for the transition to false.
	// If the SlurmCluster is already created, it will be deleted.
	// In case of transition to false a new SlurmCluster will be created.
	// +kubebuilder:default:=false
	Suspend bool `json:"suspend,omitempty"`
	// +kubebuilder:validation:Required
	Command []string `json:"command,omitempty"`
	// RestartPolicy
	RestartPolicy v1.RestartPolicy `json:"restartPolicy,omitempty"`
}

type SubmitterTemplate struct {
}

// SlurmJobStatus defines the observed state of SlurmJob
type SlurmJobStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	ClusterName         string              `json:"clusterName,omitempty"`
	JobStatus           JobStatus           `json:"jobStatus,omitempty"`
	JobDeploymentStatus JobDeploymentStatus `json:"jobDeploymentStatus,omitempty"`
	Message             string              `json:"message,omitempty"`
	// Represents time when the job was acknowledged by the Slurm cluster.
	// It is not guaranteed to be set in happens-before order across separate operations.
	// It is represented in RFC3339 form
	StartTime *metav1.Time `json:"startTime,omitempty"`
	// Represents time when the job was ended.
	EndTime            *metav1.Time        `json:"endTime,omitempty"`
	SlurmClusterStatus *SlurmClusterStatus `json:"slurmClusterStatus,omitempty"`
	// observedGeneration is the most recent generation observed for this SlurmJob. It corresponds to the
	// SlurmJob's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// SlurmJob is the Schema for the slurmjobs API
type SlurmJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SlurmJobSpec   `json:"spec,omitempty"`
	Status SlurmJobStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SlurmJobList contains a list of SlurmJob
type SlurmJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SlurmJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SlurmJob{}, &SlurmJobList{})
}
