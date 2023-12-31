package router

import (
	"log"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vjerci/polls-app/pkg/config"
)

type Router struct {
	Register    echo.HandlerFunc
	Login       echo.HandlerFunc
	GoogleLogin echo.HandlerFunc

	PollList    echo.HandlerFunc
	PollDetails echo.HandlerFunc

	PollCreate echo.HandlerFunc
	PollVote   echo.HandlerFunc

	MiddlewareWithAuth echo.MiddlewareFunc
	ErrorHandler       echo.HTTPErrorHandler
}

func (handler *Router) Build() *echo.Echo {
	router := echo.New()

	router.HTTPErrorHandler = handler.ErrorHandler

	router.Use(middleware.Logger())
	router.Use(middleware.Gzip())

	if config.Get().Debug {
		log.Println("running server with cors enabled")

		//nolint:exhaustivestruct
		router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*", "http://localhost:3000/"},
			AllowHeaders: []string{
				echo.HeaderOrigin,
				echo.HeaderContentType,
				echo.HeaderAcceptEncoding,
				echo.HeaderAuthorization,
			},
		}))
	}

	router.GET("/healthcheck", func(echoContext echo.Context) error {
		return echoContext.String(http.StatusOK, "OK")
	})
	router.Use(middleware.Static("./front/dist"))

	endpoint := router.Group("/api")

	auth := endpoint.Group("/auth")
	auth.PUT("/register", handler.Register)
	auth.POST("/login", handler.Login)
	auth.POST("/google/login", handler.GoogleLogin)

	poll := endpoint.Group("/poll", handler.MiddlewareWithAuth)
	poll.GET("", handler.PollList)
	poll.GET("/:id", handler.PollDetails)
	poll.PUT("", handler.PollCreate)
	poll.POST("/:id/vote", handler.PollVote)

	return router
}
