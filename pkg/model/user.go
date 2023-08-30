package model

type UserRepository interface {
	GetUser(userID string) (name string, err error)
}
