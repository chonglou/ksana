package ksana_orm

import (
	"errors"
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"
	//"time"
)

var re_sql_file = regexp.MustCompile("(?P<date>[\\d]{14})(?P<table>[_a-z0-9]+).sql$")

type Model struct {
	bean  interface{}
	table string
}

func (m *Model) Table() string {
	return m.table
}

func (m *Model) check(path string) (string, error) {

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}

	for _, f := range files {
		fn := f.Name()
		if !re_sql_file.MatchString(fn) {
			continue
		}

		if m.table == re_sql_file.FindStringSubmatch(fn)[2] {
			return fn, nil
		}
	}
	return "", nil
}

func (m *Model) ingnore(tag string) bool {
	return tag == "-"
}

func (m *Model) column(field reflect.StructField) (string, string, error) {
	tag := field.Tag.Get("sql")
	if m.ingnore(tag) {
		return "", "", nil
	}

	col, idx := "", ""

	switch field.Type.Kind() {
	case reflect.String:
	case reflect.Bool:
	case reflect.Int:
	case reflect.Int64:
	case reflect.Float32:
	case reflect.Float64:
	case reflect.Struct:
		// if _, ok := field.Type.Interface().(time.Time); ok {
		//
		// }
	default:
		// if _, ok := field.Type.Interface().([]byte); ok {
		//
		// }
		return "", "", errors.New("Ingnore column " + field.Name)
	}
	return col, idx, nil
}

func (m *Model) Register(db *Database) error { // todo

	bt := reflect.TypeOf(m.bean)
	logger.Info("Load bean " + bt.Name())
	m.table = strings.Replace(bt.String(), ".", "_", -1)

	fn, err := m.check(db.path)
	if err != nil {
		return err
	}
	if fn != "" {
		logger.Info("Find migration " + fn)
		return nil
	}

	var columns, indexes []string

	for i := 0; i < bt.NumField(); i++ {
		col, idx, err := m.column(bt.Field(i))
		if err != nil {
			logger.Warning(err.Error())
			continue
		}
		if col != "" {
			columns = append(columns, col)
		}

		if idx != "" {
			indexes = append(indexes, idx)
		}
	}

	return db.AddMigration(
		db.version(), m.table,
		db.AddTable(m.table, columns...)+strings.Join(indexes, ""),
		db.RemoveTable(m.table))
}
