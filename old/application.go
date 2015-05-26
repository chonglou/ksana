package ksana

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	orm "github.com/chonglou/ksana/orm"
	web "github.com/chonglou/ksana/web"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

const VERSION = "v20150510"

var logger, _ = OpenLogger("ksana-app")

type Application interface {
	Start() error
	Router() web.Router
	Model() orm.Model
	Mount(path string, engine web.Engine)
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

		model, err := orm.New("db/migrate", &config.Database)
		if err != nil {
			return nil, err
		}

		redis := Redis{}
		err = redis.Open(&config.Redis)
		if err != nil {
			return nil, err
		}

		config.file = *cfg
		if a == *act {
			app = &application{
				config: &config,
				action: *act,
				router: web.New("app/views"),
				model:  model,
				redis:  &redis,
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

	router web.Router
	model  orm.Model
	redis  *redis.Connection
}

func (app *application) Mount(mount string, e web.Engine) {
	e.Router(mount, app.Router())
}

func (app *application) Model() orm.Model {
	return app.model
}

func (app *application) Router() web.Router {
	return app.router
}

func (app *application) Start() error {
	var err error

	switch app.action {
	case "server":
		err = app.server()
	case "migrate":
		err = app.model.Db().Migrate()
	case "rollback":
		err = app.model.Db().Rollback()
	case "routes":
		app.routes()
	case "db":
		err = app.model.Db().Shell()
	case "redis":
		err = app.redis.Shell()
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
