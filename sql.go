package ksana

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var re_sql_script = regexp.MustCompile("\\$[\\d]+")

type Sql struct {
	dialect Dialect
}

func (p *Sql) Created() string {
	return p.column("created", p.dialect.DATETIME(), false, p.dialect.Now())
}

func (p *Sql) Updated() string {
	return p.column("updated", p.dialect.DATETIME(), true, "")
}

func (p *Sql) Id(uuid bool) string {
	if uuid {
		return fmt.Sprintf(
			"id %s NOT NULL PRIMARY KEY DEFAULT %s",
			p.dialect.UUID(), p.dialect.Uuid())
	}
	return fmt.Sprintf("id %s", p.dialect.SERIAL())
}

func (p *Sql) Bool(name string, def bool) string {
	return p.column(name, p.dialect.BOOLEAN(), false, p.dialect.Boolean(def))
}

func (p *Sql) String(name string, fix bool, size int, big, null bool, def string) string {
	var ts string
	switch {
	case big:
		ts = "TEXT"
	case fix:
		ts = fmt.Sprintf("CHAR(%d)", size)
	default:
		ts = fmt.Sprintf("VARCHAR(%d)", size)
	}
	if def != "" {
		def = fmt.Sprintf("'%s'", def)
	}
	return p.column(name, ts, null, def)
}

func (p *Sql) Int32(name string, null bool, def int) string {
	return p.column(name, "INT", null, fmt.Sprintf("%d", def))
}

func (p *Sql) Int64(name string, null bool, def int64) string {
	return p.column(name, "BIGINT", null, fmt.Sprintf("%d", def))
}

func (p *Sql) Bytes(name string, fix bool, size int, big, null bool) string {
	if big {
		return p.column(name, p.dialect.BLOB(), null, "")
	} else {
		return p.column(name, p.dialect.BYTES(fix, size), null, "")
	}

}

func (p *Sql) Date(name string, null bool, def string) string {
	var ds string
	switch def {
	case "now":
		ds = p.dialect.CurDate()
	default:
		ds = def
	}
	return p.column(name, "DATE", null, ds)
}

func (p *Sql) Time(name string, null bool, def string) string {
	var ds string
	switch def {
	case "now":
		ds = p.dialect.CurTime()
	default:
		ds = def
	}
	return p.column(name, "TIME", null, ds)
}

func (p *Sql) Datetime(name string, null bool, def string) string {
	var ds string
	switch def {
	case "now":
		ds = p.dialect.Now()
	default:
		ds = def
	}
	return p.column(name, p.dialect.DATETIME(), null, ds)
}

func (p *Sql) Float32(name string, def float32) string {
	return p.column(name, p.dialect.FLOAT(), false, fmt.Sprintf("%f", def))
}

func (p *Sql) Float64(name string, def float64) string {
	return p.column(name, p.dialect.DOUBLE(), false, fmt.Sprintf("%f", def))
}

func (p *Sql) column(name string, _type string, null bool, def string) string {
	ns, ds := "", ""
	if !null {
		ns = " NOT NULL"
	}
	if def != "" {
		ds = fmt.Sprintf(" DEFAULT %s", def)
	}
	return fmt.Sprintf("%s %s%s%s", name, _type, ns, ds)
}

func (p *Sql) CreateTable(table string, columns ...string) string {
	return fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s(%s);", table, strings.Join(columns, ", "))
}

func (p *Sql) DropTable(table string) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s;", table)
}

func (p *Sql) CreateIndex(name, table string, unique bool, columns ...string) string {
	idx := "INDEX"
	if unique {
		idx = "UNIQUE INDEX"
	}
	return fmt.Sprintf(
		"CREATE %s %s ON %s (%s);", idx, name, table, strings.Join(columns, ", "))

}

func (p *Sql) DropIndex(name string) string {
	return fmt.Sprintf("DROP INDEX %s;", name)
}

func (p *Sql) Create() string {
	return p.dialect.CreateDatabase()
}

func (p *Sql) Drop() string {
	return p.dialect.DropDatabase()
}

func (p *Sql) Shell() (string, []string) {
	return p.dialect.Shell()
}

//-----------------------curd---------------------------------------------------

func (p *Sql) Insert(table string, columns []string) string {
	vs := make([]string, 0)
	for i := 1; i <= len(columns); i++ {
		vs = append(vs, fmt.Sprintf("$%d", i))
	}
	return p.t(fmt.Sprintf(
		"INSERT INTO %s(%s) VALUES(%s)",
		table,
		strings.Join(columns, ", "),
		strings.Join(vs, ", ")))
}

func (p *Sql) Select(table string, columns []string, where, order string, offset, limit int) string {
	return p.t(p.dialect.Select(table, columns, where, order, offset, limit))
}

func (p *Sql) Delete(table, where string) string {
	if where != "" {
		where = fmt.Sprintf(" WHERE %s", where)
	}
	return p.t(fmt.Sprintf("DELETE FROM %s%s", table, where))
}

func (p *Sql) Update(table, columns, where string) string {
	return p.t(fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, columns, where))
}

func (p *Sql) t(sq string) string {
	switch p.dialect.Driver() {
	case "postgres":
	default:
		sq = re_sql_script.ReplaceAllString(sq, "?")
	}
	logger.Debug(sq)
	return sq
}

//------------------------NEW---------------------------------------------------
func openDB(cfg *databaseConfig) (*sql.DB, *Sql, error) {
	var dlt Dialect
	var err error
	switch cfg.Driver {
	case "postgres":
		dlt = &pgDialect{config: cfg}
	default:
		err = errors.New("Unknown driver: " + cfg.Driver)
	}
	if err != nil {
		return nil, nil, err
	}

	logger.Info("Connect to database " + dlt.String())
	var db *sql.DB
	db, err = sql.Open(cfg.Driver, dlt.Resource())
	if err != nil {
		return nil, nil, err
	}

	logger.Info("Ping database")
	err = db.Ping()
	if err != nil {
		return nil, nil, err
	}

	logger.Info("Run database setup scripts")
	_, err = db.Exec(dlt.Setup())
	if err != nil {
		return nil, nil, err
	}

	return db, &Sql{dialect: dlt}, nil
}
