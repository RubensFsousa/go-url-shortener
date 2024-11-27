package config

import (
	"fmt"
	"os"

	"github.com/RubensFsousa/go-url-shortener/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPSQL() (*gorm.DB, error) {
	logger := GetLogger("PSQL")

	godotenv.Load()

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPass, dbName, dbPort,
	)

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
