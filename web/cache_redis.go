package ksana_web

import (
	"fmt"
	"github.com/chonglou/ksana/redis"
)

type RedisCacheProvider struct {
	redis *ksana_redis.Connection
}

func (rcp *RedisCacheProvider) key(key string) string {
	return fmt.Sprintf("cache://%s", key)
}

func (rcp *RedisCacheProvider) Set(key string, value interface{}, expire int64) error {
	return rcp.redis.Set(rcp.key(key), value, expire)
}

func (rcp *RedisCacheProvider) Get(key string, value interface{}) error {
	return rcp.redis.Get(rcp.key(key), value)
}

func (fcm *RedisCacheProvider) Gc() {
	logger.Info("Call gc")
}
