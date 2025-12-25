package main

import (
	"github.com/labstack/echo/v4"
	"github.com/zivlakmilos/perfin/api"
)

func main() {
	e := echo.New()

	a := api.NewApi(e)
	a.SetupRoutes()

	e.Logger.Fatal(e.Start(":9999"))
}
