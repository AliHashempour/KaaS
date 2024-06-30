package handler

import (
	"Kaas/internal/model"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"net/http"
)

type Handler struct {
	Client *kubernetes.Clientset
}

func NewHandler(clientSet *kubernetes.Clientset) *Handler {
	return &Handler{Client: clientSet}
}

func (h *Handler) GetNodes(c echo.Context) error {
	nodes, err := h.Client.CoreV1().Nodes().List(c.Request().Context(), metav1.ListOptions{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get nodes: "+err.Error())
	}

	nodeNames := make([]string, len(nodes.Items))
	for i, node := range nodes.Items {
		nodeNames[i] = node.Name
	}

	return c.JSON(http.StatusOK, nodeNames)
}

func (h *Handler) CreateApp(c echo.Context) error {
	app := new(model.Application)
	if err := c.Bind(app); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	// Create Secrets if needed
	if hasSecrets(app.Envs) {
		secret, err := h.createSecret(app)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create secret: "+err.Error())
		}
		fmt.Println("Secret created: ", secret.Name) // Optionally log the secret creation
	}

	// Create Deployment
	_, err := h.createDeployment(app)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create deployment: "+err.Error())
	}

	// Create Service
	_, err = h.createService(app)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create service: "+err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Application created successfully"})
}

func (h *Handler) createSecret(app *model.Application) (*corev1.Secret, error) {
	secretData := make(map[string][]byte)
	for _, env := range app.Envs {
		if env.IsSecret {
			secretData[env.Key] = []byte(env.Value)
		}
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: app.AppName + "-secret",
		},
		Data: secretData,
	}
	return h.Client.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})
}

func (h *Handler) createDeployment(app *model.Application) (*appsv1.Deployment, error) {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: app.AppName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(app.Replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": app.AppName},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": app.AppName},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            app.AppName,
							Image:           app.ImageAddress + ":" + app.ImageTag,
							Ports:           []corev1.ContainerPort{{ContainerPort: int32(app.ServicePort)}},
							Env:             mapEnvVars(app),
							ImagePullPolicy: corev1.PullIfNotPresent,
						},
					},
				},
			},
		},
	}
	return h.Client.AppsV1().Deployments("default").Create(context.TODO(), deployment, metav1.CreateOptions{})
}

func (h *Handler) createService(app *model.Application) (*corev1.Service, error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: app.AppName + "-service",
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": app.AppName},
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       int32(app.ServicePort),
					TargetPort: intstr.IntOrString{IntVal: int32(app.ServicePort)},
				},
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}
	return h.Client.CoreV1().Services("default").Create(context.TODO(), service, metav1.CreateOptions{})
}

func int32Ptr(i int32) *int32 {
	return &i
}

func mapEnvVars(app *model.Application) []corev1.EnvVar {
	var vars []corev1.EnvVar
	for _, env := range app.Envs {
		if env.IsSecret {
			vars = append(vars, corev1.EnvVar{
				Name: env.Key,
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: app.AppName + "-secret",
						},
						Key: env.Key,
					},
				},
			})
		} else {
			vars = append(vars, corev1.EnvVar{
				Name:  env.Key,
				Value: env.Value,
			})
		}
	}
	return vars
}

func hasSecrets(envs []model.Env) bool {
	for _, env := range envs {
		if env.IsSecret {
			return true
		}
	}
	return false
}
