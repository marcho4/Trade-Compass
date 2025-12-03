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

func NewProviderAuth(id, userID int64, providerUserID string, providerType ProviderType, email string) *ProviderAuth {
	return &ProviderAuth{
		id:             id,
		userId:         userID,
		providerUserId: providerUserID,
		providerType:   providerType,
		email:          email,
	}
}

func (p *ProviderAuth) GetUserID() int64 {
	return p.userId
}

func (p *ProviderAuth) GetProviderUserID() string {
	return p.providerUserId
}

func (p *ProviderAuth) GetProviderType() ProviderType {
	return p.providerType
}

func (p *ProviderAuth) GetEmail() string {
	return p.email
}
