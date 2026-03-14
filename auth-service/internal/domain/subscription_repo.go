package domain

type SubscriptionRepo interface {
	GetUserSubscription(userId int) error
	CreateFreeSubscription(userId int) error
	ChangeSubscriptionType(userId, newSubscriptionType int) error
}
