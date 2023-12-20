package model

import (
	"strings"

	"gorm.io/gorm"
)

type Learn struct {
	gorm.Model
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Fields    []LearnField   `gorm:"foreignKey:learn_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"fields"`
	UserID    uint           `json:"user_id"`
}
type LearnResponse struct {
	Id     int                    `json:"id"`
	Fields map[string]interface{} `json:"fields"`
}

func ConvertToLearnResponse(db *gorm.DB, learn Learn) LearnResponse {
	db.Preload("Fields").Find(&learn)
	shortenedFields := make([]LearnShortenField, len(learn.Fields))

	for i, field := range learn.Fields {
		shortenedFields[i] = ConvertToLearnShortenField(field)
	}

	result := make(map[string]interface{})

	for _, field := range shortenedFields {
		result[field.Key] = field.Value
	}

	return LearnResponse{
		Id:     int(learn.ID),
		Fields: result,
	}
}

type LearnField struct {
	gorm.Model
	LearnID uint   `gorm:"column:learn_id" json:"learn_id"`
	Key     string `gorm:"column:key" json:"key"`
	Value   string `gorm:"column:value" json:"value"`
}

type LearnShortenField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func ConvertToLearnShortenField(learnfield LearnField) LearnShortenField {
	return LearnShortenField{
		Key:   strings.ToLower(learnfield.Key),
		Value: learnfield.Value,
	}
}
