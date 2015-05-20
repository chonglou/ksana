package ksana_orm

import (
	"container/list"
	"database/sql"
	"dmr"
	"io/ioutil"
	"log"
	"strings"
)

type Database struct {
	path       string
	driver     string
	config     Config
	dialect    Dialect
	db         *sql.DB
	migrations *list.List
}

func (d *Database) AddMigration(ver, up, down string) {
	d.migrations.PushBack(migration{version: ver, up: up, down: down})
}

//---------------------sql-----------------------------------------------------
func (d *Database) Column(name, flag string, int size, null bool, def string) string {
	// todo
}

func (d *Database) Create(name string) string {
	return d.dialect.CreateDatabase(d.config.Name)
}

func (d *Database) Drop() string {
	return d.dialect.DropDatabase(d.config.Name)
}

func (d *Database) Shell() string {
	return d.dialect.Shell(d.config)
}

func (d *Database) CreateTable(table string, columns ...string) string {
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(%s)", table, strings.Join(columns, ","))
}

func (d *Database) DropTable(table string) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s", table)
}

func (d *Database) CreateTable(name, table string, unique bool, columns ...string) string {
	idx := "INDEX"
	if unique {
		idx = "UNIQUE INDEX"
	}
	return fmt.Sprintf("CREATE %s ON %s (%s)", name, table, strings.Join(columns, ","))

}

func (d *Database) DropIndex(name string) string {
	return fmt.Sprintf("DROP INDEX %s")
}

//-------------------command---------------------------------------------------

func (d *Database) Migrate() error {
}

func (d *Database) Rollback() error {
}

func (d *Database) Generate(name string) error {
	// generate migrations file
	ver := time.Now().Format("20060102150405")
	return ioutil.WriteFile(fmt.Sprintf("db/migrate/%s_%s.go"), fmt.Sprintf(`
package CHANGE_ME
import (
	orm "github.com/chonglou/ksana/orm"
)

orm.DB.AddMigration("%s_%s", "CREATE TABLE T1(id INTEGER);", "DROP TABLE T1")
	`, ver, name), 0755)

}

func (d *Database) Open(cfg *Config, path string) error {
	logger.Info("Connect to database")
	db, err := sql.Open(driver, url)
	if err != nil {
		return err
	}
	logger.Info("Ping database")
	err = db.Ping()
	if err != nil {
		return err
	}

	d.db = db
	logger.Info("Database setup successfull")
	return nil

	// todo ping
	// todo check schema tables
}

var DB = Database{}
var logger = OpenLogger("ksana-orm")
