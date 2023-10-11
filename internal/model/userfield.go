package model

import (
	"strings"

	"gorm.io/gorm"
)

type UserField struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	Key    string `gorm:"column:key" json:"key"`
	Value  string `gorm:"column:value" json:"value"`
}

type ShortenField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func ConvertToShortenField(userfield UserField) ShortenField {
	return ShortenField{
		Key:   strings.ToLower(userfield.Key),
		Value: userfield.Value,
	}
}
