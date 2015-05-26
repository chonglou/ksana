package ksana

import (
	"encoding/gob"
	"fmt"
	utils "github.com/chonglou/ksana/utils"
	"io/ioutil"
	"os"
	"sync"
	//"time"
)

type FileCacheManager struct {
	path string
	lock sync.Mutex
}

func (fcm *FileCacheManager) filename(key string) string {
	return fmt.Sprintf("%s/%x", fcm.path, utils.Md5([]byte(key)))
}

func (fcm *FileCacheManager) Set(key string, value interface{}, expireTime int64) error {
	fcm.lock.Lock()
	defer fcm.lock.Unlock()

	fn := fcm.filename(key)
	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Chmod(0600)

	en := gob.NewEncoder(f)
	en.Encode(value)

	// ex := time.Unix(time.Now().Unix()+expireTime, 0)
	// os.Chtimes(fn, ex, ex)

	//time.AfterFunc(time.Duration(expireTime), func() { os.Remove(fn) })
	return nil
}

func (fcm *FileCacheManager) Get(key string, value interface{}) error {
	fcm.lock.Lock()
	defer fcm.lock.Unlock()

	fn := fcm.filename(key)
	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	de := gob.NewDecoder(f)
	return de.Decode(value)
}

func (fcm *FileCacheManager) Gc() {
	fcm.lock.Lock()
	defer fcm.lock.Unlock()

	files, err := ioutil.ReadDir(fcm.path)
	if err != nil {
		return
	}
	for _, f := range files {
		os.Remove(f.Name())
	}
}
