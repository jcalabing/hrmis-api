package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string         `gorm:"unique"`
	Email     string         `gorm:"unique"`
	Password  string         `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Fields    []UserField    `gorm:"foreignKey:user_id"`
}
type UserResponse struct {
	Username string         `json:"username"`
	Email    string         `json:"email"`
	Fields   []ShortenField `json:"fields"`
}

func ConvertToUserResponse(db *gorm.DB, user User) UserResponse {
	db.Preload("Fields").Find(&user)
	shortenedFields := make([]ShortenField, len(user.Fields))

	for i, field := range user.Fields {
		shortenedFields[i] = ConvertToShortenField(field)
	}

	return UserResponse{
		Username: user.Username,
		Email:    user.Email,
		Fields:   shortenedFields,
	}
}
