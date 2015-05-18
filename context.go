package ksana

import (
	"database/sql"
	"github.com/fzzy/radix/extra/pool"
	"github.com/fzzy/radix/redis"
	"log/syslog"
)

type Context struct {
	config *configuration
	beans  map[string]interface{}
}

func (c *Context) Put(key string, val interface{}) {
	c.beans[key] = val
}

func (c *Context) Get(key string) interface{} {
	return c.beans[key]
}

func (c *Context) openDatabase(driver string, url string) error {
	logger.Info("Connect to database")
	db, err := sql.Open(driver, url)
	if err != nil {
		return err
	}
	logger.Info("Ping database")
	err = db.Ping()
	if err != nil {
		return err
	}
	c.beans["db"] = db
	logger.Info("Database setup successfull")
	return nil
}

func (c *Context) openRedis(url string, size int, db int) error {
	logger.Info("Connect to redis")
	p, e := pool.NewPool("tcp", url, size)
	if e != nil {
		return e
	}

	var cl *redis.Client
	cl, e = p.Get()
	if e != nil {
		return e
	}
	logger.Info("Ping redis")
	e = cl.Cmd("PING").Err
	if e != nil {
		return e
	}
	p.Put(cl)

	c.beans["redis"] = p
	logger.Info("Redis setup successfull")
	return nil
}

func (c *Context) Load(file string) error {
	logger.Info("Booting Ksana(" + VERSION + ")")

	err := readConfig(c.config, file)
	if err != nil {
		return err
	}
	var log *syslog.Writer
	log, err = openLogger(c.config.Env, c.config.Name)
	if err != nil {
		return err
	}
	c.beans["logger"] = log

	if err = c.openDatabase(
		c.config.Database.Driver,
		c.config.Database.Url()); err != nil {
		return err
	}
	if err = c.openRedis(
		c.config.Redis.Url(),
		c.config.Redis.Pool,
		c.config.Redis.Db); err != nil {
		return err
	}

	return nil
}

var ctx = Context{config: &configuration{},
	beans: make(map[string]interface{}, 0)}
