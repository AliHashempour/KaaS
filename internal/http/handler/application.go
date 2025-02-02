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

type Application struct {
	Client *kubernetes.Clientset
}

func NewApplication(clientSet *kubernetes.Clientset) *Application {
	return &Application{Client: clientSet}
}

func (h *Application) GetNodes(c echo.Context) error {
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

func (h *Application) CreateApp(c echo.Context) error {
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

func (h *Application) createSecret(app *model.Application) (*corev1.Secret, error) {
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

func (h *Application) createDeployment(app *model.Application) (*appsv1.Deployment, error) {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: app.AppName,
			Labels: map[string]string{
				"app":     app.AppName,
				"monitor": "true",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(app.Replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": app.AppName},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: app.AppName,
					Labels: map[string]string{
						"app":     app.AppName,
						"monitor": "true",
					},
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

func (h *Application) createService(app *model.Application) (*corev1.Service, error) {
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

func (h *Application) GetDeploymentStatus(c echo.Context) error {
	appName := c.Param("appName")
	if appName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Application name is required")
	}

	// Retrieve the deployment
	deployment, err := h.Client.AppsV1().Deployments("default").Get(context.TODO(), appName, metav1.GetOptions{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get deployment: "+err.Error())
	}

	// Retrieve pods associated with the deployment
	podList, err := h.Client.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", appName),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list pods: "+err.Error())
	}

	// Construct the response with deployment and pod details
	appStatus := model.AppStatus{
		DeploymentName: deployment.Name,
		Replicas:       *deployment.Spec.Replicas,
		ReadyReplicas:  deployment.Status.ReadyReplicas,
		PodStatuses:    []model.PodStatus{},
	}

	for _, pod := range podList.Items {
		ps := model.PodStatus{
			Name:      pod.Name,
			Phase:     pod.Status.Phase,
			HostIP:    pod.Status.HostIP,
			PodIP:     pod.Status.PodIP,
			StartTime: pod.Status.StartTime,
		}
		appStatus.PodStatuses = append(appStatus.PodStatuses, ps)
	}

	return c.JSON(http.StatusOK, appStatus)
}

func (h *Application) GetAllDeploymentsStatus(c echo.Context) error {
	deployments, err := h.Client.AppsV1().Deployments("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get deployments: "+err.Error())
	}

	var allStatuses []model.AppStatus

	for _, deployment := range deployments.Items {
		podList, err := h.Client.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{
			LabelSelector: fmt.Sprintf("app=%s", deployment.Name),
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to list pods for deployment "+deployment.Name+": "+err.Error())
		}

		appStatus := model.AppStatus{
			DeploymentName: deployment.Name,
			Replicas:       *deployment.Spec.Replicas,
			ReadyReplicas:  deployment.Status.ReadyReplicas,
		}

		for _, pod := range podList.Items {
			ps := model.PodStatus{
				Name:      pod.Name,
				Phase:     pod.Status.Phase,
				HostIP:    pod.Status.HostIP,
				PodIP:     pod.Status.PodIP,
				StartTime: pod.Status.StartTime,
			}
			appStatus.PodStatuses = append(appStatus.PodStatuses, ps)
		}
		allStatuses = append(allStatuses, appStatus)
	}

	return c.JSON(http.StatusOK, allStatuses)
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
