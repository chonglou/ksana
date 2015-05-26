package ksana

type Dialect interface {
	SERIAL() string
	UUID() string
	BOOLEAN() string
	FLOAT() string
	DOUBLE() string
	BLOB() string
	BYTES(fix bool, size int) string
	DATETIME() string

	CurDate() string
	CurTime() string
	Now() string
	Uuid() string
	Boolean(val bool) string

	CreateDatabase(name string) string
	DropDatabase(name string) string

	Resource() string
	Shell() (string, []string)
	Setup() string
	String() string

	Select(table string, columns []string, where, order string, offset, limit int) string
}
