package db

import "github.com/jmoiron/sqlx"

type Config struct {
	Key   string `json:"key" xml:"key" form:"key" query:"key"`
	Value string `json:"value" xml:"value" form:"value" query:"value"`
}

type ConfigStore struct {
	con *sqlx.DB
}

func NewConfig() *Config {
	return &Config{}
}

func NewConfigStore(con *sqlx.DB) *ConfigStore {
	return &ConfigStore{
		con: con,
	}
}

func (store *ConfigStore) GetConfig() map[string]string {
	rows, err := store.con.Queryx("SELECT key, value FROM config")
	if err != nil {
		return nil
	}
	defer func() { _ = rows.Close() }()

	configs := make(map[string]string)
	for rows.Next() {
		var cfg Config
		if err := rows.StructScan(&cfg); err != nil {
			continue
		}
		configs[cfg.Key] = cfg.Value
	}

	return configs
}

func (store *ConfigStore) ReplaceConfig(cfg map[string]string) error {
	tx, err := store.con.Beginx()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.Exec("DELETE FROM config")
	if err != nil {
		return err
	}

	stmt, err := tx.Preparex("INSERT INTO config (key, value) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer func() { _ = stmt.Close() }()

	for k, v := range cfg {
		_, err := stmt.Exec(k, v)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (store *ConfigStore) SetValue(key, value string) error {
	_, err := store.con.Exec(`
		INSERT INTO config (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value=excluded.value
	`, key, value)
	return err
}
