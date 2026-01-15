package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Transaction struct {
	Id                string    `json:"id" xml:"id" form:"id" query:"id"`
	TransactionId     string    `json:"transactionId" xml:"transactionId" form:"transactionId" query:"transactionId"`
	AccountId         string    `json:"accountId" xml:"accountId" form:"accountId" query:"accountId"`
	Date              time.Time `json:"date" xml:"date" form:"date" query:"date"`
	Description       string    `json:"description" xml:"description" form:"description" query:"description"`
	Debit             float64   `json:"debit" xml:"debit" form:"debit" query:"debit"`
	Credit            float64   `json:"credit" xml:"credit" form:"credit" query:"credit"`
	ReceivedReceiptId string    `json:"receivedReceiptId" xml:"receivedReceiptId" form:"receivedReceiptId" query:"receivedReceiptId"`
}

type TransactionStore struct {
	con *sqlx.DB
	tx  *sqlx.Tx
}

func NewTransaction() *Transaction {
	return &Transaction{}
}

func NewTransactionStore(con *sqlx.DB) *TransactionStore {
	return &TransactionStore{
		con: con,
	}
}

func (s *TransactionStore) UseTransaction(tx *sqlx.Tx) {
	s.tx = tx
}

func (s *TransactionStore) Insert(m *Transaction) error {
	if m.Id == "" {
		m.Id = uuid.NewString()
	}

	query := `INSERT INTO transactions (
		id,
		transaction_id,
		account_id,
		date,
		description,
		debit,
		credit,
		received_receipt_id
	) VALUES (
		:id,
		:transaction_id,
		:account_id,
		:date,
		:description,
		:debit,
		:credit,
		:received_receipt_id
	)`

	if s.tx != nil {
		_, err := s.tx.NamedExec(query, m)
		if err != nil {
			return err
		}
	} else {
		_, err := s.con.NamedExec(query, m)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *TransactionStore) GetAll() ([]*Transaction, error) {
	var res []*Transaction
	err := s.con.Select(&res, "SELECT * FROM transactions")
	return res, err
}
