package handler

import (
	"github.com/labstack/echo/v4"
	"k8s.io/client-go/kubernetes"
)

type Handler struct {
	Client *kubernetes.Clientset
}

func NewHandler(clientSet *kubernetes.Clientset) *Handler {
	return &Handler{Client: clientSet}
}

func (h *Handler) CreateApp(c echo.Context) error {
	return nil
}
