package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string         `gorm:"unique"`
	Email     string         `gorm:"unique"`
	Password  string         `json:"-"`
	Active    string         `gorm:"column:active" json:"active"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Fields    []UserField    `gorm:"foreignKey:user_id"`
}
type UserResponse struct {
	Username string            `json:"username"`
	Email    string            `json:"email"`
	Fields   map[string]string `json:"fields"`
}

func ConvertToUserResponse(db *gorm.DB, user User) UserResponse {
	db.Preload("Fields").Find(&user)
	shortenedFields := make([]ShortenField, len(user.Fields))

	for i, field := range user.Fields {
		shortenedFields[i] = ConvertToShortenField(field)
	}

	result := make(map[string]string)

	for _, field := range shortenedFields {
		result[field.Key] = field.Value
	}

	return UserResponse{
		Username: user.Username,
		Email:    user.Email,
		Fields:   result,
	}
}
