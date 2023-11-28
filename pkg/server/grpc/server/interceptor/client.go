package interceptor

import "github.com/vjerci/polls-app/pkg/domain/util/auth"

type Client struct {
	AuthRepo AuthRepo
	UserRepo UserRepo
}

type AuthRepo interface {
	Decode(input auth.AccessToken) (userID string, err error)
}

type UserRepo interface {
	GetUser(userID string) (name string, err error)
}
