package ksana

import (
	"bytes"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Sql struct {
	driver string
}

func (s *Sql) Float(string name, null false, def string) {
	var buf bytes.Buffer
	buf.WriteString(name)
	switch s.driver {
	case "postgres":
		buf.WriteString(" REAL")
	default:
		buf.WriteString(" FLOAT")
	}
	if !null {
		buf.WriteString(" NOT NULL")
	}
	if def != "" {
		buf.WriteString(" DEFAULT ")
		buf.WriteString(def)
	}
	return buf.String()
}

func (s *Sql) Double(string name, null false, def string) {
	var buf bytes.Buffer
	buf.WriteString(name)
	switch s.driver {
	case "postgres":
		buf.WriteString(" DOUBLE PRECISION")
	default:
		buf.WriteString(" DOUBLE")
	}
	if !null {
		buf.WriteString(" NOT NULL")
	}
	if def != "" {
		buf.WriteString(" DEFAULT ")
		buf.WriteString(def)
	}
	return buf.String()
}

func (s *Sql) Blob(string name, null false) string {
	var buf bytes.Buffer
	buf.WriteString(name)
	switch s.driver {
	case "postgres":
		buf.WriteString(" BYTEA")
	default:
		buf.WriteString(" BLOB")
	}
	if !null {
		buf.WriteString(" NOT NULL")
	}
	return buf.String()
}

func (s *Sql) Byte(string name, size int, null false) string {
	var buf bytes.Buffer
	buf.WriteString(name)
	switch s.driver {
	case "postgres":
		buf.WriteString(" BIT")
	default:
		buf.WriteString(" VARBINARY")
	}
	buf.WriteString("(")
	buf.WriteString(size)
	buf.WriteString(")")
	if !null {
		buf.WriteString(" NOT NULL")
	}
	return buf.String()
}

func (s *Sql) Bool(string name, null false, def bool) string {
	var buf bytes.Buffer
	buf.WriteString(name)
	switch s.driver {
	case "postgres":
		buf.WriteString(" BOOLEAN")
	default:
		buf.WriteString(" TINYINT")
	}

	if !null {
		buf.WriteString(" NOT NULL")
	}
	if def != nil {
		buf.WriteString(" DEFAULT")
		switch s.driver {
		case "postgres":
			if def {
				buf.WriteString(" TRUE")
			} else {
				buf.WriteString(" FALSE")
			}
		default:
			if def {
				buf.WriteString(" 1")
			} else {
				buf.WriteString(" 0")
			}
		}
	}

	return buf.String()
}

func (s *Sql) Text(string name, null false, def string) string {
	var buf bytes.Buffer
	buf.WriteString(name)
	buf.WriteString(" TEXT")
	if !null {
		buf.WriteString(" NOT NULL")
	}
	if def != nil {
		buf.WriteString(" DEFAULT '")
		buf.WriteString(def)
		buf.WriteString("'")
	}

	return buf.String()
}

func (s *Sql) Time(string name, null false, def time.Time) string {
	var buf bytes.Buffer
	buf.WriteString(name)
	buf.WriteString(" TIME")
	if !null {
		buf.WriteString(" NOT NULL")
	}
	if def != nil {
		buf.WriteString(" DEFAULT '")
		buf.WriteString(def.Format("15:04:05"))
		buf.WriteString("'")
	}

	return buf.String()
}

func (s *Sql) Date(string name, null false, def time.Time) string {
	var buf bytes.Buffer
	buf.WriteString(name)
	buf.WriteString(" DATE")
	if !null {
		buf.WriteString(" NOT NULL")
	}
	if def != nil {
		buf.WriteString(" DEFAULT '")
		buf.WriteString(def.Format("2006-01-02"))
		buf.WriteString("'")
	}

	return buf.String()
}

func (s *Sql) Datetime(string name, null false, def time.Time) string {
	var buf bytes.Buffer
	buf.WriteString(name)
	switch s.driver {
	case "postgres":
		buf.WriteString(" TIMESTAMP ")
	default:
		buf.WriteString(" DATETIME ")
	}
	if !null {
		buf.WriteString(" NOT NULL ")
	}
	if def != nil {
		buf.WriteString(" DEFAULT '")
		buf.WriteString(def.Format("2006-01-02 15:04:05"))
		buf.WriteString("'")
	}

	return buf.String()
}

func (s *Sql) Created() string {
	switch s.driver {
	case "postgres":
		return "created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"
	default:
		return "created TIMESTAMP NOT NULL DEFAULT NOW()"
	}
}

func (s *Sql) Char(name string, size int, null bool, def string) string {
	var buf bytes.Buffer
	buf.WriteString(name)
	buf.WriteString(" CHAR")

	buf.WriteString("(")
	buf.WriteString(strconv.Itoa(size))
	buf.WriteString(")")

	if !null {
		buf.WriteString(" NOT NULL ")
	}
	if def != "" {
		buf.WriteString("DEFAULT '")
		buf.WriteString(def)
		buf.WriteString("'")
	}
	return buf.String()
}

func (s *Sql) String(name string, size int, null bool, def string) string {
	var buf bytes.Buffer
	buf.WriteString(name)
	buf.WriteString(" VARCHAR(")

	buf.WriteString(strconv.Itoa(size))
	buf.WriteString(")")

	if !null {
		buf.WriteString(" NOT NULL ")
	}

	if def != "" {
		buf.WriteString("DEFAULT '")
		buf.WriteString(def)
		buf.WriteString("'")
	}
	return buf.String()
}

func (s *Sql) Int64(name string, null bool, def string) string {

	var buf bytes.Buffer
	buf.WriteString(name)
	buf.WriteString(" BIGINT ")
	if !null {
		buf.WriteString("NOT NULL ")
	}
	if def != "" {
		buf.WriteString("DEFAULT ")
		buf.WriteString(def)
	}
	return buf.String()
}

func (s *Sql) Int32(name string, null bool, def string) string {

	var buf bytes.Buffer
	buf.WriteString(name)
	buf.WriteString(" INTEGER ")
	if !null {
		buf.WriteString("NOT NULL ")
	}
	if def != "" {
		buf.WriteString("DEFAULT ")
		buf.WriteString(def)
	}
	return buf.String()
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
			return "id BIGINT NOT NULL AUTO_INCREMENT"
		}
	}
}

func (s *Sql) CreateTable(name, id, columns ...string) {
	var buf bytes.Buffer
	buf.WriteString("CREATE TABLE IF NOT EXISTS ")
	buf.WriteString(name)
	buf.WriteString("(")
	buf.WriteString(strings.Join(columns, ","))
	buf.WriteString(")")
	return buf.String()
}

func (s *Sql) DropTable(name) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s", name)
}

func (s *Sql) CreateIndex(name, table string, unique bool, columns ...string) string {
	var buf bytes.Buffer
	buf.WriteString("CREATE ")
	if unique {
		buf.WriteString("UNIQUE ")
	}
	buf.WriteString("INDEX ")
	buf.WirteString(name)
	buf.WriteString("_idx ON ")
	buf.WriteString(table)
	buf.WriteString(" (")
	buf.WriteString(strings.Join(columns, ","))
	buf.WriteString(")")
	return buf.String()
}

func (s *Sql) DropIndex(name) string {
	return "DROP INDEX IF EXISTS %s_idx"
}

func (s *Sql) CreateDatabase(name) string {
	switch s.driver {
	case "mysql":
		return fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8", name)
	default:
		return fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", name)
	}

}

func (s *Sql) DropDatabase(name) string {
	return fmt.Sprintf("DROP DATABASE IF EXISTS %s", name)
}
