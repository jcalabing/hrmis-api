package model

import (
	"gorm.io/gorm"
)

type Skill struct {
	gorm.Model
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	UserID     uint           `json:"user_id"`
	Skillhobby string         `gorm:"column:skillhobby" json:"skillhobby"`
}
type SkillResponse struct {
	Id         int    `json:"id"`
	Skillhobby string `json:"skillhobby"`
}

func ConvertToSkillResponse(db *gorm.DB, children []Skill) []SkillResponse {
	var listedChild []SkillResponse

	for _, child := range children {
		newChild := SkillResponse{
			Id:         int(child.ID),
			Skillhobby: child.Skillhobby,
		}
		listedChild = append(listedChild, newChild)
	}

	return listedChild
}
