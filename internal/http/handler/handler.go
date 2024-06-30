package handler

import (
	"github.com/labstack/echo/v4"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	return nil
}
