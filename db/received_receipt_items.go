package db

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ReceivedReceiptItem struct {
	Id        string  `db:"id" json:"id"`
	ReceiptId string  `db:"receipt_id" json:"receiptId"`
	Name      string  `db:"name" json:"name"`
	Price     float64 `db:"price" json:"price"`
	Quantity  float64 `db:"quantity" json:"quantity"`
	Amount    float64 `db:"amount" json:"amount"`
	Account   string  `db:"account_id" json:"account"`
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
		price,
		quantity,
		amount,
		account_id
	) VALUES (
		:id,
		:receipt_id,
		:name,
		:price,
		:quantity,
		:amount,
		:account_id
	)`, item)
	return err
}

func (s *ReceivedReceiptItemStore) InsertWithTx(tx *sqlx.Tx, m *ReceivedReceiptItem) error {
	if m.Id == "" {
		m.Id = uuid.NewString()
	}
	_, err := tx.NamedExec(`INSERT INTO received_receipt_items (
		id,
		receipt_id,
		name,
		price,
		quantity,
		amount,
		account_id
	) VALUES (
		:id,
		:receipt_id,
		:name,
		:price,
		:quantity,
		:amount,
		:account_id
	)`, m)
	return err
}

func (s *ReceivedReceiptItemStore) GetAll() ([]*ReceivedReceiptItem, error) {
	var res []*ReceivedReceiptItem
	err := s.con.Select(&res, "SELECT * FROM received_receipt_items")
	return res, err
}

func (s *ReceivedReceiptItemStore) GetAllForReceipt(receiptId string) ([]*ReceivedReceiptItem, error) {
	var res []*ReceivedReceiptItem
	err := s.con.Select(&res, "SELECT * FROM received_receipt_items WHERE receipt_id=:receipt_id", receiptId)
	return res, err
}
