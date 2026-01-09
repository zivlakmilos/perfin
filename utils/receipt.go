package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type FiscalReceiptInfo struct {
	InvoiceRequest struct {
		PosTime                *string `json:"posTime"`
		TaxId                  string  `json:"taxId"`
		BusinessName           string  `json:"businessName"`
		LocationName           string  `json:"locationName"`
		Address                string  `json:"address"`
		City                   string  `json:"city"`
		AdministrativeUnit     string  `json:"administrativeUnit"`
		Buyer                  *string `json:"buyer"`
		BuyerCostCenter        *string `json:"buyerCostCenter"`
		Cashier                *string `json:"cashier"`
		RequestedBy            string  `json:"requestedBy"`
		ReferentDocumentNumber *string `json:"referentDocumentNumber"`
		InvoiceType            int     `json:"invoiceType"`
		TransactionType        int     `json:"transactionType"`
		Payments               []struct {
			PaymentType            int     `json:"paymentType"`
			PaymentTypeDescription string  `json:"paymentTypeDescript"`
			Amount                 float64 `json:"amount"`
		} `json:"payments"`
	} `json:"invoiceRequest"`
	InvoiceResult struct {
		TotalAmount             float64 `json:"totalAmount"`
		TransactionTypeCounter  int     `json:"transactionTypeCounter"`
		TotalCounter            int     `json:"totalCounter"`
		InvoiceCounterExtension string  `json:"invoiceCounterExtension"`
		InvoiceNumber           string  `json:"invoiceNumber"`
		SignedBy                string  `json:"signedBy"`
		SdcTime                 string  `json:"sdcTime"`
	} `json:"invoiceResult"`
	Journal string `json:"journal"`
	IsValid bool   `json:"isValid"`
}

type FiscalReceiptItem struct {
	Title    string
	Price    float64
	Quantity float64
	Amount   float64
}

func GetFiscalReceiptInfo(url string) (*FiscalReceiptInfo, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var result FiscalReceiptInfo
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	for i := range result.InvoiceRequest.Payments {
		description := ""
		switch result.InvoiceRequest.Payments[i].PaymentType {
		case 0:
			description = "Other cashless payment"
		case 1:
			description = "Cash"
		case 2:
			description = "Payment card"
		case 3:
			description = "Cheque"
		case 4:
			description = "Bank transfer"
		case 5:
			description = "Voucher"
		case 6:
			description = "Instant payment"
		}

		result.InvoiceRequest.Payments[i].PaymentTypeDescription = description
	}

	return &result, nil
}

func ParseFiscalReceiptItems(receipt *FiscalReceiptInfo) ([]*FiscalReceiptItem, error) {
	journal := receipt.Journal
	lines := strings.Split(journal, "\n")
	var items []*FiscalReceiptItem
	for i, line := range lines {
		line = strings.TrimRight(line, "\r")
		if strings.Contains(line, "Назив") && strings.Contains(line, "Цена") && strings.Contains(line, "Кол.") && strings.Contains(line, "Укупно") {
			i++
			for ; i < len(lines); i += 2 {
				nameLine := strings.TrimRight(lines[i], "\r")
				if nameLine == "" || strings.HasPrefix(nameLine, "-") {
					break
				}
				if i+1 >= len(lines) {
					break
				}
				valLine := strings.TrimRight(lines[i+1], "\r")
				fields := strings.Fields(valLine)
				if len(fields) < 3 {
					continue
				}
				fields[0] = strings.ReplaceAll(fields[0], ",", ".")
				fields[1] = strings.ReplaceAll(fields[1], ",", ".")
				fields[2] = strings.ReplaceAll(fields[2], ",", ".")
				nameLine = strings.Trim(nameLine, " ")
				price, _ := strconv.ParseFloat(fields[0], 64)
				quantity, _ := strconv.ParseFloat(fields[1], 64)
				amount, _ := strconv.ParseFloat(fields[2], 64)
				item := &FiscalReceiptItem{
					Title:    nameLine,
					Price:    price,
					Quantity: quantity,
					Amount:   amount,
				}
				items = append(items, item)
			}
			break
		}
	}

	return items, nil
}
