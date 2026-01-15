package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ReceivedReceipt struct {
	Id             string    `db:"id" json:"id"`
	TaxId          string    `db:"tax_id" json:"taxId"`
	BusinessName   string    `db:"business_name" json:"businessName"`
	Date           time.Time `db:"date" json:"date"`
	TotalAmount    float64   `db:"total_amount" json:"totalAmount"`
	PaymentAccount string    `db:"payment_account_id" json:"paymentAccount"`
	Url            string    `db:"url" json:"url"`
	Items          []*ReceivedReceiptItem
}

type ReceivedReceiptStore struct {
	con *sqlx.DB
	tx  *sqlx.Tx
}

func NewReceivedReceiptStore(con *sqlx.DB) *ReceivedReceiptStore {
	return &ReceivedReceiptStore{
		con: con,
	}
}

func (s *ReceivedReceiptStore) UseTransaction(tx *sqlx.Tx) {
	s.tx = tx
}

func (s *ReceivedReceiptStore) Insert(m *ReceivedReceipt) error {
	if m.Id == "" {
		m.Id = uuid.NewString()
	}

	query := `INSERT INTO received_receipts (
		id,
		tax_id,
		business_name,
		date,
		total_amount,
		payment_account_id,
		url
	) VALUES (
		:id,
		:tax_id,
		:business_name,
		:date,
		:total_amount,
		:payment_account_id,
		:url
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

	itemsStore := NewReceivedReceiptItemStore(s.con)
	itemsStore.UseTransaction(s.tx)
	for i := range m.Items {
		m.Items[i].ReceiptId = m.Id
		err := itemsStore.Insert(m.Items[i])
		if err != nil {
			return err
		}
	}

	return nil
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

func (s *ReceivedReceiptStore) Get(id string) (*ReceivedReceipt, error) {
	var res ReceivedReceipt

	err := s.con.Get(&res, "SELECT * FROM received_receipts WHERE id=:id", id)
	if err != nil {
		return nil, err
	}

	itemsStore := NewReceivedReceiptItemStore(s.con)
	items, err := itemsStore.GetAll()
	if err != nil {
		return nil, err
	}

	res.Items = append(res.Items, items...)

	return &res, nil
}
