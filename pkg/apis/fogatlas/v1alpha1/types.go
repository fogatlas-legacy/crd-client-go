/*
Derived work:
Copyright 2019 FBK
This file has been modified in order to add resources needed by geok8s.

Original work:

Copyright 2017 The Kubernetes Authors.

License of both original work and derivative one:

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
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// FAStatus is a type for the possibile status values
type FAStatus int32

// Possible values of FAStaus
const (
	New FAStatus = iota
	Synced
	Failed
	Changed
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FADepl is a specification for a FADepl resource
// +kubebuilder:subresource:status
type FADepl struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FADeplSpec   `json:"spec"`
	Status FADeplStatus `json:"status,omitempty"`
}

// FADeplMicroservice represent a FADepl microservice
type FADeplMicroservice struct {
	Name         string            `json:"name"`
	Regions      []*FARegion       `json:"regions,omitempty"`
	MIPSRequired resource.Quantity `json:"mipsrequired,omitempty"`
	// +kubebuilder:validation:EmbeddedResource
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:validation:Type=object
	Deployment appsv1.Deployment `json:"deployment"`
}

// FARegion represents a region inside a FADepl
type FARegion struct {
	RegionRequired string `json:"regionrequired,omitempty"`
	RegionSelected string `json:"regionselected,omitempty"`
	Replicas       int32  `json:"replicas,omitempty"`
	Image          string `json:"image,omitempty"`
	CPU2MIPSMilli  int64  `json:"cpu2mipsmilli,omitempty"`
}

// FADeplDataFlow represent a data flow
type FADeplDataFlow struct {
	Name              string            `json:"name,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
	BandwidthRequired resource.Quantity `json:"bandwidthrequired"`
	LatencyRequired   resource.Quantity `json:"latency"`
	SourceID          string            `json:"sourceid"`
	DestinationID     string            `json:"destinationid"`
}

// FADeplSpec is the spec for a FADepl resource
type FADeplSpec struct {
	ExternalEndpoints []string              `json:"externalendpoints"`
	Microservices     []*FADeplMicroservice `json:"microservices"`
	DataFlows         []*FADeplDataFlow     `json:"dataflows,omitempty"`
	Algorithm         string                `json:"algorithm"`
}

// FADeplStatus is the status for a FADepl resource
type FADeplStatus struct {
	Placements     []*FAPlacement     `json:"placements,omitempty"`
	LinksOccupancy []*FALinkOccupancy `json:"linksoccupancy,omitempty"`
	CurrentStatus  FAStatus           `json:"currentstatus"`
}

// FAPlacement maps micrsoervices on regions
type FAPlacement struct {
	Regions      []*FARegion `json:"regions"`
	Microservice string      `json:"microservice"`
}

// FALinkOccupancy stores the link occupancy
type FALinkOccupancy struct {
	LinkID          string            `json:"linkid"`
	BwAllocated     resource.Quantity `json:"bwallocated"`
	PrevBwAllocated resource.Quantity `json:"prevbwallocated"`
	IsChanged       bool              `json:"ischanged"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FADeplList is a list of FADepl resources
type FADeplList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []FADepl `json:"items"`
}

//END OF FADepl

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FedFAApp is a specification for a FedFAApp resource
// +kubebuilder:subresource:status
type FedFAApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FedFAAppSpec `json:"spec"`
	Status FADeplStatus `json:"status"`
}

// FedFAAppChunk represents a chunk (piece) of a federated application
type FedFAAppChunk struct {
	Name            string                    `json:"name"`
	FederatedFADepl unstructured.Unstructured `json:"chunk"`
}

// FedFAAppSpec is the spec for a FADepl resource
type FedFAAppSpec struct {
	Name              string            `json:"name"`
	Description       string            `json:"description"`
	ApplicationChunks []*FedFAAppChunk  `json:"applicationchunks"`
	DataFlows         []*FADeplDataFlow `json:"c2cdataflows"`
	Algorithm         string            `json:"algorithm"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FedFAAppList is a list of FedFAApp resources
type FedFAAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []FedFAApp `json:"items"`
}

//END OF FedFAApp

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Region is a specification for a Region resource
// +kubebuilder:subresource:status
type Region struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec RegionSpec `json:"spec"`
	// RegionStatus is not used at the moment
	Status RegionStatus `json:"status,omitempty"`
}

// RegionSpec is the spec for a Region resource
type RegionSpec struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Location    string     `json:"location"`
	Tier        int32      `json:"tier"`
	Type        RegionType `json:"type,omitempty"`
	CPUModel    string     `json:"cpumodel"`
	CPU2MIPS    int64      `json:"cpu2mips"`
}

// RegionType is the type of a Region resource. Could be nodes, clusters, hostcluster
type RegionType string

// Possibile values of RegionType
const (
	Nodes       RegionType = "nodes"
	Clusters    RegionType = "clusters"
	Hostcluster RegionType = "hostcluster"
)

// RegionStatus is the status for a Region resource
type RegionStatus struct {
	CurrentStatus int32 `json:"currentStatus"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RegionList is a list of Region resources
type RegionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Region `json:"items"`
}

//END OF Region

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ExternalEndpoint is a specification for a ExternalEndpoint resource
// +kubebuilder:subresource:status
type ExternalEndpoint struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ExternalEndpointSpec `json:"spec"`
	// RegionStatus is not used at the moment
	Status ExternalEndpointStatus `json:"status,omitempty"`
}

// ExternalEndpointSpec is the spec for a ExternalEndpoint resource
type ExternalEndpointSpec struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IPAddress   string `json:"ipaddress"`
	Location    string `json:"location"`
	RegionID    string `json:"regionid"`
}

// ExternalEndpointStatus is the status for a ExternalEndpoint resource
type ExternalEndpointStatus struct {
	CurrentStatus int32 `json:"currentStatus"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ExternalEndpointList is a list of ExternalEndpoint resources
type ExternalEndpointList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ExternalEndpoint `json:"items"`
}

//END OF ExternalEndpoint

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Link is a specification for a Link resource
// +kubebuilder:subresource:status
type Link struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LinkSpec   `json:"spec"`
	Status LinkStatus `json:"status,omitempty"`
}

// LinkSpec is the spec for a Link resource
type LinkSpec struct {
	ID        string            `json:"id"`
	EndpointA string            `json:"endpointa"`
	EndpointB string            `json:"endpointb"`
	Bandwidth resource.Quantity `json:"bandwidth"`
	// in milliseconds
	Latency resource.Quantity `json:"latency"`
	Status  string            `json:"status"`
}

// LinkStatus is the status for a Link resource
type LinkStatus struct {
	BwAllocated   resource.Quantity `json:"bwallocated"`
	CurrentStatus int32             `json:"currentstatus"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LinkList is a list of Link resources
type LinkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Link `json:"items"`
}

//END OF Link

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DynamicNode is a specification for a DynamicNode resource
// +kubebuilder:subresource:status
type DynamicNode struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DynamicNodeSpec   `json:"spec"`
	Status DynamicNodeStatus `json:"status,omitempty"`
}

// DynamicNodeSpec is the spec for a DynamicNode resource
type DynamicNodeSpec struct {
	ID              string        `json:"id"`
	Description     string        `json:"description"`
	Location        string        `json:"location"`
	Tier            int32         `json:"tier"`
	IPAddress       string        `json:"ip_address"`
	RegionID        string        `json:"region_id"`
	AnsibleSSHUser  string        `json:"ansible_ssh_user"`
	Architecture    string        `json:"architecture"`
	Processor       string        `json:"processor"`
	Gpu             string        `json:"gpu"`
	Memory          string        `json:"memory"`
	Storage         string        `json:"storage"`
	OperatingSystem string        `json:"operating_system"`
	Status          DynamicStatus `json:"status"`
}

// DynamicStatus defined type represents the possibile statuses of a dynamic node
type DynamicStatus string

// Possible values of a dynamic node
const (
	Baremetal      DynamicStatus = "baremetal"
	Vpnset         DynamicStatus = "vpnset"
	Clustered      DynamicStatus = "clustered"
	FreeExt        DynamicStatus = "free_ext"
	Free           DynamicStatus = "free"
	Advertised     DynamicStatus = "advertised"
	ReservedExt    DynamicStatus = "reserved_ext"
	Reserved       DynamicStatus = "reserved"
	ProvisionedExt DynamicStatus = "provisioned_ext"
	Provisioned    DynamicStatus = "provisioned"
	Provisioning   DynamicStatus = "provisioning"
	Unprovisioning DynamicStatus = "unprovisioning"
	Clustering     DynamicStatus = "clustering"
	Unclustering   DynamicStatus = "unclustering"
)

// DynamicNodeStatus is the status for a DynamicNode resource
type DynamicNodeStatus struct {
	CurrentStatus int32 `json:"currentStatus"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DynamicNodeList is a list of DynamicNode resources
type DynamicNodeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []DynamicNode `json:"items"`
}

//END OF DynamicNode
