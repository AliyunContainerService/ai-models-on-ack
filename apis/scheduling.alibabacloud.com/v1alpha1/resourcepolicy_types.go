/*
Copyright 2021 Alibaba Cloud.

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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type Policy string
type TryNextUnitsOptions struct {
	// +kubebuilder:default:=LackResourceOrExceedMax
	// +kubebuilder:validation:Enum={"TimeoutOrExceedMax","ExceedMax","LackResourceOrExceedMax","LackResourceAndNoterminating","Timeout"}
	Policy Policy `json:"policy,omitempty" protobuf:"bytes,1,opt,name=policy"`
	// if policy is TimeoutOrExceedMax and timeout is not set, we will see this timeout as 15 min
	Timeout *metav1.Duration `json:"timeout,omitempty" protobuf:"bytes,2,opt,name=timeout"`
}

var (
	TimeoutOrExceedMax           Policy = "TimeoutOrExceedMax"
	LackResourceOrExceedMax      Policy = "LackResourceOrExceedMax"
	LackResourceAndNoterminating Policy = "LackResourceAndNoterminating"
	ExceedMax                    Policy = "ExceedMax"
	Timeout                      Policy = "Timeout"
)

// ResourcePolicySpec defines the desired state of ResourcePolicy
type ResourcePolicySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	ConsumerRef      *v1.ObjectReference               `json:"consumerRef,omitempty" protobuf:"bytes,1,opt,name=consumerRef"`
	Strategy         SchedulingStrategy                `json:"strategy,omitempty" protobuf:"bytes,2,opt,name=strategy"`
	Units            []Unit                            `json:"units,omitempty" protobuf:"bytes,3,opt,name=units"`
	Selector         map[string]string                 `json:"selector,omitempty" protobuf:"bytes,4,rep,name=selector"`
	LabelExpressions []metav1.LabelSelectorRequirement `json:"matchExpressions,omitempty" protobuf:"bytes,5,rep,name=matchExpressions"`
	WhenTryNextUnits TryNextUnitsOptions               `json:"whenTryNextUnits,omitempty" protobuf:"bytes,6,opt,name=whenTryNextUnits"`
	// +kubebuilder:default:=false
	IgnorePreviousPod bool `json:"ignorePreviousPod,omitempty" protobuf:"bytes,7,opt,name=ignorePreviousPod"`
	// +kubebuilder:default:=true
	IgnoreTerminatingPod bool `json:"ignoreTerminatingPod,omitempty" protobuf:"bytes,8,opt,name=ignoreTerminatingPod"`
	// +kubebuilder:default:=AfterAllUnits
	// +kubebuilder:validation:Enum:={BeforeNextUnit,AfterAllUnits}
	PreemptPolicy PreemptPolicy `json:"preemptPolicy,omitempty" protobuf:"bytes,9,opt,name=preemptPolicy"`
	//
	MatchLabelKeys []string `json:"matchLabelKeys,omitempty" protobuf:"bytes,10,rep,name=matchLabelKeys"`
}

// ResourcePolicyStatus defines the observed state of ResourcePolicy
type ResourcePolicyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ResourcePolicy is the Schema for the resourcepolicies API
type ResourcePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourcePolicySpec   `json:"spec,omitempty"`
	Status ResourcePolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ResourcePolicyList contains a list of ResourcePolicy
type ResourcePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ResourcePolicy `json:"items"`
}

type SchedulingStrategy string

const (
	Prefer   SchedulingStrategy = "prefer"
	Required SchedulingStrategy = "required"
)

type ResourceType string

const (
	ECI ResourceType = "eci"
	ECS ResourceType = "ecs"
	ACS ResourceType = "acs"
)

type PreemptPolicy string

const BeforeNextUnit PreemptPolicy = "BeforeNextUnit"
const AfterAllUnits PreemptPolicy = "AfterAllUnits"

type Unit struct {
	Resource       ResourceType      `json:"resource,omitempty" protobuf:"bytes,1,opt,name=resource"`
	Max            *int32            `json:"max,omitempty" protobuf:"varint,2,opt,name=max"`
	NodeSelector   map[string]string `json:"nodeSelector,omitempty" protobuf:"bytes,3,rep,name=nodeSelector"`
	SpotInstance   bool              `json:"spotInstance,omitempty" protobuf:"varint,4,opt,name=spotInstance"`
	ScaleInWeight  *int32            `json:"scaleInWeight,omitempty" protobuf:"varint,5,opt,name=scaleInWeight"`
	PodLabels      map[string]string `json:"podLabels,omitempty" protobuf:"bytes,6,rep,name=podLabels"`
	PodAnnotations map[string]string `json:"podAnnotations,omitempty" protobuf:"bytes,7,rep,name=podAnnotations"`
}

func init() {
	SchemeBuilder.Register(&ResourcePolicy{}, &ResourcePolicyList{})
}
