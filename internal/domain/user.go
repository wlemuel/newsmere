package domain

import "gorm.io/gorm"

// User represents a user
type User struct {
	gorm.Model
	Name    string `json:"name" validate:"required"`
	Pass    string `json:"pass"`
	Avatar  string `json:"avatar"`
	Token   string `json:"token"`
	IsAdmin bool   `json:"is_admin"`
}

// UserToken represents a token for user
type UserToken struct {
	gorm.Model
	UserId  uint   `json:"user_id"`
	Token   string `json:"string"`
}

type UserTokenRepository interface {
	CreateToken(u *User) (UserToken, error)
	GetToken(u *User, token string) (UserToken, error)
	DeleteToken(token string) error
}

type UserRepository interface {
	UserTokenRepository
	CreateUser(name, password string) (User, error)
	GetUserByToken(token string) (User, error)
	DeleteUser(id uint) error
	UpdateUser(u *User) error
}
