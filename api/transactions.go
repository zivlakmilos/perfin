package api

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zivlakmilos/perfin/db"
)

func (a *Api) GetTransactions(c echo.Context) error {
	store := db.NewTransactionStore(db.GetInstance())

	transactions, err := store.GetAll()
	if err != nil {
		return a.ReturnError(c, http.StatusInternalServerError, "retreiving transactions failed")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"transactions": transactions,
	})
}

func (a *Api) GetTransaction(c echo.Context) error {
	store := db.NewTransactionStore(db.GetInstance())

	transaction, err := store.Get(c.Param("id"))
	if err != nil {
		if err == sql.ErrNoRows {
			return a.ReturnError(c, http.StatusNotFound, "transaction not found")
		}
		return a.ReturnError(c, http.StatusInternalServerError, "retreiving transaction failed")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"transaction": transaction,
	})
}
