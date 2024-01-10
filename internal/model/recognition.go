package model

import (
	"gorm.io/gorm"
)

type Recognition struct {
	gorm.Model
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	UserID      uint           `json:"user_id"`
	Recognition string         `gorm:"column:recognition" json:"recognition"`
}
type RecognitionResponse struct {
	Id          int    `json:"id"`
	Recognition string `json:"recognition"`
}

func ConvertToRecognitionResponse(db *gorm.DB, children []Recognition) []RecognitionResponse {
	var listedChild []RecognitionResponse

	for _, child := range children {
		newChild := RecognitionResponse{
			Id:          int(child.ID),
			Recognition: child.Recognition,
		}
		listedChild = append(listedChild, newChild)
	}

	return listedChild
}
