package ksana

type Bean interface{}

type Engine interface {
	Router(path string, router Router)
	Bean(Bean)
}
