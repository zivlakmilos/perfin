package db

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ReceivedReceiptItem struct {
	Id        string  `db:"id" json:"id"`
	ReceiptId string  `db:"receipt_id" json:"receiptId"`
	Name      string  `db:"name" json:"name"`
	Quantity  float64 `db:"quantity" json:"quantity"`
	Amount    float64 `db:"amount" json:"amount"`
}

type ReceivedReceiptItemStore struct {
	con *sqlx.DB
}

func NewReceivedReceiptItemStore(con *sqlx.DB) *ReceivedReceiptItemStore {
	return &ReceivedReceiptItemStore{con: con}
}

func (s *ReceivedReceiptItemStore) Insert(item *ReceivedReceiptItem) error {
	if item.Id == "" {
		item.Id = uuid.NewString()
	}
	_, err := s.con.NamedExec(`INSERT INTO received_receipt_items (
		id,
		receipt_id,
		name,
		quantity,
		amount
	) VALUES (
		:id,
		:receipt_id,
		:name,
		:quantity,
		:amount
	)`, item)
	return err
}

func (s *ReceivedReceiptItemStore) GetAll() ([]*ReceivedReceiptItem, error) {
	var res []*ReceivedReceiptItem
	err := s.con.Select(&res, "SELECT * FROM received_receipt_items")
	return res, err
}
