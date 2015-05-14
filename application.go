package ksana

import (
	"errors"
	"flag"
	"fmt"
	//"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	// "os/signal"
	//"time"
)

const VERSION = "v20150510"

type Application struct {
	mux        *RegexpMux
	migrations map[string][]string
}

func (app *Application) Start() {

	cf := flag.String("config", "context.xml", "configuration filename")
	server := flag.Bool("server", false, "runing server")
	migrate := flag.Bool("migrate", false, "migrate database")

	flag.Parse()

	var err error

	switch {
	case *server:
		err = app.server(*cf)
	case *migrate:
		err = app.migrate(*cf)
	default:
		err = errors.New(
			fmt.Sprintf("Unknown action, please use `%s -h` for more options.",
				os.Args[0]))
	}
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (app *Application) config(file string) error {
	ctx := &Context{}
	err := ctx.load(file)
	if err != nil {
		return err
	}
	app.migrations = make(map[string][]string, 0)
	app.mux = &RegexpMux{
		handlers: make(map[*regexp.Regexp]HandlerFunc, 0),
		ctx:      ctx}

	return nil
}

func (app *Application) migrate(file string) error {
	return nil
}

func (app *Application) server(file string) error {
	err := app.config(file)
	if err != nil {
		return err
	}

	log.Printf("=> Booting Ksana(%s)", VERSION)
	log.Printf(
		"=> Application starting in %s on http://0.0.0.0:%v",
		app.mux.ctx.Mode,
		app.mux.ctx.Port)
	log.Printf("=> Run `cat %s` for more startup options", file)
	log.Println("=> Ctrl-C to shutdown server")

	return http.ListenAndServe(fmt.Sprintf(":%d", app.mux.ctx.Port), app.mux)
}
