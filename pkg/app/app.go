package app

import (
	"errors"

	echo "github.com/labstack/echo/v4"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/config"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/db"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/model"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/domain/util/auth"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/http/api"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/http/route"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/server/schema"
)

var ErrDBConnect = errors.New("failed to connect to database")

func New(settings config.Config) (*echo.Echo, error) {
	dbClient := db.New(settings.PostgresURL)

	if err := dbClient.Connect(); err != nil {
		return nil, errors.Join(ErrDBConnect, err)
	}

	authClient := auth.New(settings.JWTSigningKey)
	modelClient := &model.Client{
		RegisterDB:    dbClient,
		AuthDB:        authClient,
		UserDB:        dbClient,
		PollListDB:    dbClient,
		PollCreateDB:  dbClient,
		PollDetailsDB: dbClient,
		PollVoteDB:    dbClient,
	}
	apiClient := api.New()

	schemaMap := schema.NewSchemaMap()

	routeHandler := route.Handler{
		Register:    apiClient.Register(modelClient, schemaMap),
		Login:       apiClient.Login(modelClient, schemaMap),
		PollList:    apiClient.PollList(modelClient, schemaMap),
		PollCreate:  apiClient.PollCreate(modelClient, schemaMap),
		PollDetails: apiClient.PollDetails(modelClient, schemaMap),
		PollVote:    apiClient.PollVote(modelClient, schemaMap),
	}

	return routeHandler.Build(), nil
}
