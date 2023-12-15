package model

import (
	"gorm.io/gorm"
)

type Children struct {
	gorm.Model
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	UserID      uint           `json:"user_id"`
	Fullname    string         `gorm:"column:fullname" json:"fullname"`
	Dateofbirth string         `gorm:"column:dateofbirth" json:"dateofbirth"`
}
type ChildrenResponse struct {
	Id          int    `json:"id"`
	Fullname    string `json:"fullname"`
	Dateofbirth string `json:"dateofbirth"`
}

func ConvertToChildrenResponse(db *gorm.DB, children []Children) []ChildrenResponse {
	var listedChild []ChildrenResponse

	for _, child := range children {
		newChild := ChildrenResponse{
			Id:          int(child.ID),
			Fullname:    child.Fullname,
			Dateofbirth: child.Dateofbirth,
		}
		listedChild = append(listedChild, newChild)
	}

	return listedChild
}
