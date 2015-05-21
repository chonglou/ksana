package ksana_web

import (
	"github.com/chonglou/ksana/redis"
	"log"
	"testing"
)

const cache_key = "test_cache"

var cache_value = map[string]interface{}{"aaa": "111", "bbb": 222, "ccc": 1.2} //[]string{"aaa", "bbb", "ccc"}

func TestRedisCache(t *testing.T) {

	r := ksana_redis.Connection{}
	err := r.Open(&ksana_redis.Config{Host: "localhost", Port: 6379, Db: 2, Pool: 12})
	if err != nil {
		t.Errorf("Open redis error: %v", err)
	}

	cache_t(&RedisCacheManager{redis: &r}, t)
}

func TestFileCache(t *testing.T) {
	cache_t(&FileCacheManager{path: "/tmp/ksana/tmp/cache"}, t)
}

func cache_t(cm CacheManager, t *testing.T) {
	cache_value := map[string]interface{}{"aaa": "111", "bbb": 222, "ccc": 1.2}

	err := cm.Set(cache_key, cache_value, 3600*24)
	if err != nil {
		t.Errorf("Error on cache set: %v", err)
	}

	val := make(map[string]interface{}, 0)

	err = cm.Get(cache_key, &val)
	if err != nil {
		t.Errorf("Error on cache get: %v", err)
	}

	if val["bbb"] == cache_value["bbb"] {
		log.Printf("Cache: %v, %v", cache_value, val)
	} else {
		t.Errorf("Want %v, get %v", cache_value, val)
	}
}
