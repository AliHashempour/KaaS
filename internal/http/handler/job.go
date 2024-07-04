package handler

import (
	"Kaas/internal/repository"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

type JobHandler struct {
	Repo repository.JobRepository
}

func NewJobHandler(repo repository.JobRepository) *JobHandler {
	return &JobHandler{Repo: repo}
}

func (h *JobHandler) GetAppHealth(c echo.Context) error {
	appName := c.Param("appName")
	if appName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Application name is required")
	}

	status, err := h.Repo.GetAppHealth(appName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "No status found for the application")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}

	return c.JSON(http.StatusOK, status)
}
