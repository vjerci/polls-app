package auth

import "github.com/vjerci/polls-app/pkg/domain/util/auth"

type AuthRepository interface {
	CreateToken(userID string) (auth.AccessToken, error)
}
