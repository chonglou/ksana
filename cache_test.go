package ksana

import (
	"log"
	"testing"
)

const cache_key = "test_cache"

var cache_value = map[string]interface{}{"aaa": "111", "bbb": 222, "ccc": 1.2} //[]string{"aaa", "bbb", "ccc"}

func TestFileCache(t *testing.T) {

	cm := CacheManager{provider: &FileCacheProvider{path: "/tmp/ksana/tmp/cache"}}

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
