package db

import (
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
		manager = &Manager{
			db: db,
		}
	}
	return manager
}

func (m *Manager) GetDb() *gorm.DB {
	return m.db
}
