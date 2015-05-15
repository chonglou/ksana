package main

import (
	"github.com/chonglou/ksana"
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
		"users",
		ksana.Controller{
			Index:   fns,
			Show:    fns,
			New:     fns,
			Create:  fns,
			Edit:    fns,
			Update:  fns,
			Destroy: fns})

	if err = app.Start(); err != nil {
		log.Fatalf(err.Error())
	}
}
