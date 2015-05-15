package ksana

import (
	"database/sql"
)

type Model struct {
	db *sql.DB
}

func (m *Model) Ping() error {
	return m.db.Ping()
}
