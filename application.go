package ksana

import (
	"container/list"
	"errors"
	"flag"
	"fmt"
	//"io/ioutil"
	"log"
	"net/http"
	"os"
	// "os/signal"
	//"time"
)

const VERSION = "v20150510"

type Application interface {
	Start() error
	Router() Router
}

func New() (Application, error) {
	cfg := flag.String("config", "context.xml", "configuration filename")
	act := flag.String("action", "server", "server | migrate | routes")
	flag.Parse()

	var err error
	var app Application

	for _, a := range []string{"server", "migrate", "routes"} {
		if a == *act {
			ctx := Context{}
			err = ctx.load(*cfg)
			if err == nil {
				app = &application{
					config: *cfg,
					action: *act,

					ctx: &ctx,
					router: &router{
						routes: list.New(),
						ctx:    &ctx},
					migrations: make(map[string]string, 0)}
			}
			break
		}
	}
	if app == nil {
		err = errors.New(
			fmt.Sprintf("Unknown action, please use `%s -h` for more options.",
				os.Args[0]))
	}

	return app, err
}

type application struct {
	config     string
	action     string
	ctx        *Context
	router     Router
	migrations map[string]string
}

func (app *application) Router() Router {
	return app.router
}

func (app *application) Start() error {

	var err error

	switch app.action {
	case "server":
		err = app.server()
	case "migrate":
		err = app.migrate()
	default:
	}

	return err
}

func (app *application) migrate() error {
	return nil
}

func (app *application) server() error {
	log.Printf("=> Booting Ksana(%s)", VERSION)
	log.Printf(
		"=> Application starting in %s on http://0.0.0.0:%v",
		app.ctx.Mode,
		app.ctx.Port)
	log.Printf("=> Run `cat %s` for more startup options", app.config)
	log.Println("=> Ctrl-C to shutdown server")

	return http.ListenAndServe(fmt.Sprintf(":%d", app.ctx.Port), app.router)
}
