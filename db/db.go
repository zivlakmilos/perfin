package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var dbInstance *sqlx.DB

func CreateConnection(url string) error {
	if dbInstance == nil {
		con, err := sqlx.Open("sqlite3", url)
		if err != nil {
			return err
		}

		dbInstance = con
	}

	return nil
}

func GetInstance() *sqlx.DB {
	return dbInstance
}

func Open(url string) (*sqlx.DB, error) {
	return sqlx.Open("sqlite3", url)
}
