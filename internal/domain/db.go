package domain

// DBRepository represents the database repository contract
type DBRepository interface {
	ArticleRepository
	GroupRepository
	SubscriptionRepository
}
