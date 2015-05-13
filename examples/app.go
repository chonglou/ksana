package main

import (
	"github.com/chonglou/ksana"
	_ "github.com/lib/pq"
)

func main() {
	var app = ksana.Application{}
	app.Start()
}
