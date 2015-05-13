package ksana

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	// "os/signal"
	"time"
)

const VERSION = "v20150510"

type Application struct {
	ctx *Context
}

func (app *Application) Start() {

	cf := flag.String("config", "context.xml", "configuration filename")
	server := flag.Bool("server", false, "runing server")
	generate := flag.String("generate", "", "config | migration | model | controller | test")
	name := flag.String("name", "", "name for generate")

	//act := flag.String("run", "server", "generate | migrate | config | server | database | redis | worker")

	flag.Parse()

	var err error
	switch {
	case *server:
		err = app.server(*cf)
	case *generate != "":
		err = app.generate(*generate, *name)
	default:
		err = errors.New(fmt.Sprintf("Unknown action, please use `%s -h` for more options.", os.Args[0]))
	}
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (app *Application) generate(mode, name string) error {
	const path = "db/migrations"

	var err error
	log.Printf("Generating %s for %s", name, mode)
	switch mode {
	case "migration":
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}

		fn := fmt.Sprintf(
			"%s/%s_%s.sql",
			path,
			time.Now().Format("2006010215040507"),
			name)
		log.Printf("Create file %s", fn)
		err = ioutil.WriteFile(fn, []byte("/*\n* File: "+fn+"\n*/\n"), 0644)
		if err != nil {
			return err
		}

		if _, err1 := os.Stat("db/seeds.sql"); os.IsNotExist(err1) {
			fn = "db/seeds.sql"
			log.Printf("Create file %s", fn)
			err = ioutil.WriteFile(fn, []byte("/*\n* File: "+fn+"\n*/\n"), 0644)
		}

	default:
		return errors.New("Unknown gererate type " + mode)
	}
	log.Println("Done!!!")
	return nil
}

func (app *Application) config(file string) error {
	ctx := &Context{}
	err := ctx.load(file)
	if err != nil {
		return err
	}
	app.ctx = ctx
	return nil
}

func (app *Application) server(cfg string) error {
	err := app.config(cfg)
	if err != nil {
		return err
	}

	log.Printf("=> Booting Ksana(%s)", VERSION)
	log.Printf(
		"=> Application starting in %s on http://0.0.0.0:%v",
		app.ctx.Mode,
		app.ctx.Port)
	log.Println("=> Run `cat context.xml` for more startup options")
	log.Println("=> Ctrl-C to shutdown server")

	return http.ListenAndServe(fmt.Sprintf(":%d", app.ctx.Port), nil)
}
