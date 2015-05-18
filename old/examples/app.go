package main

import (
	"github.com/chonglou/ksana"
	"github.com/chonglou/ksana/auth"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func sayHello(wrt http.ResponseWriter) error {
	wrt.Write([]byte("Hello,"))
	return nil
}

func main() {

	app, err := ksana.New()
	if err != nil {
		log.Fatalf(err.Error())
	}

	//---------------HTTP----------------------------------------------
	router := app.Router()
	router.Get("/hello$", sayHello, func(wrt http.ResponseWriter) {
		wrt.Write([]byte(" Ksana(HTTP GET)!"))
	})
	router.Any("/test$", func(wrt http.ResponseWriter) {
		wrt.Write([]byte("Hello, Ksans(HTTP ANY)"))
	})

	fns := []ksana.Handler{
		sayHello,
		func(wrt http.ResponseWriter) {
			wrt.Write([]byte("Ksana"))
		},
		func(wrt http.ResponseWriter) {
			wrt.Write([]byte("(HTTP RESOURCES)"))
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

	mig := app.Migration()

	mig.Add("201505151051", "CREATE TABLE T1(f1 INT)", "DROP TABLE T1")
	mig.Add("201505151052", "CREATE TABLE T2(f1 INT)", "DROP TABLE T2")
	mig.Add("201505151053", "CREATE TABLE T3(f1 INT)", "DROP TABLE T3")
	mig.Add("201505151054", "CREATE TABLE T4(f1 INT)", "DROP TABLE T4")
	mig.Add("201505151055", "CREATE TABLE T5(f1 INT)", "DROP TABLE T5")

	//------------------Engine-----------------------------------
	ae := auth.AuthEngine{}
	app.Mount("/auth", &ae)

	//-------------------SERVER----------------------------

	if err = app.Start(); err != nil {
		log.Fatalf(err.Error())
	}
}
