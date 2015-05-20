package ksana_orm

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chonglou/ksana/utils"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type Database struct {
	path    string
	config  *Config
	dialect Dialect
	db      *sql.DB
}

func (d *Database) Tx() (*sql.Tx, error) {
	return d.db.Begin()
}

func (d *Database) AddMigration(ver, name, up, down string) error {
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
func (d *Database) Created() string {
	return d.column("created", d.dialect.DATETIME(), false, d.dialect.Now())
}

func (d *Database) Updated() string {
	return d.column("updated", d.dialect.DATETIME(), true, "")
}

func (d *Database) Id(uuid bool) string {
	if uuid {
		return fmt.Sprintf(
			"id %s NOT NULL PRIMARY KEY DEFAULT %s",
			d.dialect.UUID(), d.dialect.Uuid())
	}
	return fmt.Sprintf("id %s", d.dialect.SERIAL())
}

func (d *Database) Bool(name string, def bool) string {
	return d.column(name, d.dialect.BOOLEAN(), false, d.dialect.Boolean(def))
}

func (d *Database) String(name string, fix bool, size int, big, null bool, def string) string {
	var ts string
	switch {
	case big:
		ts = "TEXT"
	case fix:
		ts = fmt.Sprintf("CHAR(%d)", size)
	default:
		ts = fmt.Sprintf("VARCHAR(%d)", size)
	}
	if def != "" {
		def = fmt.Sprintf("'%s'", def)
	}
	return d.column(name, ts, null, def)
}

func (d *Database) Int32(name string, null bool, def int) string {
	return d.column(name, "INT", null, fmt.Sprintf("%d", def))
}

func (d *Database) Int64(name string, null bool, def int64) string {
	return d.column(name, "BIGINT", null, fmt.Sprintf("%d", def))
}

func (d *Database) Bytes(name string, fix bool, size int, big, null bool) string {
	if big {
		return d.column(name, d.dialect.BLOB(), null, "")
	} else {
		return d.column(name, d.dialect.BYTES(fix, size), null, "")
	}

}

func (d *Database) Date(name string, null bool, def string) string {
	var ds string
	switch def {
	case "now":
		ds = d.dialect.CurDate()
	default:
		ds = def
	}
	return d.column(name, "DATE", null, ds)
}

func (d *Database) Time(name string, null bool, def string) string {
	var ds string
	switch def {
	case "now":
		ds = d.dialect.CurTime()
	default:
		ds = def
	}
	return d.column(name, "TIME", null, ds)
}

func (d *Database) Datetime(name string, null bool, def string) string {
	var ds string
	switch def {
	case "now":
		ds = d.dialect.Now()
	default:
		ds = def
	}
	return d.column(name, d.dialect.DATETIME(), null, ds)
}

func (d *Database) Float32(name string, null bool, def float32) string {
	return d.column(name, d.dialect.FLOAT(), null, fmt.Sprintf("%f", def))
}

func (d *Database) Float64(name string, null bool, def float64) string {
	return d.column(name, d.dialect.DOUBLE(), null, fmt.Sprintf("%f", def))
}

func (d *Database) column(name string, _type string, null bool, def string) string {
	ns, ds := "", ""
	if !null {
		ns = " NOT NULL"
	}
	if def != "" {
		ds = fmt.Sprintf(" DEFAULT %s", def)
	}
	return fmt.Sprintf("%s %s%s%s", name, _type, ns, ds)
}

func (d *Database) AddTable(table string, columns ...string) string {
	return fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s(%s)", table, strings.Join(columns, ", "))
}

func (d *Database) RemoveTable(table string) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s", table)
}

func (d *Database) AddIndex(name, table string, unique bool, columns ...string) string {
	idx := "INDEX"
	if unique {
		idx = "UNIQUE INDEX"
	}
	return fmt.Sprintf(
		"CREATE %s %s ON %s (%s)", idx, name, table, strings.Join(columns, ", "))

}

func (d *Database) RemoveIndex(name string) string {
	return fmt.Sprintf("DROP INDEX %s")
}

func (d *Database) Create(name string) string {
	return d.dialect.CreateDatabase(d.config.Name)
}

func (d *Database) Drop() string {
	return d.dialect.DropDatabase(d.config.Name)
}

func (d *Database) Shell() error {
	cmd, args := d.dialect.Shell(d.config)
	return ksana_utils.Shell(cmd, args...)
}

//-------------------command---------------------------------------------------
func (m *Database) readMigration(mig *migration, file string) error {
	f, e := os.Open(m.path + "/" + file)
	if e != nil {
		return e
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(mig)
}

func (d *Database) Migrate() error {
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

func (d *Database) Rollback() error {

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
func (d *Database) version() string {
	return time.Now().Format("20060102150405")
}
func (d *Database) Generate(name string) error {
	return d.AddMigration(
		d.version(),
		name,
		d.AddTable(name, d.Id(false), d.Created()),
		d.RemoveTable(name))

}

func (d *Database) Open(path string, cfg *Config) error {
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

	logger.Info("Database setup successfull")
	return nil

}

//-----------------------------------------------------------------------------
var migrations_table_name = "schema_migrations"
var logger, _ = ksana_utils.OpenLogger("ksana-orm")
