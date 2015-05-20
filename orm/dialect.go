package ksana_orm

type Dialect interface {
	Serial(name string)
	Uuid(name string)
	Boolean(name string, def bool)
	Float(name string, def float32)
	Double(name string, def float64)
	Created(),
	Updated(),
	Blob(name string, null bool, def string)


	CreateDatabase(name string) string
	DropDatabase(name string) string
	Shell(cfg *Config) string
}
