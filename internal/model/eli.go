package model

import (
	"strings"

	"gorm.io/gorm"
)

type Eli struct {
	gorm.Model
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Fields    []EliField     `gorm:"foreignKey:eli_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"fields"`
	UserID    uint           `json:"user_id"`
}
type EliResponse struct {
	Id     int                    `json:"id"`
	Fields map[string]interface{} `json:"fields"`
}

func ConvertToEliResponse(db *gorm.DB, eli Eli) EliResponse {
	db.Preload("Fields").Find(&eli)
	shortenedFields := make([]EliShortenField, len(eli.Fields))

	for i, field := range eli.Fields {
		shortenedFields[i] = ConvertToEliShortenField(field)
	}

	result := make(map[string]interface{})

	for _, field := range shortenedFields {
		result[field.Key] = field.Value
	}

	return EliResponse{
		Id:     int(eli.ID),
		Fields: result,
	}
}

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
