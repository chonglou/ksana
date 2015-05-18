package main

import (
	"github.com/chonglou/ksana"
	//"github.com/chonglou/ksana/auth"
	_ "github.com/lib/pq"
	"log"
	"time"
)

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

	//------------------Engine-----------------------------------
	// ae := auth.AuthEngine{}
	// app.Mount("/auth", &ae)

	//-------------------SERVER----------------------------

	if err = app.Start(); err != nil {
		log.Fatalf(err.Error())
	}
}
