package initiator

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewDB(cfg *Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.DbUrl)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func CloseDB(db *sql.DB) error {
	if db == nil {
		return nil
	}

	return db.Close()
}