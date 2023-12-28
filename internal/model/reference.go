package model

import (
	"gorm.io/gorm"
)

type Reference struct {
	gorm.Model
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	UserID    uint           `json:"user_id"`
	Fullname  string         `gorm:"column:fullname" json:"fullname"`
	Address   string         `gorm:"column:address" json:"address"`
	Telephone string         `gorm:"column:telephone" json:"telephone"`
}
type ReferenceResponse struct {
	Id        int    `json:"id"`
	Fullname  string `json:"fullname"`
	Address   string `json:"address"`
	Telephone string `json:"telephone"`
}

func ConvertToReferenceResponse(db *gorm.DB, references []Reference) []ReferenceResponse {
	var listedChild []ReferenceResponse

	for _, child := range references {
		newChild := ReferenceResponse{
			Id:        int(child.ID),
			Fullname:  child.Fullname,
			Address:   child.Address,
			Telephone: child.Telephone,
		}
		listedChild = append(listedChild, newChild)
	}

	return listedChild
}
