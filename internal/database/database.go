package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

func InitializeDB() (*gorm.DB, error) {

	time.Sleep(5 * time.Second)
	dsn := "host=kaas-api-postgres port=5432 user=postgres dbname=monitoring password=postgres sslmode=disable TimeZone=Asia/Tehran"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
