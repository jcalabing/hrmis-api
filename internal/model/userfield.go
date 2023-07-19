package model

import "gorm.io/gorm"

type UserField struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

type ShortenField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func ConvertToShortenField(userfield UserField) ShortenField {
	return ShortenField{
		Key:   userfield.Key,
		Value: userfield.Value,
	}
}
