package models

type Url struct {
	ID         uint
	DecodedUrl string `gorm:"not null"`
	CodedUrl   string `gorm:"unique;not null"`
}
