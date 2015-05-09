package ksana

import (
	"database/sql"
	"encoding/xml"
	"github.com/fzzy/radix/extra/pool"
	"github.com/fzzy/radix/redis"
	"io/ioutil"
	"log"
	"log/syslog"
	"os"
)

type Context struct {
	Db     *sql.DB
	pool   *pool.Pool
	Logger *syslog.Writer
}

type property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type bean struct {
	Id         string     `xml:"id,attr"`
	Inner      bool       `xml:"inner,attr"`
	Properties []property `xml:"property"`
}

type configuration struct {
	XMLName xml.Name `xml:"ksana"`

	Name string `xml:"name,attr"`
	Mode string `xml:"mode,attr"`

	Port int    `xml:"port,attr"`

	Beans []bean `xml:"bean"`
}

func (c *Context) Init() {
	const fn = "context.xml"

	xf, err := os.Open(fn)
	if err != nil {
		log.Fatalf("Error on open %s: %v", fn, err)
	}
	defer xf.Close()

	data, _ := ioutil.ReadAll(xf)

	cfg := configuration{}
	err = xml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("Error on parse %s: %v", fn, err)
	}

	log.Printf("=> Booting Ksana %s", VERSION)
	log.Printf("=> Application starting in %s on http://0.0.0.0:%v\n", cfg.Mode, cfg.Port)
	log.Println("=> Run `cat context.xml` for more startup options")
	log.Println("=> Ctrl-C to shutdown server")

	c.openLogger(cfg.Name)

}

func (c *Context) openLogger(tag string) {
	var level syslog.Priority
	if os.Getenv("KSANA_ENVIRONMENT") == "production" {
		level = syslog.LOG_INFO
	} else {
		level = syslog.LOG_DEBUG
	}
	logger, err := syslog.New(level, tag)
	if err != nil {
		log.Fatalf("Error on init logger: %v", err)
	}
  logger.Info("Start...")
	c.Logger = logger
}

func (c *Context) openDatabase(adapter string, url string) {
	db, err := sql.Open(adapter, url)
	if err != nil {
		log.Fatalf("Error on open database connect: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error on database ping: %v", err)
	}
	c.Db = db
}

func (c *Context) openRedis(url string, size int) {
	p, e := pool.NewPool("tcp", url, size)
	if e != nil {
		log.Fatalf("Error on open redis connection pool: %v", e)
	}

	var cl *redis.Client
	cl, e = p.Get()
	if e != nil {
		log.Fatalf("Error on open redis connection pool: %v", e)
	}
	e = cl.Cmd("PING").Err
	if e != nil {
		log.Fatalf("Error on open redis ping: %v", e)
	}
	p.Put(cl)

	c.pool = p
}

type RedisFunc func(*redis.Client) (interface{}, error)

func (c *Context) Redis(f RedisFunc) (interface{}, error) {
	cl, e := c.pool.Get()
	defer c.pool.Put(cl)

	if e != nil {
		c.Logger.Err("Error on get redis connection: %v" + e.Error())
	}
	return f(cl)
}
