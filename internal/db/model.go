package db

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

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
}

type Group struct {
	gorm.Model
	Name     string
	Articles []Article
	Source   string
}

type Topic struct {
	gorm.Model
	Name   string
	Topics []Topic
	Groups []Group
}

type Subscription struct {
	gorm.Model
	Name        string
	Description string
	High        int
	Low         int
	Source      string
}

type Tag struct {
	gorm.Model
	Name string
}
