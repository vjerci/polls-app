package route

import (
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vjerci/golang-vuejs-sample-app/pkg/config"
)

type Handler struct {
	Register echo.HandlerFunc
	Login    echo.HandlerFunc

	PollList    echo.HandlerFunc
	PollDetails echo.HandlerFunc

	PollCreate echo.HandlerFunc
	PollVote   echo.HandlerFunc
}

func (handler *Handler) Build() *echo.Echo {
	router := echo.New()

	router.Use(middleware.Logger())
	router.Use(middleware.Gzip())

	if config.Get().Debug {
		//nolint:exhaustivestruct
		router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAcceptEncoding},
		}))
	}

	router.Use(middleware.Static("./front/dist"))

	endpoint := router.Group("/api")

	auth := endpoint.Group("/auth")
	auth.PUT("/register", handler.Register)
	auth.POST("/login", handler.Login)

	endpoint.GET("/poll", handler.PollList)
	endpoint.GET("/poll/:id", handler.PollDetails)

	endpoint.PUT("/poll", handler.PollCreate)
	endpoint.POST("/poll/:id/vote", handler.PollVote)

	return router
}
