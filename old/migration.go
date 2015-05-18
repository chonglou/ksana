package ksana

import (
	"bytes"
	"container/list"
	"database/sql"
	"errors"
	"log"
	"sync"
)

const migrations_table_name = "schema_migrations"

type Migration interface {
	Add(v, u, d string)
	Migrate() error
	Rollback() error
}

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
	m.items.PushBack(migrationItem{
		version: v,
		up:      u,
		down:    d})
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
			log.Println(mi.up)
			_, e := m.db.Exec(mi.up)
			if e != nil {
				return e
			}

			sq := SQL.Insert(migrations_table_name, "version")
			logger.Info(sq)
			_, e = m.db.Exec(sq, mi.version)
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
	log.Printf("Current version: %s", v)

	for it := m.items.Front(); it != nil; it = it.Next() {
		mi := it.Value.(migrationItem)

		if v == mi.version {
			log.Println(mi.down)
			_, e = m.db.Exec(mi.down)
			if e != nil {
				return e
			}
			log.Printf("Rollback %s", mi.version)

			s1 := SQL.Delete(migrations_table_name, "version = $1")
			logger.Info(s1)
			_, e = m.db.Exec(s1, mi.version)
			return e
		}
	}

	return errors.New("Not find version: " + v)
}

func (m *migration) version() (string, error) {
	var buf bytes.Buffer
	SQL.CreateTable(&buf, migrations_table_name, SQL.Id(false), SQL.String("version", 16, true, false, ""))
	sq := buf.String()
	logger.Debug(sq)

	_, e := m.db.Exec(sq)
	if e != nil {
		return "", e
	}

	var r *sql.Rows

	sq = SQL.Select(migrations_table_name, []string{"version"}, "", SQL.Order("id", false), 0, 1)
	logger.Debug(sq)

	r, e = m.db.Query(sq)
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
