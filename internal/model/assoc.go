package model

import (
	"gorm.io/gorm"
)

type Assoc struct {
	gorm.Model
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	UserID       uint           `json:"user_id"`
	Organization string         `gorm:"column:organization" json:"organization"`
}

type AssocResponse struct {
	Id           int    `json:"id"`
	Organization string `json:"organization"`
}

func ConvertToAssocResponse(db *gorm.DB, association []Assoc) []AssocResponse {
	var listedAssoc []AssocResponse

	for _, assoc := range association {
		newAssoc := AssocResponse{
			Id:           int(assoc.ID),
			Organization: assoc.Organization,
		}
		listedAssoc = append(listedAssoc, newAssoc)
	}

	return listedAssoc
}
