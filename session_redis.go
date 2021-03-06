package ksana

import (
	"fmt"
	"time"
)

type RedisSessionStore struct {
	SessionStore
	key         string
	maxLifeTime int64
	redis       *Redis
}

func (rss *RedisSessionStore) save() error {
	return rss.redis.Set(rss.key, rss.value, rss.maxLifeTime)
}

func (rss *RedisSessionStore) Set(key, value interface{}) error {
	rss.value[key] = value
	return rss.save()
}

func (rss *RedisSessionStore) Get(key interface{}) interface{} {
	if v, ok := rss.value[key]; ok {
		return v
	}
	rss.redis.Expire(rss.key, rss.maxLifeTime)
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
	redis       *Redis
	maxLifeTime int64
}

func (rsp *RedisSessionProvider) key(sid string) string {
	return fmt.Sprintf("session://%s", sid)
}

func (rsp *RedisSessionProvider) Init(sid string) (Session, error) {

	ss := &RedisSessionStore{
		SessionStore: SessionStore{
			sid:          sid,
			value:        make(map[interface{}]interface{}, 0),
			timeAccessed: time.Now()},
		key:         rsp.key(sid),
		maxLifeTime: rsp.maxLifeTime,
		redis:       rsp.redis}
	err := ss.save()

	return ss, err
}

func (rsp *RedisSessionProvider) Read(sid string) (Session, error) {
	key := rsp.key(sid)

	val := make(map[interface{}]interface{}, 0)
	if err := rsp.redis.Get(key, &val); err == nil {
		return &RedisSessionStore{
			SessionStore: SessionStore{
				sid:          sid,
				value:        val,
				timeAccessed: time.Now()},
			key:         key,
			maxLifeTime: rsp.maxLifeTime,
			redis:       rsp.redis}, nil
	} else {
		return nil, err
	}

}

func (rsp *RedisSessionProvider) Destroy(sid string) error {
	return rsp.redis.Del(rsp.key(sid))
}
