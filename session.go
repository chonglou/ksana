package ksana

import (
	//	"log"
	"sync"
)

type Session interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionID() string
}

type sessionProvider interface {
	init(sid string) (Session, error)
	read(sid string) (Session, error)
	destroy(sid string) error
	gc(maxLifeTime int64)
}

type sessionManager struct {
	cookieName  string
	lock        sync.Mutex
	provider    sessionProvider
	maxLifeTime int64
}

func newSessionManager(provideName, cookieName string,
	maxLifeTime int64) (*sessionManager, error) {
	// provider, ok := providers[provideName]
	//
	// if !ok {
	// 	log.Fatalf("session: unknown provide %s", provideName)
	// }

	return &sessionManager{
		//provider:    provider,
		cookieName:  cookieName,
		maxLifeTime: maxLifeTime}, nil
}
