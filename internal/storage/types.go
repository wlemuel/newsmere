package storage

import (
	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var manager *Manager

type Manager struct {
	db *gorm.DB
}

func init() {
	_ = GetManager()
}

func GetManager() *Manager {
	if manager == nil {
		db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		err = db.AutoMigrate(
			&User{},
			&Article{},
			&Group{},
			&Topic{},
			&Subscription{},
		)
		if err != nil {
			panic("failed to automigrate tables")
		}

		manager = &Manager{
			db: db,
		}
	}
	return manager
}

func GetDb() *gorm.DB {
	return GetManager().db
}

type User struct {
	gorm.Model
	Name    string
	Pass    string
	IsAdmin bool
	Active  bool
}

type Article struct {
	gorm.Model
	MsgID   string
	DocType string
	Bytes   int
	Lines   int
	Title   string
	Headers datatypes.JSON
	GroupId uint
	Tags    []Tag
}

type Group struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex:idx_group_name_source"`
	Description string
	Articles    []Article
	Source      string `gorm:"uniqueIndex:idx_group_name_source"`
	TopicId     uint
	Enabled     bool `gorm:"default:true"`
	Low         int
	High        int
}

type Topic struct {
	gorm.Model
	Name    string
	Topics  []Topic
	Groups  []Group
	TopicId uint
}

type Subscription struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex:idx_sub_name_source"`
	Description string
	High        int
	Low         int
	Source      string `gorm:"uniqueIndex:idx_sub_name_source"`
}

type Tag struct {
	gorm.Model
	Name      string
	ArticleId uint
}
