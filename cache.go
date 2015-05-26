package ksana

type CacheManager interface {
	Get(key string, value interface{}) error
	Set(key string, value interface{}, expireTime int64) error
}
