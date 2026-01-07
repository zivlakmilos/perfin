package db

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ItemAccountMapping struct {
	Id        string `db:"id" json:"id"`
	ItemName  string `db:"item_name" json:"itemName"`
	AccountId string `db:"account_id" json:"accountId"`
}

type ItemAccountMappingStore struct {
	con *sqlx.DB
}

func NewItemAccountMapping() *ItemAccountMapping {
	return &ItemAccountMapping{}
}

func NewItemAccountMappingStore(con *sqlx.DB) *ItemAccountMappingStore {
	return &ItemAccountMappingStore{
		con: con,
	}
}

func (s *ItemAccountMappingStore) Insert(m *ItemAccountMapping) error {
	if m.Id == "" {
		m.Id = uuid.NewString()
	}

	_, err := s.con.NamedExec(`INSERT INTO item_account_mapping (
		id,
		item_name,
		account_id
	) VALUES (
		:id,
		:item_name,
		:account_id
	)`, m)
	if err != nil {
		return err
	}

	return nil
}

func (s *ItemAccountMappingStore) GetAll() ([]*ItemAccountMapping, error) {
	var res []*ItemAccountMapping
	err := s.con.Select(&res, "SELECT * FROM item_account_mapping")
	return res, err
}

func (s *ItemAccountMappingStore) GetByItemName(itemName string) ([]*ItemAccountMapping, error) {
	var res []*ItemAccountMapping
	err := s.con.Select(&res, "SELECT * FROM item_account_mapping WHERE item_name = ?", itemName)
	return res, err
}
