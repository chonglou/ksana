package ksana

import (
	"database/sql"
	"github.com/fzzy/radix/extra/pool"
	"github.com/fzzy/radix/redis"
	"log"
	"log/syslog"
	"os"
)

var beans = make(map[string]interface{})

func Register(name string, bean interface{}) {
	if bean == nil {
		log.Fatalf("Register bean is nil")
	}
	if _, dup := beans[name]; dup {
		log.Fatalf("Register called twice for bean " + name)
	}
	beans[name] = bean
}

func Load(name string) {
	xf, err := os.Open(file)
	if err != nil {
		log.Fatalf("Error on open config file: %v", err)
	}
	defer xf.Close()

	var data []byte
	data, err = ioutil.ReadAll(xf)
	if err != nil {
		log.Fatalf("Error on read config file: %v", err)
	}

	err = xml.Unmarshal(data, cfg)
	if err != nil {
		log.Fatalf("Error on parse config file: %v", err)
	}
	return nil
}

//------------------------------------------------------------------------------

type Context struct {
	Redis  *pool.Pool
	Db     *sql.DB
	Logger *syslog.Writer
	Sm     *SessionManager
	Cm     *CacheManager

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

	c.Logger.Info("============ Ksana starting(" + VERSION + ") ============")

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
		case "session":
			pr := b.getString("provider", "file")
			pt := b.getString("path", "tmp/sessions")
			switch pr {
			case "file":
				c.Sm = &SessionManager{provider: &FileSessionProvider{path: pt}}
				c.Logger.Info("Set session to file path: " + pt)
				err = os.MkdirAll(pt, 0700)
			default:
				c.Logger.Warning("Unknown session provider " + pr)
			}
		case "cache":
			pr := b.getString("provider", "file")
			pt := b.getString("path", "tmp/cache")
			switch pr {
			case "file":
				c.Cm = &CacheManager{provider: &FileCacheProvider{path: pt}}
				c.Logger.Info("Set cache to file path: " + pt)
				err = os.MkdirAll(pt, 0700)
			default:
				c.Logger.Warning("Unknown cache provider " + pr)
			}

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