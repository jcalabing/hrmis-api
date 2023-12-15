package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username   string         `gorm:"unique"`
	Email      string         `gorm:"unique"`
	Password   string         `json:"-"`
	Active     string         `gorm:"column:active" json:"active"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Fields     []UserField    `gorm:"foreignKey:user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"fields"`
	Educations []Edu          `gorm:"foreignKey:user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"educations"`
	Children   []Children     `gorm:"foreignKey:user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"children"`
}
type UserResponse struct {
	Id         int                `json:"id"`
	Username   string             `json:"username"`
	Email      string             `json:"email"`
	Fields     map[string]string  `json:"fields"`
	Educations []interface{}      `json:"educations"`
	Children   []ChildrenResponse `json:"children"`
}

func ConvertToUserResponse(db *gorm.DB, user User) UserResponse {
	db.Preload("Fields").Preload("Educations").Preload("Children").Find(&user)
	shortenedFields := make([]UserShortenField, len(user.Fields))

	//user profile fields
	for i, field := range user.Fields {
		shortenedFields[i] = ConvertToUserShortenField(field)
	}

	profileFieldResult := make(map[string]string)

	for _, field := range shortenedFields {
		profileFieldResult[field.Key] = field.Value
	}

	//user education fields
	var eduResponse []interface{}

	for _, eduData := range user.Educations {
		convertedData := ConvertToEduResponse(db, eduData).Fields
		convertedData["id"] = eduData.ID
		eduResponse = append(eduResponse, convertedData)
	}

	//user children fields
	childResponse := ConvertToChildrenResponse(db, user.Children)

	return UserResponse{
		Id:         int(user.ID),
		Username:   user.Username,
		Email:      user.Email,
		Fields:     profileFieldResult,
		Educations: eduResponse,
		Children:   childResponse,
	}
}
