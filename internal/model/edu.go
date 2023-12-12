package model

import (
	"gorm.io/gorm"
)

type Edu struct {
	gorm.Model
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Fields    []EduField     `gorm:"foreignKey:edu_id"`
}
type EduResponse struct {
	Id     int               `json:"id"`
	Fields map[string]string `json:"fields"`
}

func ConvertToEduResponse(db *gorm.DB, edu Edu) EduResponse {
	db.Preload("Fields").Find(&edu)
	shortenedFields := make([]EduShortenField, len(edu.Fields))

	for i, field := range edu.Fields {
		shortenedFields[i] = ConvertToEduShortenField(field)
	}

	result := make(map[string]string)

	for _, field := range shortenedFields {
		result[field.Key] = field.Value
	}

	return EduResponse{
		Id:     int(edu.ID),
		Fields: result,
	}
}
