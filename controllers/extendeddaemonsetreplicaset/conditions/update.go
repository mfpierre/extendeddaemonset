// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

package conditions

import (
	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	datadoghqv1alpha1 "github.com/DataDog/extendeddaemonset/api/v1alpha1"
)

// NewExtendedDaemonSetReplicaSetCondition returns new ExtendedDaemonSetReplicaSetCondition instance
func NewExtendedDaemonSetReplicaSetCondition(conditionType datadoghqv1alpha1.ExtendedDaemonSetReplicaSetConditionType, conditionStatus corev1.ConditionStatus, now metav1.Time, reason, message string, supportLastUpdate bool) datadoghqv1alpha1.ExtendedDaemonSetReplicaSetCondition {
	return datadoghqv1alpha1.ExtendedDaemonSetReplicaSetCondition{
		Type:               conditionType,
		Status:             conditionStatus,
		LastUpdateTime:     now,
		LastTransitionTime: now,
		Reason:             reason,
		Message:            message,
	}
}

// UpdateExtendedDaemonSetReplicaSetStatusCondition used to update a specific ExtendedDaemonSetReplicaSetConditionType
func UpdateExtendedDaemonSetReplicaSetStatusCondition(status *datadoghqv1alpha1.ExtendedDaemonSetReplicaSetStatus, now metav1.Time, t datadoghqv1alpha1.ExtendedDaemonSetReplicaSetConditionType, conditionStatus corev1.ConditionStatus, desc string, writeFalseIfNotExist, supportLastUpdate bool) {
	idCondition := getIndexForConditionType(status, t)
	if idCondition >= 0 {
		if status.Conditions[idCondition].Status != conditionStatus {
			status.Conditions[idCondition].LastTransitionTime = now
			status.Conditions[idCondition].Status = conditionStatus
			status.Conditions[idCondition].LastUpdateTime = now
		}
		if supportLastUpdate {
			status.Conditions[idCondition].LastUpdateTime = now
		}
		status.Conditions[idCondition].Message = desc
	} else if conditionStatus == corev1.ConditionTrue || writeFalseIfNotExist {
		// Only add if the condition is True
		status.Conditions = append(status.Conditions, NewExtendedDaemonSetReplicaSetCondition(t, conditionStatus, now, "", desc, supportLastUpdate))
	}
}

// UpdateErrorCondition used to update the ExtendedDaemonSetReplicaSet status error condition
func UpdateErrorCondition(status *datadoghqv1alpha1.ExtendedDaemonSetReplicaSetStatus, now metav1.Time, err error, desc string) {
	if err != nil {
		UpdateExtendedDaemonSetReplicaSetStatusCondition(status, now, datadoghqv1alpha1.ConditionTypeReconcileError, corev1.ConditionTrue, desc, false, true)
	} else {
		UpdateExtendedDaemonSetReplicaSetStatusCondition(status, now, datadoghqv1alpha1.ConditionTypeReconcileError, corev1.ConditionFalse, desc, false, true)
	}
}

func getIndexForConditionType(status *datadoghqv1alpha1.ExtendedDaemonSetReplicaSetStatus, t datadoghqv1alpha1.ExtendedDaemonSetReplicaSetConditionType) int {
	idCondition := -1
	if status == nil {
		return idCondition
	}
	for i, condition := range status.Conditions {
		if condition.Type == t {
			idCondition = i
			break
		}
	}
	return idCondition
}

// GetExtendedDaemonSetReplicaSetStatusCondition return the condition struct corresponding to the ExtendedDaemonSetReplicaSetConditionType provided in argument.
// return nil if not found
func GetExtendedDaemonSetReplicaSetStatusCondition(status *datadoghqv1alpha1.ExtendedDaemonSetReplicaSetStatus, t datadoghqv1alpha1.ExtendedDaemonSetReplicaSetConditionType) *datadoghqv1alpha1.ExtendedDaemonSetReplicaSetCondition {
	idCondition := getIndexForConditionType(status, t)
	if idCondition == -1 {
		return nil
	}
	return &status.Conditions[idCondition]
}
