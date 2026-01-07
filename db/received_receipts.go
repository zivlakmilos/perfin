package db

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ReceivedReceipt struct {
	Id           string  `db:"id" json:"id"`
	TaxId        string  `db:"tax_id" json:"taxId"`
	BusinessName string  `db:"business_name" json:"businessName"`
	Date         string  `db:"date" json:"date"`
	TotalAmount  float64 `db:"total_amount" json:"totalAmount"`
	Url          string  `db:"url" json:"url"`
	Items        []*ReceivedReceiptItem
}

type ReceivedReceiptStore struct {
	con *sqlx.DB
}

func NewReceivedReceiptStore(con *sqlx.DB) *ReceivedReceiptStore {
	return &ReceivedReceiptStore{con: con}
}

func (s *ReceivedReceiptStore) Insert(r *ReceivedReceipt) error {
	if r.Id == "" {
		r.Id = uuid.NewString()
	}
	_, err := s.con.NamedExec(`INSERT INTO received_receipts (
		id,
		tax_id,
		business_name,
		date,
		total_amount,
		url
	) VALUES (
		:id,
		:tax_id,
		:business_name,
		:date,
		:total_amount,
		:url
	)`, r)

	return err
}

func (s *ReceivedReceiptStore) GetAll() ([]*ReceivedReceipt, error) {
	var res []*ReceivedReceipt

	err := s.con.Select(&res, "SELECT * FROM received_receipts")
	if err != nil {
		return nil, err
	}

	itemsStore := NewReceivedReceiptItemStore(s.con)
	items, err := itemsStore.GetAll()
	if err != nil {
		return nil, err
	}

	itemsByReceipt := make(map[string][]*ReceivedReceiptItem)
	for _, item := range items {
		itemsByReceipt[item.ReceiptId] = append(itemsByReceipt[item.ReceiptId], item)
	}

	for _, receipt := range res {
		receipt.Items = itemsByReceipt[receipt.Id]
	}

	return res, nil
}
