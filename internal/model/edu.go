package model

import (
	"strings"

	"gorm.io/gorm"
)

type Edu struct {
	gorm.Model
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Fields    []EduField     `gorm:"foreignKey:edu_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"fields"`
	UserID    uint           `json:"user_id"`
}
type EduResponse struct {
	Id     int                    `json:"id"`
	Fields map[string]interface{} `json:"fields"`
}

func ConvertToEduResponse(db *gorm.DB, edu Edu) EduResponse {
	db.Preload("Fields").Find(&edu)
	shortenedFields := make([]EduShortenField, len(edu.Fields))

	for i, field := range edu.Fields {
		shortenedFields[i] = ConvertToEduShortenField(field)
	}

	result := make(map[string]interface{})

	for _, field := range shortenedFields {
		result[field.Key] = field.Value
	}

	return EduResponse{
		Id:     int(edu.ID),
		Fields: result,
	}
}

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
