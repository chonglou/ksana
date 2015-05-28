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
	actions := []string{"server", "migrate", "rollback", "routes", "db", "redis", "nginx"}
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
			for _, b := range []interface{}{db, sq, &redis} {
				Map(b)
			}

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
	case "nginx":
		app.nginx()
	default:
	}

	return err
}

func (app *application) nginx() {
	hn, _ := os.Hostname()
	wd, _ := os.Getwd()

	var buf bytes.Buffer
	fmt.Fprintf(
		&buf,
		`
upstream ksana.conf {
	server http://localhost:%d fail_timeout=0;
}
`,
		app.config.Web.Port)
	fmt.Fprintf(
		&buf,
		`
server {
	listen 443;
	ssl  on;
	ssl_certificate  ssl/ksana.crt;
	ssl_certificate_key  ssl/ksana.key;
	ssl_session_timeout  5m;
	ssl_protocols  SSLv2 SSLv3 TLSv1;
	ssl_ciphers  RC4:HIGH:!aNULL:!MD5;
	ssl_prefer_server_ciphers  on;

	client_max_body_size 4G;
	keepalive_timeout 10;

	server_name %s;

	root %s/public;
	try_files $uri/index.html $uri @ksana.conf;

	location @ksana.conf {
		proxy_set_header X-Forwarded-Proto https;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
		proxy_set_header Host $http_host;
		proxy_set_header  X-Real-IP $remote_addr;
		proxy_redirect off;
		proxy_pass http://ksana.conf;
		# limit_req zone=one;
		access_log log/ksana.access.log;
		error_log log/ksana.error.log;
	}

	location ^~ /assets/ {
		gzip_static on;
		expires max;
		add_header Cache-Control public;
	}

	location = /50x.html {
		root html;
	}

	location = /404.html {
		root html;
	}

	location @503 {
		error_page 405 = /system/maintenance.html;
		if (-f $document_root/system/maintenance.html) {
			rewrite ^(.*)$ /system/maintenance.html break;
		}
		rewrite ^(.*)$ /503.html break;
	}

	if ($request_method !~ ^(GET|HEAD|PUT|PATCH|POST|DELETE|OPTIONS)$ ){
		return 405;
	}

	if (-f $document_root/system/maintenance.html) {
		return 503;
	}

	location ~ \.(php|jsp|asp)$ {
		return 405;
	}

}
		`, hn, wd)
	fmt.Fprintf(&buf, "\n")
	buf.WriteTo(os.Stdout)
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
