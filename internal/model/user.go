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
}
type UserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func ConvertToUserResponse(user User) UserResponse {
	return UserResponse{
		Username: user.Username,
		Email:    user.Email,
	}
}
