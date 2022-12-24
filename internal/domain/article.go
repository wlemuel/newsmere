package domain

import (
	"gorm.io/gorm"
)

// Article is representing the Article data struct
type Article struct {
	gorm.Model
	GroupId uint   `json:"group_id"`
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
	Author  string `json:"author"`
}

// ArticleRepository represent the article's repository contract
type ArticleRepository interface {
	GetArticleById(id uint) (Article, error)
	StoreArticle(a *Article) error
	DeleteArticle(id uint) error
}

// DBRepository represents the database repository contract
type DBRepository interface {
	ArticleRepository
	GroupRepository
}
