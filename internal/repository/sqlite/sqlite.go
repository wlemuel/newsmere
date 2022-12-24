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

func (r *sqliteRepo) GetArticles(groupId uint, prevId uint) (res []domain.Article, err error) {
	result := r.DB.Limit(10).Find(&res, "group_id = ? and id > ?", groupId, prevId)
	err = result.Error

	return
}

func (r *sqliteRepo) GetArticleById(id uint) (res domain.Article, err error) {
	result := r.DB.First(&res, id)
	err = result.Error

	return
}

func (r *sqliteRepo) StoreArticle(a *domain.Article) error {
	result := r.DB.Create(a)
	return result.Error
}

func (r *sqliteRepo) DeleteArticle(id uint) error {
	result := r.DB.Delete(&domain.Article{}, id)
	return result.Error
}

func (r *sqliteRepo) GetGroupById(id uint) (res domain.Group, err error) {
	result := r.DB.First(&res, id)
	err = result.Error

	return
}

func (r *sqliteRepo) StoreGroup(g *domain.Group) error {
	result := r.DB.Create(g)
	return result.Error
}

func (r *sqliteRepo) DeleteGroup(id uint) error {
	result := r.DB.Delete(&domain.Group{}, id)
	return result.Error
}

func (r *sqliteRepo) CreateToken(u *domain.User) (res domain.UserToken, err error) {
	result := r.DB.First(&res)
	err = result.Error

	return
}

func (r *sqliteRepo) GetToken(u *domain.User, token string) (res domain.UserToken, err error) {
	result := r.DB.First(&res, "user_id = ?", u.ID)
	err = result.Error

	return
}
