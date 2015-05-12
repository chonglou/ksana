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
	"strconv"
)

type Context struct {
	Port int
	Name string
	Mode string

	Db     *sql.DB
	Logger *syslog.Writer

	pool *pool.Pool
}

type property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type bean struct {
	Name       string     `xml:"name,attr"`
	Class      string     `xml:"class,attr"`
	Properties []property `xml:"property"`
}

func (b *bean) getString(name string) string {
	for _, b := range b.Properties {
		if b.Name == name {
			return b.Value
		}
	}
	return ""
}

func (b *bean) getInt(name string) int {
	i, err := strconv.Atoi(b.getString(name))
	if err != nil {
		log.Fatal("Bad property %s", name)
	}
	return i
}

type configuration struct {
	XMLName xml.Name `xml:"ksana"`

	Name string `xml:"name,attr"`
	Mode string `xml:"mode,attr"`

	Port int `xml:"port,attr"`

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

	c.Port = cfg.Port
	c.Name = cfg.Name
	c.Mode = cfg.Mode

	c.openLogger(cfg.Name)

	for _, b := range cfg.Beans {
		switch b.Name {
		case "database":
			c.openDatabase(b.getString("adapter"), b.getString("url"))
		case "redis":
			c.openRedis(b.getString("url"), b.getInt("pool"))
		default:
			//todo auto create bean
			c.Logger.Warning("Unknown bean " + b.Name)
		}
	}

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
		log.Fatalf("Error on connect syslog: %v", err)
	}
	logger.Info("Start Ksana...")
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
	c.Logger.Info("Connect to database successfull")
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
	c.Logger.Info("Connect to redis successfull")
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

var glSessions *SessionManager

func init() {
	//todo generate session manager
	//glSessions, _ = newSessionManager("redis", "gsessionid", 3600)
}
