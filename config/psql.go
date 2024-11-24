package config

import (
	"github.com/RubensFsousa/go-url-shortener/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPSQL() (*gorm.DB, error) {
	logger := GetLogger("PSQL")
	dsn := "host=localhost user=postgres password=admin dbname=shortener_db port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Errorf("error to conect on database: %v", err)
		return nil, err
	}

	err = db.AutoMigrate(&models.Url{})
	if err != nil {
		logger.Errorf("migration error: %v", err)
		return nil, err
	}

	return db, nil
}
