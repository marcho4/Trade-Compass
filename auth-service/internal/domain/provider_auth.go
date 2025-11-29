package domain

type ProviderType string

const (
	// ProviderTypeGoogle ProviderType = "google"
	ProviderTypeYandex ProviderType = "yandex"
)

type ProviderAuth struct {
	id             int64
	userId         int64
	providerUserId string
	providerType   ProviderType
	email          string
}
