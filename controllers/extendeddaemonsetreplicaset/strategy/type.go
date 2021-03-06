// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

package strategy

import (
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	corev1 "k8s.io/api/core/v1"

	datadoghqv1alpha1 "github.com/DataDog/extendeddaemonset/api/v1alpha1"
)

// ReplicaSetStatus repesent the status of a ReplicaSet
type ReplicaSetStatus string

const (
	// ReplicaSetStatusActive the ReplicaSet is currently active
	ReplicaSetStatusActive ReplicaSetStatus = "active"
	// ReplicaSetStatusCanary the ReplicaSet is currently in canary mode
	ReplicaSetStatusCanary ReplicaSetStatus = "canary"
	// ReplicaSetStatusUnknown the controller is not able to define the ReplicaSet status
	ReplicaSetStatusUnknown ReplicaSetStatus = "unknown"
)

// Parameters use to store all the parameter need to a strategy
type Parameters struct {
	MinPodUpdate int32
	MaxPodUpdate int32

	EDSName          string
	Strategy         *datadoghqv1alpha1.ExtendedDaemonSetSpecStrategy
	Replicaset       *datadoghqv1alpha1.ExtendedDaemonSetReplicaSet
	ReplicaSetStatus string

	NewStatus *datadoghqv1alpha1.ExtendedDaemonSetReplicaSetStatus

	CanaryNodes []string

	NodeByName      map[string]*NodeItem
	PodByNodeName   map[*NodeItem]*corev1.Pod
	PodToCleanUp    []*corev1.Pod
	UnscheduledPods []*corev1.Pod

	Logger logr.Logger
}

// Result information returns by a strategy
type Result struct {
	// PodsToCreate list of NodeItem for Pods creation
	PodsToCreate []*NodeItem
	// PodsToDelete list of NodeItem for Pods deletion
	PodsToDelete []*NodeItem

	UnscheduledNodesDueToResourcesConstraints []string

	NewStatus *datadoghqv1alpha1.ExtendedDaemonSetReplicaSetStatus
	Result    reconcile.Result
}

// NodeList list of NodeItem
type NodeList struct {
	Items []*NodeItem
}

// NodeItem used to store all informations needs to create or delete a pod
type NodeItem struct {
	Node                     *corev1.Node
	ExtendedDaemonsetSetting *datadoghqv1alpha1.ExtendedDaemonsetSetting
}

// NewNodeItem used to create new NodeItem instance
func NewNodeItem(node *corev1.Node, edsNode *datadoghqv1alpha1.ExtendedDaemonsetSetting) *NodeItem {
	return &NodeItem{
		Node:                     node,
		ExtendedDaemonsetSetting: edsNode,
	}
}
