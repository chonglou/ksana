package ksana

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type pgDialect struct {
	config *databaseConfig
}

func (d *pgDialect) SERIAL() string {
	return "SERIAL"
}

func (d *pgDialect) UUID() string {
	return "UUID"
}

func (d *pgDialect) BOOLEAN() string {
	return "BOOLEAN"
}

func (d *pgDialect) FLOAT() string {
	return "REAL"
}

func (d *pgDialect) DOUBLE() string {
	return "DOUBLE PRECISION"
}

func (d *pgDialect) BLOB() string {
	return "BYTES"
}

func (d *pgDialect) BYTES(fix bool, size int) string {
	if fix {
		return fmt.Sprintf("BIT(%d)", size)
	}
	return fmt.Sprintf("BIT VARYING(%d)", size)
}

func (d *pgDialect) DATETIME() string {
	return "TIMESTAMP"
}

func (d *pgDialect) CurDate() string {
	return "CURRENT_DATE"
}

func (d *pgDialect) CurTime() string {
	return "CURRENT_TIME"
}

func (d *pgDialect) Now() string {
	return "CURRENT_TIMESTAMP"
}

func (d *pgDialect) Uuid() string {
	return "UUID_GENERATE_V4()"
}

func (d *pgDialect) Boolean(val bool) string {
	if true {
		return "TRUE"
	}
	return "FALSE"
}

func (d *pgDialect) CreateDatabase(name string) string {
	return fmt.Sprintf("CREATE DATABASE %s ENCODING=UTF8", name)
}

func (d *pgDialect) DropDatabase(name string) string {
	return fmt.Sprintf("DROP DATABASE %s", name)
}

func (d *pgDialect) Resource() string {
	return fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s?sslmode=%s",
		d.config.Driver, d.config.User, d.config.Password, d.config.Host, d.config.Port, d.config.Name, d.config.Ssl)

}

func (d *pgDialect) Shell() (string, []string) {
	return "psql", []string{
		"-h", d.config.Host,
		"-p", strconv.Itoa(d.config.Port),
		"-d", d.config.Name,
		"-U", d.config.User}
}

func (d *pgDialect) Setup() string {
	var buf bytes.Buffer
	for _, p := range []string{"uuid-ossp", "pgcrypto"} {
		fmt.Fprintf(&buf, "CREATE EXTENSION IF NOT EXISTS \"%s\";", p)
	}
	return buf.String()
}

func (d *pgDialect) String() string {
	return fmt.Sprintf("%s@%s:%d/%s", d.config.User, d.config.Host, d.config.Port, d.config.Name)
}

func (d *pgDialect) Select(table string, columns []string, where, order string, offset, limit int) string {
	if order != "" {
		order = fmt.Sprintf(" ORDER BY %s", order)
	}

	ofs := ""
	if offset > 0 {
		ofs = fmt.Sprintf(" OFFSET %d", offset)
	}

	lis := ""
	if limit > 0 {
		ofs = fmt.Sprintf(" LIMIT %d", offset)
	}

	return fmt.Sprintf("SELECT %s FROM %s WHERE %s%s%s%s", strings.Join(columns, ", "), table, where, order, ofs, lis)
}
