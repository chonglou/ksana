package ksana

import (
	"github.com/fzzy/radix/extra/pool"
	//	"github.com/fzzy/radix/redis"
	"log"
)

type Redis struct {
	pool *pool.Pool
}

func (r *Redis) Open(url string, size int) {
	p, e := pool.NewPool("tcp", url, size)
	if e != nil {
		log.Fatalf("Error on open redis connection pool: %v", e)
	}
	r.pool = p
}
func (r *Redis) Set(key string, val string) {
	c, e := r.pool.Get()
	defer r.pool.Put(c)

	if e != nil {
		log.Fatalf("Error on get redis connection: %v", e)
	}
	e = c.Cmd("set", key, val).Err
	if e != nil {
		log.Fatalf("Error redis set: %v", e)
	}
}
func (r *Redis) GetString(key string) string {
	c, e := r.pool.Get()
	defer r.pool.Put(c)

	if e != nil {
		log.Fatalf("Error on get redis connection:%v", e)
	}
	s, e1 := c.Cmd("get", key).Str()
	if e1 != nil {
		log.Fatalf("Error on redis get: %v", e1)
	}
	return s
}

func (r *Redis) Ping() {
	c, e := r.pool.Get()
	defer r.pool.Put(c)
	if e != nil {
		log.Fatalf("Error on get redis connection:%v", e)
	}
	e = c.Cmd("PING").Err
	if e != nil {
		log.Printf("Error on redis ping: %v", e)
	}
}
