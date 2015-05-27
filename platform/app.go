package main

import (
	"github.com/chonglou/ksana"
	"github.com/chonglou/ksana/auth"
	_ "github.com/lib/pq"
	"log"
)

func main() {

	app, err := ksana.New()
	if err != nil {
		log.Fatalf(err.Error())
	}

	engines := make(map[string]ksana.Engine, 0)
	engines["/auth"] = &kuth.AuthEngine{}
	for k, v := range engines {
		app.Mount(k, v)
	}

	if err = app.Start(); err != nil {
		log.Fatalf(err.Error())
	}
}
