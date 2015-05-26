package ksana

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
	"path/filepath"
)

var migrations_table_name = "schema_migrations"

type Migrator interface {
	Add(version, name, up, down string) error
	Migrate() error
	Rollback() error
	Generate(name string) error
}

type migration struct {
	Version string `json:"version"`
	Up      string `json:"up"`
	Down    string `json:"down"`
}

type migrator struct {
	db   *sql.DB
	sql  *Sql
	path string
}

func (p *migrator) Add(version, name, up, down string) error {

	fn := fmt.Sprintf("%s/%s_%s.json", p.path, version, name)
	_, err := os.Stat(fn)
	if err == nil {
		logger.Info("Find migrations " + fn)
	} else {
		logger.Info("Generate migrations " + fn)

		cj, err := json.MarshalIndent(
			&migration{Version: version, Up: up, Down: down},
			"", "\t")

		if err != nil {
			return err
		}
		return ioutil.WriteFile(fn, cj, 0600)

	}
	return nil
}

func (p *migrator) read(mig *migration, file string) error {
	f, e := os.Open(p.path + "/" + file)
	if e != nil {
		return e
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(mig)
}

func (p *migrator) Migrate() error {
	files, err := ioutil.ReadDir(p.path)
	if err != nil {
		return err
	}

	for _, f := range files {
		fn := f.Name()
		var rs *sql.Rows

		rs, err = p.db.Query(fmt.Sprintf(
			"SELECT id FROM %s WHERE version = $1", migrations_table_name), fn)
		if err != nil {
			return err
		}
		defer rs.Close()

		if rs.Next() {
			log.Printf("Has %s", fn)
		} else {
			if filepath.Ext(fn) != ".json"{
				log.Printf("Ingnore file %s", fn)
				continue
			}
			mig := migration{}
			err = p.read(&mig, fn)
			if err != nil {
				return err
			}

			log.Printf("Migrate version %s\n%s", fn, mig.Up)
			_, err = p.db.Exec(mig.Up)
			if err != nil {
				return err
			}
			_, err = p.db.Exec(fmt.Sprintf(
				"INSERT INTO %s(version) VALUES($1)", migrations_table_name), fn)
			if err != nil {
				return err
			}

		}
	}
	log.Printf("Done!!!")

	return nil
}

func (p *migrator) Rollback() error {

	rs, err := p.db.Query(
		fmt.Sprintf("SELECT id, version FROM %s ORDER BY id DESC LIMIT 1",
			migrations_table_name))
	if err != nil {
		return err
	}
	defer rs.Close()

	if rs.Next() {
		var id int
		var ver string
		err = rs.Scan(&id, &ver)
		if err != nil {
			return nil
		}

		mig := migration{}
		err = p.read(&mig, ver)
		if err != nil {
			return err
		}
		log.Printf("Rollback version %s\n%s", ver, mig.Down)
		_, err = p.db.Exec(mig.Down)
		if err != nil {
			return err
		}
		_, err = p.db.Exec(fmt.Sprintf(
			"DELETE FROM %s WHERE id=$1", migrations_table_name), id)
		if err != nil {
			return err
		}

	} else {
		log.Println("Empty database!")
	}

	log.Println("Done!")

	return nil
}

func (p *migrator) version() string {
	return time.Now().Format("20060102150405")
}

func (p *migrator) Generate(name string) error {
	return p.Add(
		p.version(),
		name,
		p.sql.CreateTable(name, p.sql.Id(false)),
		p.sql.DropTable(name))
}

//-----------------NEW-----------------------
func NewMigrator(path string, cfg *databaseConfig) (Migrator, error) {

	db, sq, err := openDB(&dbCfg)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(path, 0700)
	if err != nil {
		return nil,err
	}

	logger.Info("Check migrations schema table")
	s := sq.CreateTable(migrations_table_name,
		sq.Id(false),
		sq.String("version", false, 255, false, false, ""),
		sq.Created())
	logger.Debug(s)
	_, err = db.Exec(s)
	if err != nil {
		return nil, err
	}

	return &migrator{path: path, db: db, sql: sq}, nil
}
