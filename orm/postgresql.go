package ksana_orm

import (
	"bytes"
	"fmt"
	"strconv"
)

type pgDialect struct {
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

func (d *pgDialect) Resource(cfg *Config) string {
	return fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Driver, cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.Ssl)

}

func (d *pgDialect) Shell(cfg *Config) (string, []string) {
	return "psql", []string{
		"-h", cfg.Host,
		"-p", strconv.Itoa(cfg.Port),
		"-d", cfg.Name,
		"-U", cfg.User}
}

func (d *pgDialect) Setup() string {
	var buf bytes.Buffer
	for _, p := range []string{"uuid-ossp", "pgcrypto"} {
		fmt.Fprintf(&buf, "CREATE EXTENSION IF NOT EXISTS \"%s\";", p)
	}
	return buf.String()
}

func (d *pgDialect) String(cfg *Config) string {
	return fmt.Sprintf("%s@%s:%d/%s", cfg.User, cfg.Host, cfg.Port, cfg.Name)
}
