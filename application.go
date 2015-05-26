package ksana

import (
	"bytes"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

const VERSION = "v20150510"

type Application interface {
	Start() error
	Router() Router
	Model() Model
	Migrator() Migrator
	Mount(path string, engine Engine)
}

func New() (Application, error) {
	actions := []string{"server", "migrate", "rollback", "routes", "db", "redis"}
	cfg := flag.String("c", "config.json", "configuration file name")
	act := flag.String("r", "server", "running: "+strings.Join(actions, " | "))
	flag.Parse()

	var err error
	var app Application
	config := configuration{}

	err = readConfig(&config, *cfg)
	if err != nil {
		return nil, err
	}

	for _, a := range actions {

		db, sq, err := OpenDB(&config.Database)
		if err != nil {
			return nil, err
		}
		var mig Migrator
		mig, err = NewMigrator("migrate", db, sq)
		if err != nil {
			return nil, err
		}

		redis := Redis{}
		err = redis.Open(&config.Redis)
		if err != nil {
			return nil, err
		}

		var rtr Router
		rtr, err = NewRouter("views")
		if err != nil {
			return nil, err
		}

		config.file = *cfg
		if a == *act {
			app = &application{
				config:   &config,
				action:   *act,
				router:   rtr,
				model:    &model{sql: sq},
				migrator: mig,
				redis:    &redis,
				db:       db,
				sql:      sq,
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
	config *configuration
	action string

	router   Router
	model    Model
	migrator Migrator
	redis    *Redis
	db       *sql.DB
	sql      *Sql
}

func (app *application) Mount(mount string, e Engine) {
	e.Router(mount, app.Router())
	e.Migration(app.Migrator(), app.Sql())
}

func (app *application) Model() Model {
	return app.model
}

func (app *application) Router() Router {
	return app.router
}

func (app *application) Sql() *Sql {
	return app.sql
}

func (app *application) Db() *sql.DB {
	return app.db
}

func (app *application) Migrator() Migrator {
	return app.migrator
}

func (app *application) Start() error {
	var err error

	switch app.action {
	case "server":
		err = app.server()
	case "migrate":
		err = app.migrator.Migrate()
	case "rollback":
		err = app.migrator.Rollback()
	case "routes":
		app.routes()
	case "db":
		cmd, args := app.sql.Shell()
		err = Shell(cmd, args...)
	case "redis":
		cmd, args := app.redis.Shell()
		err = Shell(cmd, args...)
	default:
	}

	return err
}

func (app *application) shell(cmd string, args ...string) error {
	bin, err := exec.LookPath(cmd)
	if err != nil {
		return err
	}
	return syscall.Exec(bin, append([]string{cmd}, args...), os.Environ())
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
		app.config.Env,
		app.config.Web.Port)
	log.Printf("=> Run `cat %s` for more startup options", app.config.file)
	log.Println("=> Ctrl-C to shutdown server")

	return http.ListenAndServe(fmt.Sprintf(":%d", app.config.Web.Port), app.router)
}
