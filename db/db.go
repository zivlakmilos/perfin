package db

import (
	"strings"
	"unicode"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var dbInstance *sqlx.DB

func init() {
	sqlx.NameMapper = func(s string) string {
		var result strings.Builder
		for i, r := range s {
			if unicode.IsUpper(r) {
				if i > 0 {
					result.WriteRune('_')
				}
				result.WriteRune(unicode.ToLower(r))
			} else {
				result.WriteRune(r)
			}
		}
		return result.String()
	}
}

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
