package ksana

import (
	"fmt"
)

type RedisCacheManager struct {
	redis *Redis
}

func (rcm *RedisCacheManager) key(key string) string {
	return fmt.Sprintf("cache://%s", key)
}

func (rcm *RedisCacheManager) Set(key string, value interface{}, expire int64) error {
	return rcm.redis.Set(rcm.key(key), value, expire)
}

func (rcm *RedisCacheManager) Get(key string, value interface{}) error {
	return rcm.redis.Get(rcm.key(key), value)
}

func (rcm *RedisCacheManager) Gc() {
	logger.Info("Call gc")
}
