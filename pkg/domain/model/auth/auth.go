package auth

import "github.com/vjerci/golang-vuejs-sample-app/pkg/domain/util/auth"

type AuthRepository interface {
	CreateToken(userID string) (auth.AccessToken, error)
}
