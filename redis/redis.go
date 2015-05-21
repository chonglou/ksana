package ksana_redis

import (
	"bytes"
	"encoding/gob"
	"fmt"
	utils "github.com/chonglou/ksana/utils"
	"github.com/fzzy/radix/extra/pool"
	"github.com/fzzy/radix/redis"
	"strconv"
)

var logger, _ = utils.OpenLogger("ksana-redis")

type Config struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	Db   int    `json:"db"`
	Pool int    `json:"pool"`
}

type Connection struct {
	config *Config
	pool   *pool.Pool
}

func (r *Connection) Open(cfg *Config) error {
	logger.Info(fmt.Sprintf("Connect to redis %s:%d/%d", cfg.Host, cfg.Port, cfg.Db))

	df := func(network, addr string) (*redis.Client, error) {
		client, err := redis.Dial(network, addr)
		if err != nil {
			return nil, err
		}
		err = client.Cmd("PING").Err
		if err != nil {
			return nil, err
		}
		err = client.Cmd("SELECT", cfg.Db).Err
		if err != nil {
			return nil, err
		}
		return client, nil
	}

	p, e := pool.NewCustomPool(
		"tcp",
		fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		cfg.Pool, df)
	if e != nil {
		return e
	}
	logger.Info("Connection setup successfull")

	r.pool = p
	r.config = cfg
	return nil

}

func (r *Connection) Shell() error {
	return utils.Shell("telnet", []string{r.config.Host, strconv.Itoa(r.config.Port)}...)
}

func (r *Connection) cmd(f func(*redis.Client) error) error {
	c, e := r.pool.Get()
	if e != nil {
		return e
	}
	defer r.pool.Put(c)
	return f(c)
}

func (r *Connection) Set(key string, val interface{}, expire int64) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(val)
	if err != nil {
		return err
	}

	return r.cmd(func(c *redis.Client) error {
		if expire > 0 {
			return c.Cmd("SET", key, buf.Bytes(), "EX", expire).Err
		}
		return c.Cmd("SET", key, buf.Bytes()).Err

	})

}

func (r *Connection) Del(key string) error {
	return r.cmd(func(c *redis.Client) error {
		return c.Cmd("DEL", key).Err
	})
}

func (r *Connection) Get(key string, val interface{}) error {
	var buf bytes.Buffer
	enc := gob.NewDecoder(&buf)
	err := r.cmd(func(c *redis.Client) error {
		s, e := c.Cmd("get", key).Bytes()
		if e != nil {
			return e
		}
		buf.Write(s)
		return nil
	})

	if err != nil {
		return err
	}

	return enc.Decode(val)
}

func (r *Connection) Expire(key string, time int64) error {
	return r.cmd(func(c *redis.Client) error {
		return c.Cmd("EXPIRE", key, time).Err
	})
}

// func (r *Connection) Cache(key string, val interface{}, f func(interface{}) error, expire int64) error {
// 	err := r.Get(key, val)
// 	if err == nil {
// 		return nil
// 	}
//
// 	err = f(val)
// 	if err != nil {
// 		return err
// 	}
//
// 	go func() {
// 		err := r.Set(key, val)
// 		if err == nil {
// 			r.cmd(func(c *redis.Client) error {
// 				return c.Cmd("expire", key, expire).Err
// 			})
// 		}
//
// 	}()
// 	return nil
// }
