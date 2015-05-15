package ksana

import (
	"container/list"
	"database/sql"
	"errors"
	"log"
	"sync"
)

const migrations_table_name = "schema_migrations"

type migrationItem struct {
	version string
	up      string
	down    string
}

type migration struct {
	db    *sql.DB
	items *list.List
	lock  sync.RWMutex
}

func (m *migration) Add(v, u, d string) {
	m.items.PushBack(migrationItem{version: v, up: u, down: d})
}

func (m *migration) Migrate() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	v, e := m.version()
	if e != nil {
		return e
	}

	is := false
	if v == "" {
		is = true
		log.Printf("Empty database")
	} else {
		log.Printf("Current version: %s", v)
	}

	for it := m.items.Front(); it != nil; it = it.Next() {
		mi := it.Value.(migrationItem)

		if v == mi.version {
			is = true
			continue
		}

		if is {
			log.Printf("Migrate %s", mi.version)
			_, e := m.db.Exec(mi.up)
			if e != nil {
				return e
			}
			_, e = m.db.Exec(
				"INSERT INTO "+migrations_table_name+"(version) VALUES($1)",
				mi.version)
			if e != nil {
				return e
			}
		}

	}

	return nil
}

func (m *migration) Rollback() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	v, e := m.version()
	if e != nil {
		return e
	}
	if v == "" {
		return errors.New("Empty database")
	}

	for it := m.items.Front(); it != nil; it = it.Next() {
		mi := it.Value.(migrationItem)

		if v == mi.version {
			_, e = m.db.Exec(mi.down)
			if e != nil {
				return e
			}
			log.Printf("Rollback %s", mi.version)
			_, e = m.db.Exec(
				"DELETE FROM "+migrations_table_name+" WHERE version = $1",
				mi.version)
			return e
		}
	}

	return errors.New("Not find version: " + v)
}

func (m *migration) version() (string, error) {
	_, e := m.db.Exec(
		"CREATE TABLE IF NOT EXISTS " +
			migrations_table_name +
			"(id SERIAL, version VARCHAR(16) NOT NULL UNIQUE)")

	if e != nil {
		return "", e
	}

	var r *sql.Rows
	r, e = m.db.Query(
		"SELECT version FROM " +
			migrations_table_name +
			" ORDER BY id DESC LIMIT 1")
	if e != nil {
		return "", e
	}
	defer r.Close()

	if r.Next() {
		var v string
		e = r.Scan(&v)
		return v, e
	}
	return "", nil
}
