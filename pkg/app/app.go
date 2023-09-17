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

	apiClient := api.New(
		newModel(dbClient, authClient),
		newSchemaMap(),
	)

	routeHandler := route.Handler{
		Register:    apiClient.Register,
		Login:       apiClient.Login,
		PollList:    apiClient.PollList,
		PollCreate:  apiClient.PollCreate,
		PollDetails: apiClient.PollDetails,
		PollVote:    apiClient.PollVote,
	}

	return routeHandler.Build(), nil
}

func newModel(dbClient *db.Client, authClient model.AuthRepository) *api.Models {
	return &api.Models{
		Login: &model.LoginModel{
			AuthDB: authClient,
			UserDB: dbClient,
		},
		Register: &model.RegisterModel{
			RegisterDB: dbClient,
			AuthDB:     authClient,
		},
		PollList: &model.PollListModel{
			PollListDB: dbClient,
		},
		PollDetails: &model.PollDetailsModel{
			PollDetailsDB: dbClient,
		},
		PollCreate: &model.PollCreateModel{
			PollCreateDB: dbClient,
		},
		PollVote: &model.PollVoteModel{
			PollVoteDB: dbClient,
		},
	}
}

func newSchemaMap() *api.SchemaMap {
	return &api.SchemaMap{
		Login:       &schema.LoginSchemaMap{},
		Register:    &schema.RegisterSchemaMap{},
		PollList:    &schema.PollListSchemaMap{},
		PollDetails: &schema.PollDetailsSchemaMap{},
		PollCreate:  &schema.PollCreateSchemaMap{},
		PollVote:    &schema.PollVoteSchemaMap{},
	}
}
