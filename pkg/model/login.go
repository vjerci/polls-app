package model

import "github.com/vjerci/golang-vuejs-sample-app/pkg/util/login"

type LoginRepository interface {
	CreateToken(userID string, groupID string) (login.AccessToken, error)
}
