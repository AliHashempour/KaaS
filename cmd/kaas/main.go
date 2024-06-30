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

	requestHandler := handler.NewHandler(clientSet)

	e.GET("/", requestHandler.GetNodes)

	e.POST("/create", requestHandler.CreateApp)

	e.Logger.Fatal(e.Start(":8080"))
}
