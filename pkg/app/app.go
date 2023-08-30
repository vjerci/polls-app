package app

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/api"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/config"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/db"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/model"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/route"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/util/login"
)

var ErrDbConnect = errors.New("failed to connect to database")

func New(settings config.Config) (*echo.Echo, error) {
	dbClient := db.New(settings.PostgresURL)
	err := dbClient.Connect()
	if err != nil {
		return nil, errors.Join(ErrDbConnect, err)
	}

	loginClient := login.New(settings.JWTSigningKey)
	modelClient := &model.Client{
		RegisterDB: dbClient,
		LoginDB:    loginClient,
		UserDB:     dbClient,
	}
	apiClient := api.New()

	routeHandler := route.Handler{
		Register: apiClient.Register(modelClient),
		Login:    apiClient.Login(modelClient),
	}

	return routeHandler.Build(), nil
}
