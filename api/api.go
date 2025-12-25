package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Api struct {
	e *echo.Echo
}

func NewApi(e *echo.Echo) *Api {
	return &Api{
		e: e,
	}
}

func (a *Api) SetupRoutes() {
	e := a.e
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!")
	})
}
