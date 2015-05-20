package ksana_web

import (
	"sync"
)

type CacheProvider interface {
	Get(key string, value interface{}) error
	Set(key string, value interface{}, expireTime int64) error
	Gc()
}

type CacheManager struct {
	lock     sync.Mutex
	provider CacheProvider
}

func (cm *CacheManager) Get(key string, value interface{}) error {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	return cm.provider.Get(key, value)
}

func (cm *CacheManager) Set(key string, value interface{}, expireTime int64) error {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	return cm.provider.Set(key, value, expireTime)
}

func (cm *CacheManager) Gc() {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	cm.Gc()
}
