package api

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zivlakmilos/perfin/db"
)

func (a *Api) GetAccounts(c echo.Context) error {
	store := db.NewAccountStore(db.GetInstance())

	accounts, err := store.GetAll()
	if err != nil {
		return a.ReturnError(c, http.StatusInternalServerError, "retreiving account failed")
	}

	var res []*db.Account
	parentMapping := make(map[string]*db.Account)

	for _, account := range accounts {
		if account.ParentId == "" {
			res = append(res, account)
		} else {
			parentMapping[account.ParentId].Childrens = append(parentMapping[account.ParentId].Childrens, *account)
		}
		parentMapping[account.Id] = account
	}

	return c.JSON(http.StatusOK, map[string]any{
		"accounts": res,
	})
}

func (a *Api) GetAccount(c echo.Context) error {
	store := db.NewAccountStore(db.GetInstance())

	account, err := store.Get(c.Param("id"))
	if err != nil {
		if err == sql.ErrNoRows {
			return a.ReturnError(c, http.StatusNotFound, "account not found")
		}
		return a.ReturnError(c, http.StatusInternalServerError, "retreiving account failed")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"account": account,
	})
}

func (a *Api) CreateAccount(c echo.Context) error {
	var account db.Account
	err := c.Bind(&account)
	if err != nil {
		return a.ReturnError(c, http.StatusInternalServerError, "request parsing failed")
	}

	store := db.NewAccountStore(db.GetInstance())

	err = store.Insert(&account)
	if err != nil {
		return a.ReturnError(c, http.StatusInternalServerError, "creating account failed")
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"account": account,
	})
}

func (a *Api) UpdateAccount(c echo.Context) error {
	var account db.Account
	err := c.Bind(&account)
	if err != nil {
		return a.ReturnError(c, http.StatusInternalServerError, "request parsing failed")
	}

	if account.Id != c.Param("id") {
		return a.ReturnError(c, http.StatusInternalServerError, "account id missmatch")
	}

	store := db.NewAccountStore(db.GetInstance())

	err = store.Update(&account)
	if err != nil {
		if err == sql.ErrNoRows {
			return a.ReturnError(c, http.StatusNotFound, "account not found")
		}
		return a.ReturnError(c, http.StatusInternalServerError, "updateing account failed")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"account": account,
	})
}
