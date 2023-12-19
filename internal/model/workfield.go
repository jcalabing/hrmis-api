package model

import (
	"strings"

	"gorm.io/gorm"
)

type WorkField struct {
	gorm.Model
	WorkID uint   `gorm:"column:edu_id" json:"edu_id"`
	Key    string `gorm:"column:key" json:"key"`
	Value  string `gorm:"column:value" json:"value"`
}

type WorkShortenField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func ConvertToWorkShortenField(edufield WorkField) WorkShortenField {
	return WorkShortenField{
		Key:   strings.ToLower(edufield.Key),
		Value: edufield.Value,
	}
}
