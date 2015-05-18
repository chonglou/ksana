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

func (m *migration) Add(b Bean) error {
	bt := reflect.TypeOf(b)
	logger.Info("Load bean " + bt.Name())
	table := m.en(strings.Replace(bt.String(), ".", "_", -1))

	files, err := ioutil.ReadDir(m.path)
	if err != nil {
		return err
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
			logger.Info("Find migration file: " + fn)
			return nil
		}
	}

	fn := time.Now().Format("20060102150405") + table + ".sql"

	var up bytes.Buffer
	var down bytes.Buffer

	fmt.Fprintf(&up, "CREATE TABLE IF NOT EXISTS %s(", table)
	fmt.Fprintf(&down, "DROP TABLE IF EXISTS %s", table)

	for i := 0; i < bt.NumField(); i++ {
		f := bt.Field(i)
		fn := f.Name
		ft := f.Type.Name()
		tag := f.Tag.Get("sql")
		tags := make(map[string]string, 0)
		tags["type"] = ft
		if tag != "" {
			for _, it := range strings.Split(tag, ";") {
				ss := strings.Split(it, ":")
				if len(ss) != 2 {
					return errors.New("Error struct tag format: " + it)
				}
				tags[ss[0]] = ss[1]
			}
		}
		logger.Debug(fmt.Sprintf("Find field: %s %s %v", fn, ft, tags))

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
			logger.Debug("Ingnore field: " + fn + " " + ft)
		}

	}

	logger.Info("Generate migration file: " + fn)

	sq := make(map[string]string, 0)
	sq["up"] = up.String()
	sq["down"] = down.String()
	var cj []byte
	cj, err = json.MarshalIndent(sq, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(m.path+"/"+fn, cj, 0600)

}

func (m *migration) Migrate() error {
	return nil
}

func (m *migration) Rollback() error {
	return nil
}
