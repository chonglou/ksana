package ksana

import (
	"database/sql"
	"github.com/fzzy/radix/extra/pool"
	"github.com/fzzy/radix/redis"
	"log/syslog"
)

type Context struct {
	Config *configuration
	Beans  map[string]interface{}
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
	c.Beans["db"] = db
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

	c.Beans["redis"] = p
	logger.Info("Redis setup successfull")
	return nil
}

func (c *Context) Load(file string) error {
	err := readConfig(c.Config, file)
	if err != nil {
		return err
	}
	var log *syslog.Writer
	log, err = openLogger(c.Config.Env, c.Config.Name)
	if err != nil {
		return err
	}
	c.Beans["logger"] = log

	if err = c.openDatabase(
		c.Config.Database.Driver,
		c.Config.Database.Url); err != nil {
		return err
	}
	if err = c.openRedis(
		c.Config.Redis.Url,
		c.Config.Redis.Pool,
		c.Config.Redis.Db); err != nil {
		return err
	}

	return nil
}

var ctx = Context{Config: &configuration{},
	Beans: make(map[string]interface{}, 0)}
