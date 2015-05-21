package main

import (
	"github.com/chonglou/ksana"
	//"github.com/chonglou/ksana/auth"
	web "github.com/chonglou/ksana/web"
	_ "github.com/lib/pq"
	"log"
	"time"
	//"bytes"
)

type User1 struct {
	Id   int
	Name string
	Nick string
}

type Log1 struct {
	Id      string
	Message string
	Created time.Time
}

func sayHello(req *web.Request, res *web.Response) error {
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
	router.Any("/test$", func(req *web.Request, res *web.Response) error {
		res.Text([]byte("Hello,"))
		return nil
	}, func(req *web.Request, res *web.Response) error {
		res.Text([]byte(" Ksans(HTTP ANY)!!!"))
		return nil
	})

	fns := []web.Handler{
		sayHello,
		func(req *web.Request, res *web.Response) error {
			res.Text([]byte("Ksana"))
			return nil
		},
		func(req *web.Request, res *web.Response) error {
			res.Text([]byte("(HTTP RESOURCES)"))
			return nil
		}}

	router.Resources(
		"/tags",
		web.Controller{
			Index:   fns,
			Show:    fns,
			New:     fns,
			Create:  fns,
			Edit:    fns,
			Update:  fns,
			Destroy: fns})

	//-------------------DATABASE----------------------------------------

	// mod := app.Model()
	// db := mod.Db()
	// db.AddMigration
	//
	// ver := db.Version()
	// var up, down bytes.Buffer
	// for _, b := range []interface{}{User1{}, Log1{}} {
	// 	if !mod.Check(db.path, b){
	// 			ups, downs, err = mod.Table(b)
	// 			if err != nil {
	// 				log.Fatalf("Error on register bean: %v", err)
	// 			}
	// 			up.Write([]bytes(ups))
	// 			down.Write([]bytes(downs))
	// 	}
	//
	// 	up.Write(mod.Index(b, "Name", "Nick"))
	//
	// }
	// db.AddMigration()

	//------------------Engine-----------------------------------
	// ae := auth.AuthEngine{}
	// app.Mount("/auth", &ae)

	//-------------------SERVER----------------------------

	if err = app.Start(); err != nil {
		log.Fatalf(err.Error())
	}
}
