package ksana

type Engine interface {
	Router(path string, router Router)
	Migration(Migrator, *Sql)
}
