package kube

import (
	"github.com/labstack/echo/v4"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"net/http"
	"os"
	"path/filepath"
)

func GetKubeConfig() (*kubernetes.Clientset, error) {

	// Try to use in-cluster configuration first
	config, err := rest.InClusterConfig()
	if err != nil {
		// Fall back to out-of-cluster configuration if in-cluster fails
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to find home directory: "+err.Error())
		}
		kubeConfigPath := filepath.Join(homeDir, ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to configure Kubernetes client: "+err.Error())
		}
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientSet, nil
}
