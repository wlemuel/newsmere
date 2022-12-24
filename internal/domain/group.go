package domain

import "gorm.io/gorm"

// Group is representing the Group data struct
type Group struct {
	gorm.Model
	TopicId uint   `json:"topic_id"`
	Name    string `json:"name"`
	Repo    string `json:"repo"`
	SubId   uint   `json:"sub_id"`
}

// GroupRepository represent the group's repository contract
type GroupRepository interface {
	GetGroupById(id uint) (Group, error)
	StoreGroup(g *Group) error
	DeleteGroup(id uint) error
}
