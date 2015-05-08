package ksana

import (
	"github.com/fzzy/radix/extra/pool"
	"github.com/fzzy/radix/redis"
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

func (r *Redis) Run(command string, args ...interface{}) (*redis.Reply, error) {
	c, e := r.pool.Get()
	defer r.pool.Put(c)

	if e != nil {
		log.Fatalf("Error on get redis connection: %v", e)
	}

	v := c.Cmd(command, args)

	e = v.Err
	if e != nil {
		log.Printf("Error redis set: %v", e)
	}
	return v, e
}
