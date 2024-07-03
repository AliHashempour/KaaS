package model

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
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

type MonitorStatus struct {
	ID           int       `json:"id"`
	AppName      string    `json:"app_name"`
	FailureCount int       `json:"failure_count"`
	SuccessCount int       `json:"success_count"`
	LastFailure  time.Time `json:"last_failure"`
	LastSuccess  time.Time `json:"last_success"`
	CreatedAt    time.Time `json:"created_at"`
}
