package sqlite

import (
	"newsmere/internal/domain"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type sqliteRepo struct {
	DB *gorm.DB
}

func NewSqliteRepo() domain.DBRepository {
	dsn := "test.db"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database 'test.db'")
	}

	return &sqliteRepo{
		DB: db,
	}
}

func (r *sqliteRepo) GetArticleById(id uint) (res domain.Article, err error) {
	r.DB.First(&res, id)
	err = nil

	return
}

func (r *sqliteRepo) StoreArticle(a *domain.Article) error {
	return nil
}

func (r *sqliteRepo) DeleteArticle(id uint) error {
	return nil
}

func (r *sqliteRepo) GetGroupById(id uint) (res domain.Group, err error) {
	r.DB.First(&res, id)
	err = nil

	return
}

func (r *sqliteRepo) StoreGroup(g *domain.Group) error {
	return nil
}

func (r *sqliteRepo) DeleteGroup(id uint) error {
	return nil
}
