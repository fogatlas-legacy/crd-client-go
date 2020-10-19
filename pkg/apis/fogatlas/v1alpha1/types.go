/*
Modification:
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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type FAStatus int32

const (
	New FAStatus = iota
	Synced
	Failed
	Changed
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FADepl is a specification for a FADepl resource
type FADepl struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FADeplSpec   `json:"spec"`
	Status FADeplStatus `json:"status"`
}

type FADeplMicroservice struct {
	Name         string            `json:"name"`
	Regions      []*FARegion       `json:"regions,omitempty"`
	MIPSRequired int64             `json:"mipsrequired,omitempty"`
	Deployment   appsv1.Deployment `json:"deployment"`
}

type FARegion struct {
	RegionRequired string `json:"regionrequired,omitempty"`
	RegionSelected string `json:"regionselected"`
	Replicas       int32  `json:"replicas,omitempty"`
	Image          string `json:"image,omitempty"`
	CPU2MIPSMilli  int64  `json:"cpu2mipsmilli,omitempty"`
}

type FADeplDataFlow struct {
	BandwidthRequired int32  `json:"bandwidthrequired"`
	LatencyRequired   int32  `json:"latency"`
	SourceId          string `json:"sourceid"`
	DestinationId     string `json:"destinationid"`
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

type FAPlacement struct {
	Regions      []*FARegion `json:"regions"`
	Microservice string      `json:"microservice"`
}

type FALinkOccupancy struct {
	LinkId      string `json:"linkid"`
	BwAllocated int32  `json:"bwallocated"`
	IsChanged   bool   `json:"ischanged"`
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

// FedFADepl is a specification for a FedFADepl resource
type FedFADepl struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FedFADeplSpec `json:"spec"`
	Status FADeplStatus  `json:"status"`
}

type FedFADeplMicroservice struct {
	Name                string                    `json:"name"`
	FederatedDeployment unstructured.Unstructured `json:"federateddeployment"`
}

// FedFADeplSpec is the spec for a FADepl resource
type FedFADeplSpec struct {
	Name              string                   `json:"name"`
	Description       string                   `json:"description"`
	ExternalEndpoints []string                 `json:"externalendpoints"`
	Microservices     []*FedFADeplMicroservice `json:"microservices"`
	DataFlows         []*FADeplDataFlow        `json:"dataflows"`
	Algorithm         string                   `json:"algorithm"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FedFADeplList is a list of FedFADepl resources
type FedFADeplList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []FedFADepl `json:"items"`
}

//END OF FedFADepl

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Region is a specification for a Region resource
type Region struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RegionSpec   `json:"spec"`
	Status RegionStatus `json:"status"`
}

// RegionSpec is the spec for a Region resource
type RegionSpec struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Location    string     `json:"location"`
	Tier        int32      `json:"tier"`
	Type        RegionType `json:"type,omitempty"`
	CPUModel    string     `json:"cpumodel"`
	CPU2MIPS    int64      `json:"cpu2mips"`
}

// RegionType is the type of a Region resource
// nodes, clusters
type RegionType string

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
type ExternalEndpoint struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExternalEndpointSpec   `json:"spec"`
	Status ExternalEndpointStatus `json:"status"`
}

// ExternalEndpointSpec is the spec for a ExternalEndpoint resource
type ExternalEndpointSpec struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IpAddress   string `json:"ipaddress"`
	Location    string `json:"location"`
	RegionId    string `json:"regionid"`
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
type Link struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LinkSpec   `json:"spec"`
	Status LinkStatus `json:"status"`
}

// LinkSpec is the spec for a Link resource
type LinkSpec struct {
	Id         string `json:"id"`
	EndpointA  string `json:"endpointa"`
	EndpointB  string `json:"endpointb"`
	BwPeak     int32  `json:"bwpeak"`
	BwMeasured int32  `json:"bwmeasured"`
	Latency    int32  `json:"latency"`
	Status     string `json:"status"`
}

// LinkStatus is the status for a Link resource
type LinkStatus struct {
	BwAllocated   int32 `json:"bwallocated"`
	CurrentStatus int32 `json:"currentstatus"`
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
type DynamicNode struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DynamicNodeSpec   `json:"spec"`
	Status DynamicNodeStatus `json:"status"`
}

// DynamicNodeSpec is the spec for a DynamicNode resource
type DynamicNodeSpec struct {
	Id              string        `json:"id"`
	Description     string        `json:"description"`
	Location        string        `json:"location"`
	Tier            int32         `json:"tier"`
	IpAddress       string        `json:"ip_address"`
	RegionId        string        `json:"region_id"`
	AnsibleSshUser  string        `json:"ansible_ssh_user"`
	Architecture    string        `json:"architecture"`
	Processor       string        `json:"processor"`
	Gpu             string        `json:"gpu"`
	Memory          string        `json:"memory"`
	Storage         string        `json:"storage"`
	OperatingSystem string        `json:"operating_system"`
	Status          DynamicStatus `json:"status"`
}

type DynamicStatus string

const (
	Baremetal       DynamicStatus = "baremetal"
	Vpnset          DynamicStatus = "vpnset"
	Clustered       DynamicStatus = "clustered"
	Free_ext        DynamicStatus = "free_ext"
	Free            DynamicStatus = "free"
	Advertised      DynamicStatus = "advertised"
	Reserved_ext    DynamicStatus = "reserved_ext"
	Reserved        DynamicStatus = "reserved"
	Provisioned_ext DynamicStatus = "provisioned_ext"
	Provisioned     DynamicStatus = "provisioned"
	Provisioning    DynamicStatus = "provisioning"
	Unprovisioning  DynamicStatus = "unprovisioning"
	Clustering      DynamicStatus = "clustering"
	Unclustering    DynamicStatus = "unclustering"
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
