package ksana

import (
	"encoding/gob"
	"io/ioutil"
	"os"
	"time"
)

type FileSessionStore struct {
	SessionStore
	filename string
}

func (sfs *FileSessionStore) save() error {
	f, err := os.Create(sfs.filename)
	if err != nil {
		return err
	}
	defer f.Close()

	en := gob.NewEncoder(f)
	en.Encode(sfs.value)
	return nil
}

func (sfs *FileSessionStore) Set(key, value interface{}) error {
	sfs.value[key] = value
	return sfs.save()
}

func (sfs *FileSessionStore) Get(key interface{}) interface{} {
	if v, ok := sfs.value[key]; ok {
		return v
	}
	return nil
}

func (sfs *FileSessionStore) Delete(key interface{}) error {
	delete(sfs.value, key)
	return sfs.save()
}

func (sfs *FileSessionStore) SessionId() string {
	return sfs.sid
}

type FileSessionProvider struct {
	path string
}

func (fsp *FileSessionProvider) filename(sid string) string {
	return fsp.path + "/tmp/sessions/" + sid
}

func (fsp *FileSessionProvider) Init(sid string) (Session, error) {
	ss := &FileSessionStore{
		SessionStore{
			sid:          sid,
			value:        make(map[interface{}]interface{}, 0),
			timeAccessed: time.Now()},
		fsp.filename(sid)}
	err := ss.save()
	return ss, err
}

func (fsp *FileSessionProvider) Read(sid string) (Session, error) {
	fn := fsp.filename(sid)
	if _, err := os.Stat(fn); err == nil {
		val := make(map[interface{}]interface{}, 0)

		f, err := os.Open(fn)
		defer f.Close()
		os.Chtimes(fn, time.Now(), time.Now())
		if err != nil {
			return nil, err
		}
		de := gob.NewDecoder(f)
		err = de.Decode(&val)
		return &FileSessionStore{
			SessionStore{
				sid:          sid,
				value:        val,
				timeAccessed: time.Now()},
			fsp.filename(sid)}, err

	} else {
		return fsp.Init(sid)
	}

}

func (fsp *FileSessionProvider) Gc(maxLifeTime int64) {
	files, err := ioutil.ReadDir(fsp.filename("/"))
	if err != nil {
		return
	}
	for _, f := range files {
		if (f.ModTime().Unix() + maxLifeTime) < time.Now().Unix() {
			os.Remove(f.Name())
		}
	}
}

func (fsp *FileSessionProvider) Destroy(sid string) error {
	return os.Remove(fsp.filename(sid))
}
