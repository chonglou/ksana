package ksana

import (
	"fmt"
)

type RedisCacheManager struct {
	redis *Redis
}

func (p *RedisCacheManager) key(key string) string {
	return fmt.Sprintf("cache://%s", key)
}

func (p *RedisCacheManager) Set(key string, value interface{}, expire int64) error {
	return p.redis.Set(p.key(key), value, expire)
}

func (p *RedisCacheManager) Get(key string, value interface{}) error {
	return p.redis.Get(p.key(key), value)
}
