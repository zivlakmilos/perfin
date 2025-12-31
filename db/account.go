package db

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AccountType string

const (
	AccountTypeAsset     AccountType = "asset"
	AccountTypeLiability AccountType = "liability"
	AccountTypeIncome    AccountType = "income"
	AccountTypeExpense   AccountType = "expense"
	AccountTypeEquity    AccountType = "equity"
)

type Account struct {
	Id          string      `json:"id" xml:"id" form:"id" query:"id"`
	AccountType AccountType `json:"accountType" xml:"accountType" form:"accountType" query:"accountType"`
	ParentId    string      `json:"parentId" xml:"parentId" form:"parentId" query:"parentId"`
	Title       string      `json:"title" xml:"title" form:"title" query:"title"`
	Childrens   []Account   `json:"childrens" xml:"childrens" form:"childrens" query:"childrens"`
}

type AccountStore struct {
	con *sqlx.DB
}

func NewAccount() *Account {
	return &Account{}
}

func NewAccountStore(con *sqlx.DB) *AccountStore {
	return &AccountStore{
		con: con,
	}
}

func (s *AccountStore) Insert(m *Account) error {
	if m.Id == "" {
		m.Id = uuid.NewString()
	}

	_, err := s.con.NamedExec(`INSERT INTO accounts (
		id,
		account_type,
		parent_id,
		title
	) VALUES (
		:id,
		:account_type,
		:parent_id,
		:title
	)`, m)
	if err != nil {
		return err
	}

	return nil
}

func (s *AccountStore) GetAll() ([]*Account, error) {
	var res []*Account
	err := s.con.Select(&res, "SELECT * FROM accounts")
	return res, err
}
