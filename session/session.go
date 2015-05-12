package ksana

import (
	"github.com/chonglou/ksana/utils"
	//"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type SessionStory struct {
	sid          string
	timeAccessed time.Time
	value        map[interface{}]interface{}
}

type Session interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionID() string
}

type SessionProvider interface {
	Init(sid string) (Session, error)
	Read(sid string) (Session, error)
	Destroy(sid string) error
	GC(maxLifeTime int64)
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
		sid := UUID()
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
		sess, _ = sm.provider.Read(sid)
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

func (sm *SessionManager) GC() {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	sm.provider.GC(sm.maxLifeTime)
	time.AfterFunc(time.Duration(sm.maxLifeTime), func() { sm.GC() })
}
