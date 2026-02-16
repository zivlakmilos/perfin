package api

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
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
	var receipt db.ReceivedReceipt
	err := c.Bind(&receipt)
	if err != nil {
		return a.ReturnError(c, http.StatusInternalServerError, "request parsing failed")
	}

	tx, err := db.StartTransaction(db.GetInstance())
	if err != nil {
		return a.ReturnError(c, http.StatusInternalServerError, "internal server error")
	}
	defer func() { _ = tx.Rollback() }()

	receiptStore := db.NewReceivedReceiptStore(db.GetInstance())
	err = receiptStore.Insert(&receipt)
	if err != nil {
		return a.ReturnError(c, http.StatusInternalServerError, "error saving receipt")
	}

	txId := uuid.NewString()

	transactionStore := db.NewTransactionStore(db.GetInstance())
	transactionStore.UseTransaction(tx)

	debit := make(map[string]float64)

	for _, item := range receipt.Items {
		debit[item.Account] += item.Amount
	}

	totalAmout := float64(0)
	for key, val := range debit {
		transaction := db.Transaction{
			Id:                "",
			TransactionId:     txId,
			AccountId:         key,
			Date:              receipt.Date,
			Description:       "fiscal receipt",
			Debit:             val,
			Credit:            0,
			ReceivedReceiptId: receipt.Id,
		}
		totalAmout += val
		err := transactionStore.Insert(&transaction)
		if err != nil {
			return a.ReturnError(c, http.StatusInternalServerError, "error saving transaction")
		}
	}

	if !utils.CompareFloat4(receipt.TotalAmount, totalAmout) {
		return a.ReturnError(c, http.StatusInternalServerError, "receipt total and item value missmatch")
	}

	transaction := db.Transaction{
		Id:                "",
		TransactionId:     txId,
		AccountId:         "",
		Date:              receipt.Date,
		Description:       "fiscal receipt",
		Debit:             0,
		Credit:            receipt.TotalAmount,
		ReceivedReceiptId: receipt.Id,
	}
	err = transactionStore.Insert(&transaction)
	if err != nil {
		return a.ReturnError(c, http.StatusInternalServerError, "error saving transaction")
	}

	err = tx.Commit()
	if err != nil {
		return a.ReturnError(c, http.StatusInternalServerError, "error saving data")
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"receipt": receipt,
	})
}
