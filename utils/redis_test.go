package ksana

import (
	"testing"
)

func TestRedis(t *testing.T) {
	key := "aaa"
	val := "bbb"

	redis := Redis{}
	redis.Open("localhost:6379", 5)
	//redis.Run('PING')
	redis.Run("SET", key, val)

	r, _ := redis.Run("GET", key)

	v1, _ := r.Str()
	if v1 != val {
		t.Errorf("val == %i, want %i", v1, val)
	}
}
