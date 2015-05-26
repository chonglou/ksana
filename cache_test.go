package ksana

import (
	"log"
	"testing"
)

const cache_key = "test_cache"

var cache_value = map[string]interface{}{"aaa": "111", "bbb": 222, "ccc": 1.2}

func TestRedisCache(t *testing.T) {
	log.Printf("==================CACHE=============================")
	r := Redis{}
	err := r.Open(&redisCfg)
	if err != nil {
		t.Errorf("Open redis error: %v", err)
	}

	cache_t(&RedisCacheManager{redis: &r}, t)
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
