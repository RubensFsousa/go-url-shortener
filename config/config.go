package config

import (
	"fmt"

	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	logger *Logger
)

func Init() error {
	var err error

	db, err = InitPSQL()
	if err != nil {
		return fmt.Errorf("error to initializing: %v", err)
	}

	return err
}

func GetPSQL() *gorm.DB {
	return db
}

func GetLogger(p string) *Logger {
	logger = newLogger(p)
	return logger
}
