package model

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AppStatus struct {
	DeploymentName string
	Replicas       int32
	ReadyReplicas  int32
	PodStatuses    []PodStatus
}

type PodStatus struct {
	Name      string
	Phase     corev1.PodPhase
	HostIP    string
	PodIP     string
	StartTime *metav1.Time
}
