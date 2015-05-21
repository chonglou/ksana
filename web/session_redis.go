package ksana_web

import (
	"fmt"
	"github.com/chonglou/ksana/redis"
	"time"
)

type RedisSessionStore struct {
	SessionStore
	key   string
	redis *ksana_redis.Connection
}

func (rss *RedisSessionStore) save() error {
	return rss.redis.Set(rss.key, rss.value, 0)
}

func (rss *RedisSessionStore) Set(key, value interface{}) error {
	rss.value[key] = value
	return rss.save()
}

func (rss *RedisSessionStore) Get(key interface{}) interface{} {
	if v, ok := rss.value[key]; ok {
		return v
	}
	return nil
}

func (rss *RedisSessionStore) Delete(key interface{}) error {
	delete(rss.value, key)
	return rss.save()
}

func (rss *RedisSessionStore) SessionId() string {
	return rss.sid
}

type RedisSessionProvider struct {
	redis *ksana_redis.Connection
}

func (rsp *RedisSessionProvider) key(sid string) string {
	return fmt.Sprintf("session://%s", sid)
}

func (rsp *RedisSessionProvider) Init(sid string) (Session, error) {
	ss := &RedisSessionStore{
		SessionStore{
			sid:          sid,
			value:        make(map[interface{}]interface{}, 0),
			timeAccessed: time.Now()},
		rsp.key(sid),
		rsp.redis}
	err := ss.save()
	return ss, err
}

func (rsp *RedisSessionProvider) Read(sid string) (Session, error) {
	val := make(map[interface{}]interface{}, 0)
	if err := rsp.redis.Get(rsp.key(sid), &val); err == nil {
		return &RedisSessionStore{
			SessionStore{
				sid:          sid,
				value:        val,
				timeAccessed: time.Now()},
			rsp.key(sid),
			rsp.redis}, nil
	} else {
		return nil, err
	}

}

func (rsp *RedisSessionProvider) Destroy(sid string) error {
	return rsp.redis.Del(rsp.key(sid))
}

func (rsp *RedisSessionProvider) Gc(maxLifeTime int64) {
	logger.Info("Session gc!!")
}
