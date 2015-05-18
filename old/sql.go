package ksana

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Sql struct {
	driver string
}

//-------------------------column-----------------------------------------------
func (s *Sql) DateOf(t time.Time) string {
	return t.Format("2006-01-02")
}

func (s *Sql) TimeOf(t time.Time) string {
	return t.Format("15:04:05")
}

func (s *Sql) DatetimeOf(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func (s *Sql) Float(name string, m, d int, null bool, def string) string {

	switch s.driver {
	case "postgres":
		return s.column(name, "REAL", null, def)
	default:
		return s.column(name, "FLOAT("+strconv.Itoa(m)+","+strconv.Itoa(d)+")", null, def)
	}
}

func (s *Sql) Numeric(name string, m, d int, null bool, def string) string {
	switch s.driver {
	case "postgres":
		return s.column(name, "NUMERIC("+strconv.Itoa(m)+","+strconv.Itoa(d)+")", null, def)
	default:
		logger.Warning("Not support numeric on " + s.driver + ", using double instead!!!")
		return s.Double(name, m, d, null, def)
	}

}

func (s *Sql) Double(name string, m, d int, null bool, def string) string {

	switch s.driver {
	case "postgres":
		return s.column(name, "DOUBLE PRECISION", null, def)
	default:
		return s.column(name, "DOUBLE("+strconv.Itoa(m)+","+strconv.Itoa(d)+")", null, def)
	}
}

func (s *Sql) Byte(name string, size int, null bool) string {

	switch s.driver {
	case "postgres":
		return s.column(name, "BIT("+strconv.Itoa(size)+")", null, "")
	default:
		return s.column(name, "VARBINARY("+strconv.Itoa(size)+")", null, "")
	}
}

func (s *Sql) Blob(name string, null bool) string {
	switch s.driver {
	case "postgres":
		return s.column(name, "BYTEA", null, "")
	default:
		return s.column(name, "BLOB", null, "")
	}
}

func (s *Sql) Bool(name string, null bool, def bool) string {

	switch s.driver {
	case "postgres":
		if def {
			return s.column(name, "BOOLEAN", null, "TRUE")
		} else {
			return s.column(name, "BOOLEAN", null, "FALSE")
		}
	default:
		if def {
			return s.column(name, "TINYINT", null, "1")
		} else {
			return s.column(name, "TINYINT", null, "0")
		}
	}
}

func (s *Sql) Text(name string, null bool, def string) string {
	if def != "" {
		def = "'" + def + "'"
	}
	return s.column(name, "TEXT", null, def)
}

func (s *Sql) Time(name string, null bool, def string) string {
	if def != "" {
		def = "'" + def + "'"
	}
	return s.column(name, "DATE", null, def)
}

func (s *Sql) Date(name string, null bool, def string) string {

	if def != "" {
		def = "'" + def + "'"
	}
	return s.column(name, "DATE", null, def)
}

func (s *Sql) Datetime(name string, null bool, def string) string {

	if def != "" {
		def = "'" + def + "'"
	}
	switch s.driver {
	case "postgres":
		return s.column(name, "TIMESTAMP", null, def)
	default:
		return s.column(name, "DATETIME", null, def)
	}
}

func (s *Sql) Char(name string, size int, null bool, def string) string {
	if def != "" {
		def = "'" + def + "'"
	}
	return s.column(name, "CHAR("+strconv.Itoa(size)+")", null, def)
}

func (s *Sql) String(name string, size int, unique bool, null bool, def string) string {
	t := "VARCHAR(" + strconv.Itoa(size) + ")"
	if unique {
		t += " UNIQUE"
	}
	if def != "" {
		def = "'" + def + "'"
	}
	return s.column(name, t, null, def)
}

func (s *Sql) Long(name string, null bool, def int64) string {
	return s.column(name, "BIGINT", null, strconv.FormatInt(def, 10))
}

func (s *Sql) Int(name string, null bool, def int) string {
	return s.column(name, "INT", null, strconv.Itoa(def))
}

func (s *Sql) Short(name string, null bool, def int) string {
	return s.column(name, "SMALLINT", null, strconv.Itoa(def))
}

func (s *Sql) column(name string, _type string, null bool, def string) string {
	ns, ds := "", ""
	if !null {
		ns = " NOT NULL"
	}
	if def != "" {
		ds = " DEFAULT " + def
	}

	return fmt.Sprintf("%s %s%s%s", name, _type, ns, ds)
}

func (s *Sql) Created() string {
	switch s.driver {
	case "postgres":
		return "created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"
	default:
		return "created TIMESTAMP NOT NULL DEFAULT NOW()"
	}
}

func (s *Sql) Id(uuid bool) string {
	switch s.driver {
	case "postgres":
		if uuid {
			return "id UUID NOT NULL PRIMARY KEY DEFAULT UUID_GENERATE_V4()"
		} else {
			return "id SERIAL"
		}
	default:
		if uuid {
			return "id CHAR(36) PRIMARY KEY NOT NULL"
		} else {
			return "id INT NOT NULL AUTO_INCREMENT"
		}
	}
}

//---------------------------Sql------------------------------------------------
func (s *Sql) Insert(name string, columns ...string) string {
	vs := make([]string, len(columns))
	switch s.driver {
	case "postgres":
		for i, _ := range columns {
			vs[i] = fmt.Sprintf("$%d", i+1)
		}
	default:
		for i, _ := range columns {
			vs[i] = "?"
		}

	}
	return fmt.Sprintf(
		"INSERT INTO %s(%s) VALUES(%s)",
		name,
		strings.Join(columns, ", "),
		strings.Join(vs, ", "))
}

func (s *Sql) Delete(name, where string) string {
	if where != "" {
		where = " WHERE " + where
	}
	return fmt.Sprintf("DELETE FROM %s%s", name, where)
}
func (s *Sql) Order(name string, asc bool) string {
	if asc {
		return name + " ASC"
	}
	return name + " DESC"

}
func (s *Sql) Select(table string, columns []string, where, order string, offset, limit int) string {
	cs, ls := "*", ""
	if columns != nil {
		cs = strings.Join(columns, ",")
	}
	if where != "" {
		where = " WHERE " + where
	}

	if order != "" {
		order = " ORDER BY " + order
	}

	switch s.driver {
	case "postgres":
		if limit > 0 {
			ls = fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
		}
	case "mysql":
		if limit > 0 {
			ls = fmt.Sprintf(" LIMIT %d,%d", offset, limit)
		}
	default:
		logger.Warning("Not support page query for database " + s.driver)
	}
	return fmt.Sprintf("SELECT %s FROM %s%s%s%s", cs, table, where, order, ls)
}
func (s *Sql) CreateTable(buf *bytes.Buffer, name string, columns ...string) {
	fmt.Fprintf(buf,
		"CREATE TABLE IF NOT EXISTS %s(%s)",
		name, strings.Join(columns, ","))
}

func (s *Sql) DropTable(buf *bytes.Buffer, name string) {
	fmt.Fprintf(buf, "DROP TABLE IF EXISTS %s", name)
}

func (s *Sql) CreateIndex(buf *bytes.Buffer,
	name, table string, unique bool, columns ...string) {

	inx := "INDEX"
	if unique {
		inx = "UNIQUE " + inx
	}
	fmt.Fprintf(buf,
		"CREATE %s %s_inx ON %s (%s)",
		inx, name, table, strings.Join(columns, ","))
}

func (s *Sql) DropIndex(buf *bytes.Buffer, name string) {
	fmt.Fprintf(buf, "DROP INDEX IF EXISTS %s_idx", name)
}

func (s *Sql) CreateDatabase(buf *bytes.Buffer, name string) {
	switch s.driver {
	case "postgres":
		fmt.Fprintf(buf, "CREATE DATABASE %s ENCODING 'utf-8'", name)
	case "mySql":
		fmt.Fprintf(buf, "CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8", name)
	default:
		logger.Warning("Not support create database on " + s.driver)
	}

}

func (s *Sql) DropDatabase(buf *bytes.Buffer, name string) {
	switch s.driver {
	case "postgres", "mySql":
		fmt.Fprintf(buf, "DROP DATABASE IF EXISTS %s", name)
	default:
		logger.Warning("Not support drop database on " + s.driver)
	}

}

func (s *Sql) UseDriver(driver string) {
	s.driver = driver
}

//------------------------------------------------------------------------------
var SQL = Sql{}
