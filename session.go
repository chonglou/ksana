package ksana

import (
	"net/http"
	"net/url"
	"sync"
	"time"
)

type SessionStore struct {
	sid          string
	value        map[interface{}]interface{}
	timeAccessed time.Time
}

type Session interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionId() string
}

type SessionProvider interface {
	Init(sid string) (Session, error)
	Read(sid string) (Session, error)
	Destroy(sid string) error
}

type SessionManager struct {
	cookieName  string
	lock        sync.Mutex
	provider    SessionProvider
	maxLifeTime int64
}

func (sm *SessionManager) Start(wrt http.ResponseWriter,
	req *http.Request) (sess Session) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	cke, err := req.Cookie(sm.cookieName)
	if err != nil || cke.Value == "" {
		sid := Uuid()
		sess, _ = sm.provider.Init(sid)
		cke := http.Cookie{
			Name:     sm.cookieName,
			Value:    url.QueryEscape(sid),
			Path:     "/",
			HttpOnly: true,
			MaxAge:   int(sm.maxLifeTime),
		}
		http.SetCookie(wrt, &cke)
	} else {
		sid, _ := url.QueryUnescape(cke.Value)
		sess, err = sm.provider.Read(sid)
		if err != nil {
			sess, _ = sm.provider.Init(sid)
		}
	}
	return sess
}

func (sm *SessionManager) Destroy(wrt http.ResponseWriter, req *http.Request) {
	cke, err := req.Cookie(sm.cookieName)
	if err != nil || cke.Value == "" {
		return
	}
	sm.lock.Lock()
	defer sm.lock.Unlock()
	sm.provider.Destroy(cke.Value)
	http.SetCookie(wrt, &http.Cookie{
		Name:     sm.cookieName,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now(),
		MaxAge:   -1,
	})
}
