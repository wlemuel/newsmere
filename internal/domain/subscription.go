package domain

import "gorm.io/gorm"

// Subscription is representing the sub data struct
type Subscription struct {
	gorm.Model
	Name string `json:"name"`
	Repo string `json:"repo"`
}
