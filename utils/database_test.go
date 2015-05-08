package ksana

import (
	"testing"
	_ "github.com/lib/pq"
)

func TestDatabase(t *testing.T) {
	d := Database{}
	d.Open("postgres", "postgres://postgres@localhost/?sslmode=disable")
	rs, err := d.Query("select NOW()")
	var now string
	if rs.Next() {
		rs.Scan(&now)
	}
	err = rs.Err()

	if err != nil {
		t.Errorf("sql error")
	}
}
