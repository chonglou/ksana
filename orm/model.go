package ksana_orm

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var re_sql_file = regexp.MustCompile("(?P<date>[\\d]{14})_(?P<table>[_a-zA-Z0-9]+).sql$")

type Model struct {
	db *Connection
}

func (m *Model) tags(field reflect.StructField) (map[string]string, error) {
	tag := field.Tag.Get("sql")

	if tag == "-" {
		return nil, nil
	}

	tags := make(map[string]string, 0)
	tags["name"] = field.Name

	if tag != "" {
		for _, it := range strings.Split(tag, ";") {
			ss := strings.Split(it, "=")
			if len(ss) != 2 {
				return nil, errors.New("Error struct tag format: " + it)
			}
			tags[ss[0]] = ss[1]
		}
	}

	return tags, nil
}

func (m *Model) table(bean interface{}) (string, reflect.Type) {
	bt := reflect.TypeOf(bean)
	return strings.Replace(bt.String(), ".", "_", -1), bt
}

func (m *Model) column(field reflect.StructField) (string, error) {
	tags, err := m.tags(field)
	if err != nil {
		return "", err
	}
	if tags == nil {
		return "", nil
	}

	switch tags["name"] {
	case "Id":
		switch field.Type.Kind() {
		case reflect.Int32:
			return m.db.Id(false), nil
		case reflect.String:
			return m.db.Id(true), nil
		default:
			return "", errors.New("Error type of id: " + field.Type.Name())
		}
	case "Created":
		return m.db.Created(), nil
	case "Updated":
		return m.db.Updated(), nil
	default:
		switch field.Type.Kind() {
		case reflect.Int32:
			val := 0
			if tags["default"] != "" {
				val, _ = strconv.Atoi(tags["default"])
			}
			return m.db.Int32(tags["name"], tags["null"] == "true", val), nil
		case reflect.Int64:
			var val int64
			if tags["default"] != "" {
				val, _ = strconv.ParseInt(tags["default"], 10, 64)
			}
			return m.db.Int64(tags["name"], tags["null"] == "true", val), nil
		case reflect.String:
			size, _ := strconv.Atoi(tags["size"])
			return m.db.String(
				tags["name"],
				tags["fix"] == "true",
				size,
				tags["big"] == "true",
				tags["null"] == "true",
				tags["default"]), nil

		case reflect.Bool:
			return m.db.Bool(tags["name"], tags["default"] == "true"), nil

		case reflect.Float32:
			var val float64
			if tags["default"] != "" {
				val, _ = strconv.ParseFloat(tags["default"], 32)
			}
			return m.db.Float32(tags["name"], float32(val)), nil
		case reflect.Float64:
			var val float64
			if tags["default"] != "" {
				val, _ = strconv.ParseFloat(tags["default"], 64)
			}
			return m.db.Float64(tags["name"], val), nil

		case reflect.Struct:
			ty := fmt.Sprintf("%s.%s", field.Type.PkgPath(), field.Type.Name())
			switch ty {
			case "time.Time":
				switch tags["type"] {
				case "date":
					return m.db.Date(
						tags["name"],
						tags["null"] == "true",
						tags["default"]), nil
				case "time":
					return m.db.Time(
						tags["name"],
						tags["null"] == "true",
						tags["default"]), nil
				default:
					return m.db.Datetime(
						tags["name"],
						tags["null"] == "true",
						tags["default"]), nil
				}
			default:
				return "", errors.New("Unsupport struct type " + ty)
			}

		case reflect.Slice:
			switch field.Type.Elem().Kind() {
			case reflect.Uint8:
				size, _ := strconv.Atoi(tags["size"])
				return m.db.Bytes(
					tags["name"],
					tags["fix"] == "true",
					size,
					tags["big"] == "true",
					tags["null"] == "true"), nil
			default:
				return "", errors.New("Unsupport slice type " + field.Type.Elem().Name())
			}

		default:
			// if _, ok := reflect.New(field.Type).Interface().(*[]byte); ok {
			//
			// }
			return "", errors.New("Ingnore column " + field.Name)
		}
	}

}

func (m *Model) Index(bean interface{}, unique bool, fields ...string) (string, string) {
	table, _ := m.table(bean)
	idx := fmt.Sprintf("%s_%s_idx", table, strings.Join(fields, "_"))
	return m.db.AddIndex(idx, table, unique, fields...), m.db.RemoveIndex(idx)

}

func (m *Model) Table(bean interface{}) (string, string, error) {
	table, bt := m.table(bean)
	logger.Info("Load bean " + bt.Name())

	var columns []string

	for i := 0; i < bt.NumField(); i++ {
		col, err := m.column(bt.Field(i))
		if err != nil {
			logger.Warning(err.Error())
			continue
		}
		if col != "" {
			columns = append(columns, col)
		}

	}

	return m.db.AddTable(table, columns...), m.db.RemoveTable(table), nil
}

func (m *Model) Check(path string, bean interface{}) (string, error) {
	table, _ := m.table(bean)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}

	for _, f := range files {
		fn := f.Name()

		if re_sql_file.MatchString(fn) {
			if table == re_sql_file.FindStringSubmatch(fn)[2] {
				return fn, nil
			}
		}

	}
	return "", nil
}

func (m *Model) Select(bean interface{}, columns []string, where, order string, offset, limit int, args ...interface{}) (*sql.Rows, error) {
	table, _ := m.table(bean)
	return m.db.Select(table, columns, where, order, offset, limit, args...)
}

func (m *Model) Delete(bean interface{}, where string, args ...interface{}) (sql.Result, error) {
	table, _ := m.table(bean)
	return m.db.Delete(table, where, args...)
}

func (m *Model) Update(bean interface{}, columns, where string, args ...interface{}) (sql.Result, error) {
	table, _ := m.table(bean)
	return m.db.Update(table, columns, where, args...)
}
