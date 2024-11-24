package models

import "gorm.io/gorm"

type Url struct {
	gorm.Model
	ID         uint
	DecodedUrl string `gorm:"not null"`
	CodedUrl   string `gorm:"unique;not null"`
}
