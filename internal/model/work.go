package model

import (
	"strings"

	"gorm.io/gorm"
)

type Work struct {
	gorm.Model
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Fields    []WorkField    `gorm:"foreignKey:work_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"fields"`
	UserID    uint           `json:"user_id"`
}
type WorkResponse struct {
	Id     int                    `json:"id"`
	Fields map[string]interface{} `json:"fields"`
}

func ConvertToWorkResponse(db *gorm.DB, work Work) WorkResponse {
	db.Preload("Fields").Find(&work)
	shortenedFields := make([]WorkShortenField, len(work.Fields))

	for i, field := range work.Fields {
		shortenedFields[i] = ConvertToWorkShortenField(field)
	}

	result := make(map[string]interface{})

	for _, field := range shortenedFields {
		result[field.Key] = field.Value
	}

	return WorkResponse{
		Id:     int(work.ID),
		Fields: result,
	}
}

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
