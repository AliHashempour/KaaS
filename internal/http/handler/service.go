package handler

import (
	"Kaas/internal/model"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"net/http"
)

type Service struct {
	Client *kubernetes.Clientset
}

func NewService(clientSet *kubernetes.Clientset) *Service {
	return &Service{Client: clientSet}
}

func (s *Service) DeployPostgres(c echo.Context) error {
	var req model.PostgresSpec
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	// Create Secret
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.AppName + "-secret",
			Namespace: "default",
			Labels:    map[string]string{"app": req.AppName},
		},
		StringData: map[string]string{
			"username": "postgres",
			"password": "postgres",
		},
	}
	_, err := s.Client.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to create secret: %v", err))
	}

	// Create ConfigMap
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.AppName + "-config",
			Namespace: "default",
			Labels:    map[string]string{"app": req.AppName},
		},
		Data: map[string]string{
			"max_connections": "100",
		},
	}
	_, err = s.Client.CoreV1().ConfigMaps("default").Create(context.TODO(), configMap, metav1.CreateOptions{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to create config map: %v", err))
	}

	// Create Service
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.AppName + "-service",
			Namespace: "default",
			Labels:    map[string]string{"app": req.AppName},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port: 5432,
				},
			},
			Selector: map[string]string{
				"app": req.AppName,
			},
		},
	}
	_, err = s.Client.CoreV1().Services("default").Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to create service: %v", err))
	}
	// Create StatefulSet
	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.AppName,
			Namespace: "default",
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": req.AppName},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": req.AppName},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "postgres",
							Image: "postgres:latest",
							Ports: []corev1.ContainerPort{{ContainerPort: 5432}},
							Env: []corev1.EnvVar{
								{Name: "POSTGRES_USER", ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{Name: req.AppName + "-secret"},
										Key:                  "username",
									},
								}},
								{Name: "POSTGRES_PASSWORD", ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{Name: req.AppName + "-secret"},
										Key:                  "password",
									},
								}},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse(req.Resources.CPU),
									corev1.ResourceMemory: resource.MustParse(req.Resources.RAM),
								},
							},
						},
					},
				},
			},
		},
	}
	_, err = s.Client.AppsV1().StatefulSets("default").Create(context.TODO(), statefulSet, metav1.CreateOptions{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to create stateful set: %v", err))
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "PostgreSQL deployed successfully"})
}
