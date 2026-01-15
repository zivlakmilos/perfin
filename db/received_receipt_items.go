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
	tx  *sqlx.Tx
}

func NewReceivedReceiptItemStore(con *sqlx.DB) *ReceivedReceiptItemStore {
	return &ReceivedReceiptItemStore{
		con: con,
	}
}

func (s *ReceivedReceiptItemStore) UseTransaction(tx *sqlx.Tx) {
	s.tx = tx
}

func (s *ReceivedReceiptItemStore) Insert(item *ReceivedReceiptItem) error {
	if item.Id == "" {
		item.Id = uuid.NewString()
	}
	query := `INSERT INTO received_receipt_items (
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
	)`

	if s.tx != nil {
		_, err := s.tx.NamedExec(query, item)
		if err != nil {
			return err
		}
	} else {
		_, err := s.con.NamedExec(query, item)
		if err != nil {
			return err
		}
	}

	return nil
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
