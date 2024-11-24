package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	logger *Logger
)

func Init() error {
	dsn := "host=localhost user=postgres password=admin dbname=shortener_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return err
}

func GetLogger(p string) *Logger {
	logger := newLogger(p)
	return logger
}
