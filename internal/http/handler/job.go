package handler

import (
	"Kaas/internal/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

type Job struct {
	DB *gorm.DB
}

func NewJobHandler(db *gorm.DB) *Job {
	return &Job{DB: db}
}

func (j *Job) GetAppHealth(c echo.Context) error {
	appName := c.Param("appName")
	if appName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Application name is required")
	}

	var status model.MonitorStatus
	result := j.DB.Where("app_name = ?", appName).First(&status)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "No status found for the application")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}

	return c.JSON(http.StatusOK, status)
}
