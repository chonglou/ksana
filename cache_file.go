package ksana

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type FileCacheProvider struct {
	path string
}

func (fcm *FileCacheProvider) filename(key string) string {
	return fmt.Sprintf("%s/%x", fcm.path, Md5([]byte(key)))
}

func (fcm *FileCacheProvider) Set(key string, value interface{}, expireTime int64) error {
	fn := fcm.filename(key)
	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	en := gob.NewEncoder(f)
	en.Encode(value)

	ex := time.Unix(time.Now().Unix()+expireTime, 0)
	os.Chtimes(fn, ex, ex)
	return nil
}

func (fcm *FileCacheProvider) Get(key string, value interface{}) error {
	fn := fcm.filename(key)
	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	de := gob.NewDecoder(f)
	return de.Decode(value)
}

func (fcm *FileCacheProvider) Gc() {
	files, err := ioutil.ReadDir(fcm.path)
	if err != nil {
		return
	}
	for _, f := range files {
		if f.ModTime().Unix() < time.Now().Unix() {
			os.Remove(f.Name())
		}
	}
}
