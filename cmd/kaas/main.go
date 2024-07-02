package main

import (
	"Kaas/internal/http/handler"
	"Kaas/internal/http/middleware"
	"Kaas/internal/kube"
	"github.com/labstack/echo/v4"
	"log"
)

func main() {
	e := echo.New()
	e.Use(middleware.InfoLogger)

	clientSet, err := kube.GetKubeConfig()
	if err != nil {
		log.Fatalf("Failed to get Kubernetes config: %v", err)
	}

	applicationHandler := handler.NewApplication(clientSet)
	serviceHandler := handler.NewService(clientSet)

	e.GET("/", applicationHandler.GetNodes)

	e.POST("/create", applicationHandler.CreateApp)

	e.GET("/status/:appName", applicationHandler.GetDeploymentStatus)

	e.GET("/status", applicationHandler.GetAllDeploymentsStatus)

	e.POST("/postgres", serviceHandler.DeployPostgres)

	e.Logger.Fatal(e.Start(":8080"))
}
