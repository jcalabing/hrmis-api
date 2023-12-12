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

type UserShortenField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func ConvertToUserShortenField(userfield UserField) UserShortenField {
	return UserShortenField{
		Key:   strings.ToLower(userfield.Key),
		Value: userfield.Value,
	}
}
