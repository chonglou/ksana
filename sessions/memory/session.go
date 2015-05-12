package ksana_session

import (
	"container/list"
	"github.com/chonglou/ksana"
	"sync"
	"time"
)

var glSessionProvider = &Provider{list: list.New()}



type SessionStore struct {
	sid          string
	timeAccessed time.Time
	value        map[interface{}]interface{}
}

func (st *SessionStore) Set(key, value interface{}) error {
	st.value[key] = value
	glSessionProvider.Update(st.sid)
	return nil
}

func (st *SessionStore) Get(key interface{}) interface{} {
	glSessionProvider.Update(st.sid)
	if v, ok := st.value[key]; ok {
		return v
	} else {
		return nil
	}
	return nil
}

func (st *SessionStore) Delete(key interface{}) error {
	delete(st.value, key)
	glSessionProvider.Update(st.sid)
	return nil
}

func (st *SessionStore) SessionId() string {
	return st.sid
}

type Provider struct {
	lock     sync.Mutex
	sessions map[string]*list.Element
	list     *list.List
}

func (pd *Provider) Init(sid string) (ksana.Session, error) {
	pd.lock.Lock()
	defer pd.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	newSess := &SessionStore{sid: sid, timeAccessed: time.Now(), value: v}
	el := pd.list.PushBack(newSess)
	pd.sessions[sid] = el
	return newSess, nil
}

func (pd *Provider) Read(sid string) (ksana.Session, error) {
	if el, ok := pd.sessions[sid]; ok {
		return el.Value.(*SessionStore), nil
	} else {
		sess, err := pd.Init(sid)
		return sess, err
	}
	return nil, nil
}

func (pd *Provider) Destroy(sid string) error {
	if el, ok := pd.sessions[sid]; ok {
		delete(pd.sessions, sid)
		pd.list.Remove(el)
		return nil
	}
	return nil
}

func (pd *Provider) Gc(maxLifeTime int64) {
	pd.lock.Lock()
	defer pd.lock.Unlock()

	for {
		el := pd.list.Back()
		if el == nil {
			break
		}
		if (el.Value.(*SessionStore).timeAccessed.Unix() + maxLifeTime) < time.Now().Unix() {
			pd.list.Remove(el)
			delete(pd.sessions, el.Value.(*SessionStore).sid)
		} else {
			break
		}
	}
}

func (pd *Provider) Update(sid string) error {
	pd.lock.Lock()
	defer pd.lock.Unlock()
	if el, ok := pd.sessions[sid]; ok {
		el.Value.(*SessionStore).timeAccessed = time.Now()
		pd.list.MoveToFront(el)
		return nil
	}
	return nil
}

func init() {
	glSessionProvider.sessions = make(map[string]*list.Element, 0)

}
