package app

import (
	"errors"

	echo "github.com/labstack/echo/v4"
	"github.com/vjerci/polls-app/pkg/config"
	"github.com/vjerci/polls-app/pkg/domain/db"
	authmodel "github.com/vjerci/polls-app/pkg/domain/model/auth"
	"github.com/vjerci/polls-app/pkg/domain/model/poll"
	"github.com/vjerci/polls-app/pkg/domain/util/auth"
	"github.com/vjerci/polls-app/pkg/domain/util/googleauth"
	"github.com/vjerci/polls-app/pkg/server/http/api"
	"github.com/vjerci/polls-app/pkg/server/http/api/middleware"
	"github.com/vjerci/polls-app/pkg/server/http/router"
	"github.com/vjerci/polls-app/pkg/server/http/schema"
)

var ErrDBConnect = errors.New("failed to connect to database")

func New(settings config.Config) (*echo.Echo, error) {
	dbClient := db.New(settings.PostgresURL)

	if err := dbClient.Connect(); err != nil {
		return nil, errors.Join(ErrDBConnect, err)
	}

	authClient := auth.New(settings.JWTSigningKey)
	googleAuthClient := googleauth.New(settings.GoogleClientID)

	apiClient := api.New(
		newModel(dbClient, authClient, googleAuthClient),
		newSchemaMap(),
	)

	middlewareClient := middleware.Client{
		AuthRepo: authClient,
		UserRepo: &authmodel.UserModel{
			UserDB: dbClient,
		},
	}

	routeHandler := router.Router{
		Register:    apiClient.Register,
		Login:       apiClient.Login,
		GoogleLogin: apiClient.GoogleLogin,

		PollList:    apiClient.PollList,
		PollCreate:  apiClient.PollCreate,
		PollDetails: apiClient.PollDetails,
		PollVote:    apiClient.PollVote,

		MiddlewareWithAuth: middlewareClient.WithAuth,
		ErrorHandler:       middlewareClient.ErrorHandler,
	}

	return routeHandler.Build(), nil
}

func newModel(dbClient *db.DB,
	authClient authmodel.AuthRepository,
	googleAuthClient authmodel.GoogleAuthRepository,
) *api.Models {
	return &api.Models{
		Login: &authmodel.LoginModel{
			AuthDB: authClient,
			UserDB: dbClient,
		},
		Register: &authmodel.RegisterModel{
			RegisterDB: dbClient,
			AuthDB:     authClient,
		},
		GoogleLogin: &authmodel.GoogleLoginModel{
			AuthDB:       authClient,
			GoogleAuthDB: googleAuthClient,
			UserDB:       dbClient,
			RegisterDB:   dbClient,
		},
		PollList: &poll.ListModel{
			ListDB:  dbClient,
			CountDB: dbClient,
		},
		PollDetails: &poll.DetailsModel{
			DetailsDB: dbClient,
		},
		PollCreate: &poll.CreateModel{
			CreateDB: dbClient,
		},
		PollVote: &poll.VoteModel{
			VoteDB: dbClient,
		},
	}
}

func newSchemaMap() *api.SchemaMap {
	return &api.SchemaMap{
		Login:       &schema.LoginSchemaMap{},
		Register:    &schema.RegisterSchemaMap{},
		GoogleLogin: &schema.GoogleLoginSchemaMap{},
		PollList:    &schema.PollListSchemaMap{},
		PollDetails: &schema.PollDetailsSchemaMap{},
		PollCreate:  &schema.PollCreateSchemaMap{},
		PollVote:    &schema.PollVoteSchemaMap{},
	}
}
