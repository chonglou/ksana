package ksana_web

import (
	"log"
	"testing"
	redis "github.com/chonglou/ksana/redis"
)

const cache_key = "test_cache"

var cache_value = map[string]interface{}{"aaa": "111", "bbb": 222, "ccc": 1.2} //[]string{"aaa", "bbb", "ccc"}

func TestRedisCache() {
	r := Redis{}
	err := r.Open(&Config{Host: "localhost", Port: 6379, Db: 2, Pool: 12})
	if err != nil {
		t.Errorf("Open redis error: %v", err)
	}
	cache(&CacheManager{provider:&RedisCacheProvider{redis: &r}})

}

func TestFileCache(t *testing.T) {
	cache(&CacheManager{provider: &FileCacheProvider{path: "/tmp/ksana/tmp/cache"}})
}

func cache(cm *CacheManager) {
	err := cm.Set(cache_key, cache_value, 3600*24)
	if err != nil {
		t.Errorf("Error on cache set: %v", err)
	}

	val := make(map[string]interface{}, 0)

	err = cm.Get(cache_key, &val)
	if err != nil {
		t.Errorf("Error on cache get: %v", err)
	}

	if len(val) == len(cache_value) {
		log.Printf("Cache: %v, %v", cache_value, val)
	} else {
		t.Errorf("Want %v, get %v", cache_value, val)
	}
}
