package api

import (
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
