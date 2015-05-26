package ksana

import (
	orm "github.com/chonglou/ksana/orm"
)

type Engine interface {
	Router(path string, router Router)
	Migration(orm.Connection)
}
