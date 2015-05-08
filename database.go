package ksana

import (
	"database/sql"
	"log"
)

type DatabaseCfg struct {
	Adapter  string `json:"adapter"`
	Url  string   `json:"url"`
}

type Database struct {
	db *sql.DB
}

func (d *Database) Open(adapter string, url string) {
	db, err := sql.Open(adapter, url)
	if err != nil {
		log.Fatalf("Error on open database connect: %v", err)
	}
	d.db = db
}

func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := d.db.Exec(query, args...)
	if err != nil {
		log.Printf("Error on run sql: %v\n%s", query)
	}
	return result, err
}

func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := d.db.Query(query, args...)
	if err != nil {
		log.Printf("Error on query: %v\n%s", err, query)
	}
	return rows, err

}
