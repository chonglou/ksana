package ksana

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var re_sql_file = regexp.MustCompile("(?P<date>[\\d]{14})(?P<table>[_a-z]+).sql$")

type Bean interface{}

type Migration interface {
	Add(Bean) error
	Migrate() error
	Rollback() error
}

type migration struct {
	path string
}

func (m *migration) en(s string) string {
	return "_" + strings.ToLower(s)
}

func (m *migration) column(
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

func (m *migration) write(
	table string,
	create, index, drop *bytes.Buffer) error {

	fn := time.Now().Format("20060102150405") + table + ".sql"
	logger.Info("Generate migration file: " + fn)

	sq := make(map[string]string, 0)
	sq["up"] = create.String() + index.String()
	sq["down"] = drop.String()

	cj, err := json.MarshalIndent(sq, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(m.path+"/"+fn, cj, 0600)
}

func (m *migration) check(b Bean) (string, error) {

	bt := reflect.TypeOf(b)
	logger.Info("Load bean " + bt.Name())
	table := m.en(strings.Replace(bt.String(), ".", "_", -1))

	files, err := ioutil.ReadDir(m.path)
	if err != nil {
		return "", err
	}

	for _, f := range files {
		fn := f.Name()
		logger.Debug("Find sql file: " + fn)
		if !re_sql_file.MatchString(fn) {
			logger.Warning("Error migration file: " + fn)
			continue
		}
		logger.Debug(re_sql_file.FindStringSubmatch(fn)[2])

		if table == re_sql_file.FindStringSubmatch(fn)[2] {
			return "", errors.New("Find migration file: " + fn)
		}
	}
	return table, nil
}

func (m *migration) Add(b Bean) error {
	table, err := m.check(b)
	if err != nil {
		return err
	}

	bt := reflect.TypeOf(b)
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

func (m *migration) Migrate() error {
	return nil
}

func (m *migration) Rollback() error {
	return nil
}
