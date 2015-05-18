package ksana

type Engine interface {
	Router(path string, router Router)
	Bean(Bean)
}
