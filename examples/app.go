package main

import (
	"github.com/chonglou/ksana"
	//"github.com/chonglou/ksana/auth"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type User1 struct {
	ksana.Bean
	Id   int    `sql:"type=serial"`
	Name string `sql:"size=255;fix=true"`
}

type Log1 struct {
	ksana.Bean
	Id      string `sql:"type=uuid"`
	Message string
	Created time.Time `sql:"type=created"`
}

func sayHello(req *ksana.Request, res *ksana.Response, ctx *ksana.Context) error {
	val := make(map[string]interface{}, 0)
	val["ok"] = true
	val["created"] = time.Now()
	res.Json(val)
	return nil
}

func main() {

	app, err := ksana.New()
	if err != nil {
		log.Fatalf(err.Error())
	}

	//---------------HTTP----------------------------------------------
	router := app.Router()

	router.Get("/hello$", sayHello)
	router.Any("/test$", func(req *ksana.Request, res *ksana.Response, ctx *ksana.Context) error {
		res.Text([]byte("Hello,"))
		return nil
	}, func(req *ksana.Request, res *ksana.Response, ctx *ksana.Context) error {
		res.Text([]byte(" Ksans(HTTP ANY)!!!"))
		return nil
	})

	fns := []ksana.Handler{
		sayHello,
		func(req *ksana.Request, res *ksana.Response, ctx *ksana.Context) error {
			res.Text([]byte("Ksana"))
			return nil
		},
		func(req *ksana.Request, res *ksana.Response, ctx *ksana.Context) error {
			res.Text([]byte("(HTTP RESOURCES)"))
			return nil
		}}

	router.Resources(
		"/tags",
		ksana.Controller{
			Index:   fns,
			Show:    fns,
			New:     fns,
			Create:  fns,
			Edit:    fns,
			Update:  fns,
			Destroy: fns})

	//-------------------DATABASE----------------------------------------

	mod := app.Model()
	for _, b := range []ksana.Bean{User1{}, Log1{}} {
		err := mod.Register(b)
		if err != nil {
			log.Fatalf("Error on register bean: %v", err)
		}

	}

	//------------------Engine-----------------------------------
	// ae := auth.AuthEngine{}
	// app.Mount("/auth", &ae)

	//-------------------SERVER----------------------------

	if err = app.Start(); err != nil {
		log.Fatalf(err.Error())
	}
}
