package ksana

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
)

const migrations_table_name = "schema_migrations"

var re_sql_file = regexp.MustCompile("(?P<date>[\\d]{14})(?P<table>[_a-z0-9]+).sql$")

type migration struct {
	Up   string `json:"up"`
	Down string `json:"down"`
}

type Bean interface{}

type Model interface {
	Register(Bean) error
	Migrate() error
	Rollback() error
}

type model struct {
	db   *sql.DB
	path string
}

func (m *model) en(s string) string {
	return "_" + strings.ToLower(s)
}

func (m *model) column(
	table string,
	create, index *bytes.Buffer,
	tags map[string]string) {

	name := tags["name"]
	if tags["index"] != "" {
		unique := ""
		if tags["unique"] == "true" {
			unique = " UNIQUE"
		}

		idx1 := strings.Split(tags["index"], ",")
		idx2 := make([]string, len(idx1))
		for k, v := range idx1 {
			idx2[k] = m.en(v)
		}

		fmt.Fprintf(
			index, "CREATE%s INDEX %s_%s_idx ON %s (%s);",
			unique, table, name, table, strings.Join(idx2, ","))
	}

	def, null, fix, long, size := "", "", tags["fix"], tags["long"], tags["size"]
	if tags["default"] != "" {
		def = tags["default"]
	}
	if tags["null"] == "false" {
		null = " NOT NULL"
	}

	ty := tags["type"]
	switch ty {
	case "serial":
		fmt.Fprintf(create, "%s SERIAL", name)
		return
	case "uuid":
		fmt.Fprintf(create, "%s UUID NOT NULL PRIMARY KEY DEFAULT UUID_GENERATE_V4()", name)
		return
	case "created":
		fmt.Fprintf(create, "%s TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP", name)
		return
	case "updated":
		fmt.Fprintf(create, "%s TIMESTAMP", name)
		return
	case "date":
		ty = "DATE"
	case "time":
		ty = "TIME"
	case "datetime":
		ty = "TIMESTAMP"
	case "int":
		ty = "INTEGER"
	case "int8":
		ty = "SMALLINT"
	case "int32":
		ty = "INTEGER"
	case "int64":
		ty = "BIGINT"
	case "uint":
		ty = "INTEGER"
	case "string": // text char varchar byte bolb
		if def != "" {
			def = "'" + def + "'"
		}
		switch {
		case long == "true":
			ty = "TEXT"
		case fix == "true":
			ty = "CHAR(" + size + ")"
		default:
			ty = "VARCHAR(" + size + ")"
		}

	case "bool": //bool
		ty = "BOOLEAN"
		if def != "" {
			if def == "true" {
				def = "TRUE"
			} else {
				def = "FALSE"
			}
		}
	case "float32": //float
		ty = "REAL"
	case "float64": //double
		ty = "DOUBLE PRECISION"
	case "uint8": //byte
		ty = "BIT(1)"
	case "byte":
		switch {
		case long == "true":
			ty = "BYTEA"
		case fix == "true":
			ty = "BIT(" + size + ")"
		default:
			ty = "BIT VARYING(" + size + ")"
		}
	default:
		logger.Debug("Ingnore")
		return
	}

	if def != "" {
		def = " DEFAULT " + def
	}
	fmt.Fprintf(create, "%s %s%s%s", name, ty, null, def)
}

