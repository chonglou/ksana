package ksana

import (
	"database/sql"
	"fmt"
	"github.com/fzzy/radix/extra/pool"
	"github.com/fzzy/radix/redis"
	"log/syslog"
	"os"
  "log"
)

type Context struct {
	Redis  *pool.Pool
	Db     *sql.DB
	Logger *syslog.Writer

	Port int
	Name string
	Mode string
}

func (c *Context) load(fn string) error {

	cfg := configuration{}
	err := loadConfiguration(fn, &cfg)
	if err != nil {
		return err
	}

	err = c.openLogger(cfg.Name)
	if err != nil {
		return err
	}

	c.Logger.Info("=> Booting Ksana " + VERSION)
	c.Logger.Info(
		fmt.Sprintf(
			"=> Application starting in %s on http://0.0.0.0:%v",
			cfg.Mode,
			cfg.Port))
	c.Logger.Info("=> Run `cat context.xml` for more startup options")
	c.Logger.Info("=> Ctrl-C to shutdown server")

	c.Port = cfg.Port
	c.Name = cfg.Name
	c.Mode = cfg.Mode

	for _, b := range cfg.Beans {
		switch b.Name {
		case "database":
			err = c.openDatabase(
				b.getString("adapter", "postgres"),
				b.getString("url", "postgres://postgres@localhost/ksana?sslmode=disable"))
		case "redis":
			err = c.openRedis(
				b.getString("url", "localhost:6379"),
				b.getInt("pool", 12),
				b.getInt("db", 0))
		default:
			c.Logger.Warning("Unknown bean " + b.Name)
		}
		if err != nil {
			return err
		}
	}
	return nil

}

func (c *Context) openLogger(tag string) error {
	var level syslog.Priority
	if os.Getenv("KSANA_ENVIRONMENT") == "production" {
		level = syslog.LOG_INFO
	} else {
		level = syslog.LOG_DEBUG
	}
	logger, err := syslog.New(level, tag)
	if err != nil {
		return err
	}
	logger.Info("Start Ksana...")
	c.Logger = logger
	return nil
}

func (c *Context) openDatabase(adapter string, url string) error {
  c.Logger.Info("Connect to database")
	db, err := sql.Open(adapter, url)
	if err != nil {
		return err
	}
  c.Logger.Info("Ping database")
	err = db.Ping()
	if err != nil {
		return err
	}
	c.Db = db
  c.Logger.Info("Database setup successfull")
	return nil
}

func (c *Context) openRedis(url string, size int, db int) error {
  c.Logger.Info("Connect to redis")
	p, e := pool.NewPool("tcp", url, size)
	if e != nil {
		return e
	}

	var cl *redis.Client
	cl, e = p.Get()
	if e != nil {
		return e
	}
  c.Logger.Info("Ping redis")
	e = cl.Cmd("PING").Err
	if e != nil {
		return e
	}
	p.Put(cl)

	c.Redis = p
	c.Logger.Info("Redis setup successfull")
	return nil
}

//-----------------------------------------------------------------------------
var Ctx *Context

func init() {
	Ctx := &Context{}
	err := Ctx.load("context.xml")
	if err != nil {
		log.Fatalf("Error on load configuration from context.xml: %v", err)
	}
}
