package ksana

import (
	"testing"
)

type P struct {
	Aaa string
	Bbb int
	Ccc float32
}

func TestRedis(t *testing.T) {
	r := Redis{}
	err := r.Open(&redisConfig{Host: "localhost", Port: 6379, Db: 2, Pool: 12})
	if err != nil {
		t.Errorf("Open redis error: %v", err)
	}

	key := "aaa"
	val := P{Aaa: "3.14", Bbb: 3, Ccc: 3.14}

	err = r.Set(key, val, 120)
	if err != nil {
		t.Errorf("Open redis set: %v", err)
	}

	val1 := P{}
	err = r.Get(key, &val1)
	if err != nil {
		t.Errorf("Open redis get: %v", err)
	}
	if val.Ccc != val1.Ccc {
		t.Errorf("Want: %v, Get %v", val, val1)
	}

	// val2 := P{}
	// err = r.Cache("cache://"+key, &val2, func(v interface{}) error {
	// 	v = &P{Aaa: "3.14", Bbb: 3, Ccc: 3.14}
	// 	return nil
	// }, 30)
	//
	// if err != nil {
	// 	t.Errorf("Open redis get: %v", err)
	// }
	// if val2.Ccc != val.Ccc {
	// 	t.Errorf("Want: %v, Get %v", val, val2)
	// }
}
