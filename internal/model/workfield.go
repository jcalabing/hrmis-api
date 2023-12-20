package model

import (
	"strings"

	"gorm.io/gorm"
)

type WorkField struct {
	gorm.Model
	WorkID uint   `gorm:"column:work_id" json:"work_id"`
	Key    string `gorm:"column:key" json:"key"`
	Value  string `gorm:"column:value" json:"value"`
}

type WorkShortenField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func ConvertToWorkShortenField(workfield WorkField) WorkShortenField {
	return WorkShortenField{
		Key:   strings.ToLower(workfield.Key),
		Value: workfield.Value,
	}
}
