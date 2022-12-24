package domain

import "gorm.io/gorm"

// Topic is representing the Topic data struct
type Topic struct {
	gorm.Model
	TopicId uint   `json:"topic_id"`
	Name    string `json:"name"`
}
