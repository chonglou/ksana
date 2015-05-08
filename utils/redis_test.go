package ksana

import (
	"testing"
)

func TestRedis(t *testing.T) {
	key := "aaa"
	val := "bbb"

	redis := Redis{}
	redis.Open("localhost:6379", 5)
	redis.Set(key, val)
	v1 := redis.GetString(key)
	if v1 != val {
		t.Errorf("val == %i, want %i", v1, val)
	}
}
