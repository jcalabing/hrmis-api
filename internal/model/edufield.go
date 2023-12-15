package model

import (
	"strings"

	"gorm.io/gorm"
)

type EduField struct {
	gorm.Model
	EduID uint   `gorm:"column:edu_id" json:"edu_id"`
	Key   string `gorm:"column:key" json:"key"`
	Value string `gorm:"column:value" json:"value"`
}

type EduShortenField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func ConvertToEduShortenField(edufield EduField) EduShortenField {
	return EduShortenField{
		Key:   strings.ToLower(edufield.Key),
		Value: edufield.Value,
	}
}
