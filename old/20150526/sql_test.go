package ksana

import (
	"database/sql"
	_ "github.com/lib/pq"
	"testing"
)

func TestSql(t *testing.T) {
	db, err := sql.Open("postgres", "postgres://postgres@localhost/ksana?sslmode=disable")
	if err != nil {
		t.Errorf("Open database: %v", err)
	}
	defer db.Close()

}
