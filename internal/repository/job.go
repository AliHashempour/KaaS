package repository

import (
	"Kaas/internal/model"
	"gorm.io/gorm"
)

type JobRepository interface {
	GetAppHealth(appName string) (*model.MonitorStatus, error)
}

type jobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) JobRepository {
	return &jobRepository{db: db}
}

func (r *jobRepository) GetAppHealth(appName string) (*model.MonitorStatus, error) {
	var status model.MonitorStatus
	result := r.db.Where("app_name = ?", appName).First(&status)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return &status, nil
}
