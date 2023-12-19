package model

import (
	"strings"

	"gorm.io/gorm"
)

type EliField struct {
	gorm.Model
	EliID uint   `gorm:"column:eli_id" json:"eli_id"`
	Key   string `gorm:"column:key" json:"key"`
	Value string `gorm:"column:value" json:"value"`
}

type EliShortenField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func ConvertToEliShortenField(elifield EliField) EliShortenField {
	return EliShortenField{
		Key:   strings.ToLower(elifield.Key),
		Value: elifield.Value,
	}
}
