package api

import (
	"github.com/labstack/echo/v4"
)

type Api struct {
	echo *echo.Echo
}

func NewApi(e *echo.Echo) *Api {
	return &Api{
		echo: e,
	}
}

func (a *Api) SetupRoutes() {
	e := a.echo

	auth := e.Group("/auth")
	auth.POST("/login", a.Login)

	accounts := e.Group("/accounts", a.AuthMiddleware)
	accounts.GET("", a.GetAccounts)
}
