package ksana

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	utils "github.com/chonglou/ksana/utils"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type Connection struct {
	path    string
	config  *Config
	dialect Dialect
	db      *sql.DB
}

func (d *Connection) Tx() (*sql.Tx, error) {
	return d.db.Begin()
}

func (d *Connection) AddMigration(ver, name, up, down string) error {
	fn := fmt.Sprintf("%s/%s_%s.sql", d.path, ver, name)
	_, err := os.Stat(fn)
	if err == nil {
		logger.Info("Find migrations " + fn)
	} else {
		logger.Info("Generate migrations " + fn)

		cj, err := json.MarshalIndent(
			&migration{Version: ver, Up: up, Down: down},
			"", "\t")

		if err != nil {
			return err
		}
		return ioutil.WriteFile(fn, cj, 0600)

	}
	return nil
}

//---------------------sql-----------------------------------------------------
}

//-------------------command---------------------------------------------------
func (m *Connection) readMigration(mig *migration, file string) error {
	f, e := os.Open(m.path + "/" + file)
	if e != nil {
		return e
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(mig)
}

func (d *Connection) Migrate() error {
	files, err := ioutil.ReadDir(d.path)
	if err != nil {
		return err
	}

	for _, f := range files {

		fn := f.Name()
		var rs *sql.Rows

		rs, err = d.db.Query(fmt.Sprintf(
			"SELECT id FROM %s WHERE version = $1", migrations_table_name), fn)
		if err != nil {
			return err
		}
		defer rs.Close()

		if rs.Next() {
			log.Printf("Has %s", fn)
		} else {
			mig := migration{}
			err = d.readMigration(&mig, fn)
			if err != nil {
				return err
			}
			log.Printf("Migrate version %s\n%s", fn, mig.Up)
			_, err = d.db.Exec(mig.Up)
			if err != nil {
				return err
			}
			_, err = d.db.Exec(fmt.Sprintf(
				"INSERT INTO %s(version) VALUES($1)", migrations_table_name), fn)
			if err != nil {
				return err
			}

		}
	}
	log.Printf("Done!!!")

	return nil
}

func (d *Connection) Rollback() error {

	rs, err := d.db.Query(
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
		err = d.readMigration(&mig, ver)
		if err != nil {
			return err
		}
		log.Printf("Rollback version %s\n%s", ver, mig.Down)
		_, err = d.db.Exec(mig.Down)
		if err != nil {
			return err
		}
		_, err = d.db.Exec(fmt.Sprintf(
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
func (d *Connection) Version() string {
	return time.Now().Format("20060102150405")
}
func (d *Connection) Generate(name string) error {
	return d.AddMigration(
		d.Version(),
		name,
		d.AddTable(name, d.Id(false), d.Created()),
		d.RemoveTable(name))

}

func (d *Connection) Open(path string, cfg *Config) error {
	err := os.MkdirAll(path, 0700)
	if err != nil {
		return err
	}
	d.path = path

	switch cfg.Driver {
	case "postgres":
		d.dialect = &pgDialect{}
	default:
		return errors.New("Not supported driver " + cfg.Driver)
	}

	logger.Info("Connect to database " + d.dialect.String(cfg))

	var db *sql.DB
	db, err = sql.Open(cfg.Driver, d.dialect.Resource(cfg))
	if err != nil {
		return err
	}

	logger.Info("Ping database")
	err = db.Ping()
	if err != nil {
		return err
	}

	logger.Info("Run setup scripts")
	_, err = db.Exec(d.dialect.Setup())
	if err != nil {
		return err
	}

	logger.Info("Check migrations schema table")
	sq := d.AddTable(migrations_table_name,
		d.Id(false),
		d.String("version", false, 255, false, false, ""),
		d.Created())
	logger.Debug(sq)
	_, err = db.Exec(sq)
	if err != nil {
		return err
	}

	d.db = db
	d.config = cfg

	logger.Info("Connection setup successfull")
	return nil

}

//curd
func (d *Connection) Select(table string, columns []string, where, order string, offset, limit int, args ...interface{}) (*sql.Rows, error) {
	sq := d.dialect.Select(table, columns, where, order, offset, limit)
	logger.Debug(sq)
	return d.db.Query(sq, args...)
}

func (d *Connection) Delete(table, where string, args ...interface{}) (sql.Result, error) {
	sq := fmt.Sprintf("DELETE FROM %s WHERE %s", table, where)
	logger.Debug(sq)
	return d.db.Exec(sq, args...)
}

func (d *Connection) Update(table, columns, where string, args ...interface{}) (sql.Result, error) {
	sq := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, columns, where)
	logger.Debug(sq)
	return d.db.Exec(sq, args...)
}

//-----------------------------------------------------------------------------
var migrations_table_name = "schema_migrations"
var logger, _ = utils.OpenLogger("ksana-orm")
