package model

import (
	"strings"

	"gorm.io/gorm"
)

type Vol struct {
	gorm.Model
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Fields    []VolField     `gorm:"foreignKey:vol_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"fields"`
	UserID    uint           `json:"user_id"`
}
type VolResponse struct {
	Id     int                    `json:"id"`
	Fields map[string]interface{} `json:"fields"`
}

func ConvertToVolResponse(db *gorm.DB, vol Vol) VolResponse {
	db.Preload("Fields").Find(&vol)
	shortenedFields := make([]VolShortenField, len(vol.Fields))

	for i, field := range vol.Fields {
		shortenedFields[i] = ConvertToVolShortenField(field)
	}

	result := make(map[string]interface{})

	for _, field := range shortenedFields {
		result[field.Key] = field.Value
	}

	return VolResponse{
		Id:     int(vol.ID),
		Fields: result,
	}
}

type VolField struct {
	gorm.Model
	VolID uint   `gorm:"column:vol_id" json:"vol_id"`
	Key   string `gorm:"column:key" json:"key"`
	Value string `gorm:"column:value" json:"value"`
}

type VolShortenField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func ConvertToVolShortenField(volfield VolField) VolShortenField {
	return VolShortenField{
		Key:   strings.ToLower(volfield.Key),
		Value: volfield.Value,
	}
}
