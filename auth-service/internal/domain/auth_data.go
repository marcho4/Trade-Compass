package domain

type AuthData struct {
	userId       int64
	email        string
	passwordHash string
}

func (a *AuthData) GetPasswordHash() string {
	return a.passwordHash
}

func (a *AuthData) GetUserID() int64 {
	return a.userId
}

func (a *AuthData) GetEmail() string {
	return a.email
}
