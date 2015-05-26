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

//---------------------sql-----------------------------------------------------
}

//-------------------command---------------------------------------------------



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

var logger, _ = utils.OpenLogger("ksana-orm")
