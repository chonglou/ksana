package ksana

import (
	"log"
	"testing"
)

const sid = "test_session_sid"

func TestRedisSession(t *testing.T) {
	log.Printf("==================SESSION=============================")
	r := Redis{}
	err := r.Open(&redisConfig{Host: "localhost", Port: 6379, Db: 2, Pool: 12})
	if err != nil {
		t.Errorf("Open redis error: %v", err)
	}

	session_t(&RedisSessionProvider{redis: &r, maxLifeTime: 600}, t)
}

func session_t(sp SessionProvider, t *testing.T) {
	sess, err := sp.Init(sid)
	if err != nil {
		t.Errorf("Session init error: %v", err)
	}
	key, val := "aaa", 1234
	err = sess.Set(key, val)
	if err != nil {
		t.Errorf("Session set error: %v", err)
	}

	s1, e1 := sp.Read(sid)
	if e1 != nil {
		t.Errorf("Session read error: %v", e1)
	}

	if s1.Get(key) == val {
		log.Printf("test session pass")
	} else {
		t.Errorf("Want %d, Get %d", val, s1.Get(key))
	}

}
