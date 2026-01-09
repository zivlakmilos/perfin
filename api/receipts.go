package api

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/zivlakmilos/perfin/utils"
)

func (a *Api) ProcessFiscalReceipt(c echo.Context) error {
	var req struct {
		ReceiptUrl string `json:"receiptUrl"`
	}
	if err := c.Bind(&req); err != nil {
		return a.ReturnError(c, 500, "request parsing failed")
	}

	receiptData, err := utils.GetFiscalReceiptInfo(req.ReceiptUrl)
	if err != nil {
		return a.ReturnError(c, 500, "processing fiscal receipt failed")
	}

	fmt.Printf("%+v\n", receiptData)
	items, err := utils.ParseFiscalReceiptItems(receiptData)
	if err != nil {
		return a.ReturnError(c, 500, "processing fiscal receipt items failed")
	}

	fmt.Printf("%+v\n", receiptData)
	for _, item := range items {
		fmt.Printf("%+v\n", item)
	}

	return c.JSON(200, map[string]any{
		"receipt": receiptData,
		"items":   items,
	})
}
