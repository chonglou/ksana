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

	switch tags["type"] {
	case "serial":
	case "uuid":
	case "created":

	case "date":
	case "time":
	case "datetime":

	case "int":
	case "int8":
	case "int32":
	case "int64":
	case "uint":
	case "uint8": //byte
	case "string": // text char varchar byte bolb
	case "bool": //bool
	case "float32": //float
	case "float64": //double
	default:
		logger.Debug("Ingnore")
	}

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

	for i := 0; i < bt.NumField(); i++ {
		f := bt.Field(i)
		tag := f.Tag.Get("sql")
		tags := make(map[string]string, 0)
		tags["type"] = f.Type.Name()
		tags["name"] = f.Name

		if tag != "" {
			for _, it := range strings.Split(tag, ";") {
				ss := strings.Split(it, ":")
				if len(ss) != 2 {
					return errors.New("Error struct tag format: " + it)
				}
				tags[ss[0]] = ss[1]
			}
		}

		logger.Debug(fmt.Sprintf("Find field: %v", tags))

		if i == 0 {
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
