/*
Copyright © 2019 The controller101 Authors.

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

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type VirtualMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMachineSpec   `json:"spec"`
	Status VirtualMachineStatus `json:"status"`
}

type VirtualMachineResource struct {
	CPU      int32 `json:"cpu"`
	Memory   int32 `json:"memory"`
	RootDisk int32 `json:"rootDisk"`
}

type VirtualMachineSpec struct {
	Action   string                 `json:"action"`
	Resource VirtualMachineResource `json:"resource"`
}

type VirtualMachinePhase string

const (
	VirtualMachineNone          VirtualMachinePhase = ""
	VirtualMachinePending       VirtualMachinePhase = "Pending"
	VirtualMachineSynchronizing VirtualMachinePhase = "Synchronizing"
	VirtualMachineSynchronized  VirtualMachinePhase = "Synchronized"
	VirtualMachineFailed        VirtualMachinePhase = "Failed"
	VirtualMachineTerminating   VirtualMachinePhase = "Terminating"
	VirtualMachineUnknown       VirtualMachinePhase = "Unknown"
)

type ResourceUsage struct {
	CPU    int32 `json:"cpu"`
	Memory int32 `json:"memory"`
}

type ServerStatus struct {
	State string        `json:"state"`
	Usage ResourceUsage `json:"usage"`
}

type VirtualMachineStatus struct {
	Phase          VirtualMachinePhase `json:"phase"`
	Reason         string              `json:"reason,omitempty"`
	Server         ServerStatus        `json:"server,omitempty"`
	LastUpdateTime metav1.Time         `json:"lastUpdateTime"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type VirtualMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []VirtualMachine `json:"items"`
}
