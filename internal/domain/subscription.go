package domain

import "gorm.io/gorm"

// Subscription is representing the sub data struct
type Subscription struct {
	gorm.Model
	Name string `json:"name"`
	High int    `json:"high"`
	Low  int    `json:"low"`
	Type string `json:"type"`
	Repo string `json:"repo"`
}

type SubscriptionRepository interface {
	CountSubsByTypeRepo(t, repo string) (int64, error)
	StoreSubs(subs []Subscription) error
}
