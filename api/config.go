package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zivlakmilos/perfin/db"
)

func (a *Api) GetConfig(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"config": a.config,
	})
}

func (a *Api) ReplaceConfig(c echo.Context) error {
	req := map[string]string{}
	if err := c.Bind(&req); err != nil {
		return a.ReturnError(c, http.StatusInternalServerError, "request parsing failed")
	}

	store := db.NewConfigStore(db.GetInstance())
	err := store.ReplaceConfig(req)
	if err != nil {
		return a.ReturnError(c, http.StatusInternalServerError, "config replace failed")
	}

	a.config = req

	return c.JSON(http.StatusOK, map[string]any{
		"config": a.config,
	})
}

func (a *Api) UpdateConfig(c echo.Context) error {
	req := struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}{}
	if err := c.Bind(&req); err != nil {
		return a.ReturnError(c, http.StatusInternalServerError, "request parsing failed")
	}

	store := db.NewConfigStore(db.GetInstance())
	err := store.SetValue(req.Key, req.Value)
	if err != nil {
		return a.ReturnError(c, http.StatusInternalServerError, "config update failed")
	}

	a.config[req.Key] = req.Value

	return c.JSON(http.StatusOK, map[string]any{
		"config": a.config,
	})
}
