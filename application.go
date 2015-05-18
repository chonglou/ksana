package ksana

import (
	"bytes"
	"container/list"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const VERSION = "v20150510"

type Application interface {
	Start() error
	Router() Router
	Migration() Migration
	Mount(path string, engine Engine)
}

func New() (Application, error) {
	actions := []string{"server", "migrate", "rollback", "routes"}
	cfg := flag.String("c", "context.xml", "configuration file name")
	act := flag.String("r", "server", "running: "+strings.Join(actions, " | "))
	flag.Parse()

	var err error
	var app Application

	err = ctx.Load(*cfg)
	if err != nil {
		return nil, err
	}

	for _, a := range actions {
		if a == *act {
			app = &application{
				config:    *cfg,
				action:    *act,
				router:    &router{routes: list.New()},
				migration: &migration{},
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
	config string
	action string

	router    Router
	migration Migration
}

func (app *application) Mount(p string, e Engine) {
	e.Router(p, app.Router())
}

func (app *application) Migration() Migration {
	return app.migration
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
		err = app.migration.Migrate()
	case "rollback":
		err = app.migration.Rollback()
	case "routes":
		app.routes()
	default:
	}

	return err
}

func (app *application) routes() {
	var buf bytes.Buffer
	app.router.Status(&buf)
	buf.WriteTo(os.Stdout)
}

func (app *application) server() error {
	log.Printf("=> Booting Ksana(%s)", VERSION)
	log.Printf(
		"=> Application starting in %s on http://0.0.0.0:%v",
		ctx.Config.Env,
		ctx.Config.Port)
	log.Printf("=> Run `cat %s` for more startup options", app.config)
	log.Println("=> Ctrl-C to shutdown server")

	return http.ListenAndServe(fmt.Sprintf(":%d", ctx.Config.Port), app.router)
}
