package ksana

import (
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)



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
	Update(sid string) error
	Gc(maxLifeTime int64)
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
	sm.provider.Gc(sm.maxLifeTime)
	time.AfterFunc(time.Duration(sm.maxLifeTime), func() { sm.GC() })
}

var glSessionProviders = make(map[string]SessionProvider)

func SessionRegister(name string, provider SessionProvider) {
	if glSessionProviders == nil {
		log.Fatalf("Session provider is nil")
	}
	if _, dup := glSessionProviders[name]; dup {
		log.Fatalf("Register called twice for provide %s", name)
	}
	glSessionProviders[name] = provider
}

func NewSessionManager(providerName, cookieName string,
	maxLifeTime int64) (*SessionManager, error) {

	provider, ok := glSessionProviders[providerName]

	if !ok {
		log.Fatalf("Unknown session provide %s", providerName)
	}

	return &SessionManager{
		provider:    provider,
		cookieName:  cookieName,
		maxLifeTime: maxLifeTime}, nil
}
