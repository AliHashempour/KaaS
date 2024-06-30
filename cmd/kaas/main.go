package main

import (
	"Kaas/internal/http/handler"
	"Kaas/internal/http/middleware"
	"Kaas/internal/kube"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func main() {
	e := echo.New()
	e.Use(middleware.InfoLogger)

	clientSet, err := kube.GetKubeConfig()
	if err != nil {
		log.Fatalf("Failed to get Kubernetes config: %v", err)
	}

	requestHandler := handler.NewHandler(clientSet)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Kaas service ")
	})

	e.POST("/create", requestHandler.CreateApp)

	e.Logger.Fatal(e.Start(":8080"))
}