func (m *model) write(
	table string,
	create, index, drop *bytes.Buffer) error {

	fn := time.Now().Format("20060102150405") + table + ".sql"
	logger.Info("Generate model file: " + fn)

	cj, err := json.MarshalIndent(
		&migration{
			Up:   create.String() + index.String(),
			Down: drop.String()},
		"", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(m.path+"/"+fn, cj, 0600)
}

func (m *model) check(table string) (string, error) {

	files, err := ioutil.ReadDir(m.path)
	if err != nil {
		return "", err
	}

	for _, f := range files {
		fn := f.Name()
		if !re_sql_file.MatchString(fn) {
			continue
		}

		if table == re_sql_file.FindStringSubmatch(fn)[2] {
			return fn, nil
		}
	}
	return "", nil
}

func (m *model) Register(b Bean) error {

	bt := reflect.TypeOf(b)
	logger.Info("Load bean " + bt.Name())
	table := m.en(strings.Replace(bt.String(), ".", "_", -1))

	fn, err := m.check(table)
	if err != nil {
		return err
	}
	if fn != "" {
		logger.Info("Find migration " + fn)
		return nil
	}

	var create bytes.Buffer
	var drop bytes.Buffer
	var index bytes.Buffer

	fmt.Fprintf(&create, "CREATE TABLE IF NOT EXISTS %s", table)
	fmt.Fprintf(&drop, "DROP TABLE IF EXISTS %s;", table)

	for i := 1; i < bt.NumField(); i++ {
		f := bt.Field(i)
		tag := f.Tag.Get("sql")
		tags := make(map[string]string, 0)
		tags["type"] = f.Type.Name()
		tags["name"] = m.en(f.Name)

		if tag != "" {
			for _, it := range strings.Split(tag, ";") {
				ss := strings.Split(it, "=")
				if len(ss) != 2 {
					return errors.New("Error struct tag format: " + it)
				}
				tags[ss[0]] = ss[1]
			}
		}

		logger.Debug(fmt.Sprintf("Find field: %v", tags))

		if i == 1 {
			create.Write([]byte("("))
		} else {
			create.Write([]byte(", "))
		}

		m.column(table, &create, &index, tags)
	}

	create.Write([]byte(");"))
	return m.write(table, &create, &index, &drop)
}

func (m *model) version() error {
	var buf bytes.Buffer
	for _, p := range []string{"uuid-ossp", "pgcrypto"} {
		fmt.Fprintf(&buf, "CREATE EXTENSION IF NOT EXISTS \"%s\";", p)
	}
	fmt.Fprintf(
		&buf,
		"CREATE TABLE IF NOT EXISTS %s(id SERIAL, version VARCHAR(255) NOT NULL UNIQUE, created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);",
		migrations_table_name)

	_, err := m.db.Exec(buf.String())
	return err
}

func (m *model) read(mig *migration, file string) error {
	f, e := os.Open(m.path + "/" + file)
	if e != nil {
		return e
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(mig)
}

func (m *model) Migrate() error {
	err := m.version()
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(m.path)
	if err != nil {
		return err
	}

	for _, f := range files {

		fn := f.Name()
		var rs *sql.Rows

		rs, err = m.db.Query(fmt.Sprintf(
			"SELECT id FROM %s WHERE version = $1", migrations_table_name), fn)
		if err != nil {
			return err
		}
		defer rs.Close()

		if rs.Next() {
			log.Printf("Has %s", fn)
		} else {
			mig := migration{}
			err = m.read(&mig, fn)
			if err != nil {
				return err
			}
			log.Printf("Migrate version %s!\n%s", fn, mig.Up)
			_, err = m.db.Exec(mig.Up)
			if err != nil {
				return err
			}
			_, err = m.db.Exec(fmt.Sprintf(
				"INSERT INTO %s(version) VALUES($1)", migrations_table_name), fn)
			if err != nil {
				return err
			}

		}
	}
	return nil
}

func (m *model) Rollback() error {
	err := m.version()
	if err != nil {
		return err
	}

	var rs *sql.Rows
	rs, err = m.db.Query(fmt.Sprintf(
		"SELECT id, version FROM %s ORDER BY id DESC LIMIT 1", migrations_table_name))
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
		err = m.read(&mig, ver)
		if err != nil {
			return err
		}
		log.Printf("Rollback version %s\n%s", ver, mig.Down)
		_, err = m.db.Exec(mig.Down)
		if err != nil {
			return err
		}
		_, err = m.db.Exec(fmt.Sprintf(
			"DELETE FROM %s WHERE id=$1", migrations_table_name), id)
		if err != nil {
			return err
		}

	} else {
		log.Println("Empty database!")
	}

	return nil
}
