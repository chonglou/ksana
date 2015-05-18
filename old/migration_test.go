package ksana

import (
	"container/list"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"testing"
)

func TestMigration(t *testing.T) {
	SQL.UseDriver("postgres")
	db, err := sql.Open("postgres", "postgres://postgres@localhost/ksana_t?sslmode=disable")

	if err != nil {
		t.Errorf("Open database: %v", err)
	}
	defer db.Close()

	var m Migration
	m = &migration{db: db, items: list.New()}

	m.Add("201505151051", "CREATE TABLE T1(f1 INT)", "DROP TABLE T1")
	m.Add("201505151052", "CREATE TABLE T2(f1 INT)", "DROP TABLE T2")
	m.Add("201505151053", "CREATE TABLE T3(f1 INT)", "DROP TABLE T3")
	m.Add("201505151054", "CREATE TABLE T4(f1 INT)", "DROP TABLE T4")
	m.Add("201505151055", "CREATE TABLE T5(f1 INT)", "DROP TABLE T5")

	log.Printf("==== Begin migrate ====")
	err = m.Migrate()
	if err != nil {
		t.Errorf("Migrate database: %v", err)
	}

	log.Printf("==== Begin rollback ====")
	err = m.Rollback()
	if err != nil {
		t.Errorf("Rollback database: %v", err)
	}
}
