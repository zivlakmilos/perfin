package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/zivlakmilos/perfin/db"
)

func (a *Api) Login(c echo.Context) error {
	user := &db.User{}

	err := c.Bind(user)
	if err != nil {
		return err
	}

	store := db.NewUserStore(db.GetInstance())
	user, err = store.Login(user.Username, user.Password)
	if err != nil {
		c.JSON(http.StatusForbidden, map[string]string{
			"error": "wrong username or password",
		})
		return err
	}
	user.Password = ""

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user,
	})

	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		fmt.Print(err)
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": tokenStr,
	})
}

func (a *Api) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")
		if len(auth) < len("Bearer ") {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "invalid token",
			})
		}

		tokenStr := auth[len("Bearer "):]
		_, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "invalid token",
			})
		}

		return next(c)
	}
}
