package ksana

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var re_sql_file = regexp.MustCompile("(?P<date>[\\d]{14})_(?P<table>[_a-zA-Z0-9]+).json$")

type Model interface {
	Index(bean interface{}, unique bool, fields ...string) (string, string)
	Table(bean interface{}) (string, string, error)
	Check(path string, bean interface{}) (string, error)
}

type model struct {
	sql *Sql
}

func (m *model) tags(field reflect.StructField) (map[string]string, error) {
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

func (m *model) table(bean interface{}) (string, reflect.Type) {
	bt := reflect.TypeOf(bean)
	return strings.Replace(bt.String(), ".", "_", -1), bt
}

func (m *model) column(field reflect.StructField) (string, error) {
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
			return m.sql.Id(false), nil
		case reflect.String:
			return m.sql.Id(true), nil
		default:
			return "", errors.New("Error type of id: " + field.Type.Name())
		}
	case "Created":
		return m.sql.Created(), nil
	case "Updated":
		return m.sql.Updated(), nil
	default:
		switch field.Type.Kind() {
		case reflect.Int32:
			val := 0
			if tags["default"] != "" {
				val, _ = strconv.Atoi(tags["default"])
			}
			return m.sql.Int32(tags["name"], tags["null"] == "true", val), nil
		case reflect.Int64:
			var val int64
			if tags["default"] != "" {
				val, _ = strconv.ParseInt(tags["default"], 10, 64)
			}
			return m.sql.Int64(tags["name"], tags["null"] == "true", val), nil
		case reflect.String:
			size, _ := strconv.Atoi(tags["size"])
			return m.sql.String(
				tags["name"],
				tags["fix"] == "true",
				size,
				tags["big"] == "true",
				tags["null"] == "true",
				tags["default"]), nil

		case reflect.Bool:
			return m.sql.Bool(tags["name"], tags["default"] == "true"), nil

		case reflect.Float32:
			var val float64
			if tags["default"] != "" {
				val, _ = strconv.ParseFloat(tags["default"], 32)
			}
			return m.sql.Float32(tags["name"], float32(val)), nil
		case reflect.Float64:
			var val float64
			if tags["default"] != "" {
				val, _ = strconv.ParseFloat(tags["default"], 64)
			}
			return m.sql.Float64(tags["name"], val), nil

		case reflect.Struct:
			ty := fmt.Sprintf("%s.%s", field.Type.PkgPath(), field.Type.Name())
			switch ty {
			case "time.Time":
				switch tags["type"] {
				case "date":
					return m.sql.Date(
						tags["name"],
						tags["null"] == "true",
						tags["default"]), nil
				case "time":
					return m.sql.Time(
						tags["name"],
						tags["null"] == "true",
						tags["default"]), nil
				default:
					return m.sql.Datetime(
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
				return m.sql.Bytes(
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

func (m *model) Index(bean interface{}, unique bool, fields ...string) (string, string) {
	table, _ := m.table(bean)
	idx := fmt.Sprintf("%s_%s_idx", table, strings.Join(fields, "_"))
	return m.sql.CreateIndex(idx, table, unique, fields...), m.sql.DropIndex(idx)

}

func (m *model) Table(bean interface{}) (string, string, error) {
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

	return m.sql.CreateTable(table, columns...), m.sql.DropTable(table), nil
}

func (m *model) Check(path string, bean interface{}) (string, error) {
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
