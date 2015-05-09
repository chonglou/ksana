package ksana

import (
	"database/sql"
	"encoding/xml"
  "io/ioutil"
	"github.com/fzzy/radix/extra/pool"
	"github.com/fzzy/radix/redis"
	"log"
  "os"
	"log/syslog"
)

type Context struct {
	Db     *sql.DB
	pool   *pool.Pool
  Logger *syslog.Writer
}

type property struct {
	name  string `xml:"name,attr"`
	value string `xml:"value,attr"`
}

type bean struct {
	id    string     `xml:"id,attr"`
	inner bool       `xml:"inner,attr"`
	items []property `xml:"property"`
}

type configuration struct {
	XMLName xml.Name `xml:"ksana"`

	name string `xml:"name,attr"`
	mode string `xml:"mode,attr"`
	port int    `xml:"port,attr"`

	items []bean `xml:"bean"`
}

func (c *Context) Init() {
	xf, err := os.Open("context.xml")
	if err != nil {
		log.Fatalf("Error on open context.xml: %v", err)
	}
	defer xf.Close()

  data, _ := ioutil.ReadAll(xf)

	cfg := configuration{}
	xml.Unmarshal(data, &cfg)

	log.Printf("=> Booting Ksana %s", VERSION)
	log.Printf("=> Application starting in %s on http://0.0.0.0:%v\n", cfg.mode, cfg.port)
	log.Println("=> Run `gails -h server` for more startup options")
	log.Println("=> Ctrl-C to shutdown server")

	c.openLogger(cfg.name)

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
		c.Logger.Err("Error on get redis connection: %v"+e.Error())
	}
	return f(cl)
}
