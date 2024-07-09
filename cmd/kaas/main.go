package main

import (
	"Kaas/internal/database"
	"Kaas/internal/http/handler"
	"Kaas/internal/http/middleware"
	"Kaas/internal/kube"
	"Kaas/internal/repository"
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

	db, err := database.InitializeDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	applicationHandler := handler.NewApplication(clientSet)
	serviceHandler := handler.NewService(clientSet)
	jobHandler := handler.NewJobHandler(repository.NewJobRepository(db))

	e.GET("/", applicationHandler.GetNodes)

	e.POST("application/create", applicationHandler.CreateApp)

	e.GET("application/status/:appName", applicationHandler.GetDeploymentStatus)

	e.GET("application/status", applicationHandler.GetAllDeploymentsStatus)

	e.POST("service/postgres", serviceHandler.DeployPostgres)

	e.GET("application/health/:appName", jobHandler.GetAppHealth)

	e.Logger.Fatal(e.Start(":8080"))
}
