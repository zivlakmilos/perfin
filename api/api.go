package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Api struct {
	echo *echo.Echo
}

func NewApi(e *echo.Echo) *Api {
	return &Api{
		echo: e,
	}
}

func (a *Api) ReturnError(c echo.Context, code int, msg string) error {
	return c.JSON(code, map[string]any{
		"error": msg,
	})
}

func (a *Api) SetupRoutes() {
	e := a.echo

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	auth := e.Group("/auth")
	auth.POST("/login", a.Login)

	accounts := e.Group("/accounts", a.AuthMiddleware)
	accounts.GET("", a.GetAccounts)

	mappings := e.Group("/mappings", a.AuthMiddleware)
	mappings.GET("/items", a.GetItemMappings)
	mappings.POST("/items", a.CreateItemMapping)

	fiscalReceipts := e.Group("/fiscal_receipts", a.AuthMiddleware)
	fiscalReceipts.POST("/process", a.ProcessFiscalReceipt)
}
