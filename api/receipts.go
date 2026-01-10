package api

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zivlakmilos/perfin/db"
	"github.com/zivlakmilos/perfin/utils"
)

func (a *Api) ProcessFiscalReceipt(c echo.Context) error {
	var req struct {
		ReceiptUrl string `json:"receiptUrl"`
	}
	if err := c.Bind(&req); err != nil {
		return a.ReturnError(c, http.StatusInternalServerError, "request parsing failed")
	}

	receiptData, err := utils.GetFiscalReceiptInfo(req.ReceiptUrl)
	if err != nil {
		return a.ReturnError(c, http.StatusInternalServerError, "processing fiscal receipt failed")
	}

	items, err := utils.ParseFiscalReceiptItems(receiptData)
	if err != nil {
		return a.ReturnError(c, http.StatusInternalServerError, "processing fiscal receipt items failed")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"receipt": receiptData,
		"items":   items,
	})
}

func (a *Api) GetFiscalReceipts(c echo.Context) error {
	store := db.NewReceivedReceiptStore(db.GetInstance())

	receipts, err := store.GetAll()
	if err != nil {
		return a.ReturnError(c, http.StatusInternalServerError, "error retreiving receipts")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"receipts": receipts,
	})
}

func (a *Api) GetFiscalReceipt(c echo.Context) error {
	store := db.NewReceivedReceiptStore(db.GetInstance())

	receipt, err := store.Get(c.Param("id"))
	if err != nil {
		if err == sql.ErrNoRows {
			return a.ReturnError(c, http.StatusNotFound, "receipt not found")
		}
		return a.ReturnError(c, http.StatusInternalServerError, "error retreiving receipt")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"receipt": receipt,
	})
}

func (a *Api) CreateFiscalReceipt(c echo.Context) error {
	var req db.ReceivedReceipt
	err := c.Bind(&req)
	if err != nil {
		return a.ReturnError(c, http.StatusInternalServerError, "request parsing failed")
	}

	store := db.NewReceivedReceiptStore(db.GetInstance())
	err = store.Insert(&req)
	if err != nil {
		return a.ReturnError(c, http.StatusInternalServerError, "error saving receipt")
	}

	// TODO: create transaction based on receipt

	return c.JSON(http.StatusCreated, map[string]any{
		"receipt": req,
	})
}
