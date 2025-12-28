package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/zivlakmilos/perfin/api"
	"github.com/zivlakmilos/perfin/db"
)

func main() {
	e := echo.New()

	err := godotenv.Load()
	if err != nil {
		e.Logger.Fatal()
	}
	db.CreateConnection(os.Getenv("DB_URL"))

	a := api.NewApi(e)
	a.SetupRoutes()

	e.Logger.Fatal(e.Start(":9999"))
}
