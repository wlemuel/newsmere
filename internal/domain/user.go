package domain

import "gorm.io/gorm"

// User represents a user
type User struct {
	gorm.Model
	Name    string `json:"name" validate:"required"`
	Pass    string `json:"pass"`
	Avatar  string `json:"avatar"`
	Token   string `json:"token"`
	IsValid bool   `json:"is_valid"`
	IsAdmin bool   `json:"is_admin"`
}

// UserToken represents a token for user
type UserToken struct {
	gorm.Model
	UserId  uint   `json:"user_id"`
	Token   string `json:"string"`
	IsValid bool   `json:"is_valid"`
}
