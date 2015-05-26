package ksana

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type pgDialect struct {
	config *dbConfig
}

func (p *pgDialect) SERIAL() string {
	return "SERIAL"
}

func (p *pgDialect) UUID() string {
	return "UUID"
}

func (p *pgDialect) BOOLEAN() string {
	return "BOOLEAN"
}

func (p *pgDialect) FLOAT() string {
	return "REAL"
}

func (p *pgDialect) DOUBLE() string {
	return "DOUBLE PRECISION"
}

func (p *pgDialect) BLOB() string {
	return "BYTES"
}

func (p *pgDialect) BYTES(fix bool, size int) string {
	if fix {
		return fmt.Sprintf("BIT(%d)", size)
	}
	return fmt.Sprintf("BIT VARYING(%d)", size)
}

func (p *pgDialect) DATETIME() string {
	return "TIMESTAMP"
}

func (p *pgDialect) CurDate() string {
	return "CURRENT_DATE"
}

func (p *pgDialect) CurTime() string {
	return "CURRENT_TIME"
}

func (p *pgDialect) Now() string {
	return "CURRENT_TIMESTAMP"
}

func (p *pgDialect) Uuid() string {
	return "UUID_GENERATE_V4()"
}

func (p *pgDialect) Boolean(val bool) string {
	if true {
		return "TRUE"
	}
	return "FALSE"
}

func (p *pgDialect) CreateDatabase() string {
	return fmt.Sprintf("CREATE DATABASE %s ENCODING=UTF8", p.config.Name)
}

func (p *pgDialect) DropDatabase() string {
	return fmt.Sprintf("DROP DATABASE %s", p.config.Name)
}

func (p *pgDialect) Resource() string {
	return fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s?sslmode=%s",
		p.config.Driver, p.config.User, p.config.Password, p.config.Host, p.config.Port, p.config.Name, p.config.Ssl)

}

func (p *pgDialect) Shell() (string, []string) {
	return "psql", []string{
		"-h", p.config.Host,
		"-p", strconv.Itoa(p.config.Port),
		"-d", p.config.Name,
		"-U", p.config.User}
}

func (p *pgDialect) Setup() string {
	var buf bytes.Buffer
	for _, p := range []string{"uuid-ossp", "pgcrypto"} {
		fmt.Fprintf(&buf, "CREATE EXTENSION IF NOT EXISTS \"%s\";", p)
	}
	return buf.String()
}

func (p *pgDialect) String() string {
	return fmt.Sprintf("%s@%s:%d/%s", p.config.User, p.config.Host, p.config.Port, p.config.Name)
}

func (p *pgDialect) Driver() string {
	return p.config.Driver
}

func (p *pgDialect) Select(table string, columns []string, where, order string, offset, limit int) string {
	if where != "" {
		where = fmt.Sprintf(" WHERE %s", where)
	}
	if order != "" {
		order = fmt.Sprintf(" ORDER BY %s", order)
	}

	ofs := ""
	if offset > 0 {
		ofs = fmt.Sprintf(" OFFSET %d", offset)
	}

	lis := ""
	if limit > 0 {
		ofs = fmt.Sprintf(" LIMIT %d", limit)
	}

	return fmt.Sprintf("SELECT %s FROM %s %s%s%s%s", strings.Join(columns, ", "), table, where, order, ofs, lis)
}
