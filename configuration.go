package ksana

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type redisC struct {
	Url  string `json:"url"`
	Db   int    `json:"db"`
	Pool int    `json:"pool"`
}

type databaseC struct {
	Driver string `json:"driver"`
	Url    string `json:"url"`
}

type sessionC struct {
	Name   string `json:"name"`
	Secret []byte `json:"secret"`
}

type configuration struct {
	Name     string    `json:"name"`
	Port     int       `json:"port"`
	Env      string    `json:"env"`
	Home     string    `json:"home"`
	Key      []byte    `json:"key"`
	Password []byte    `json:"password"`
	Session  sessionC  `json:"session"`
	Redis    redisC    `json:"redis"`
	Database databaseC `json:"database"`
}

func writeConfig(cfg *configuration, file string) error {
	cj, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, cj, 0600)
}

func readConfig(cfg *configuration, file string) error {
	f, e := os.Open(file)
	if e != nil {
		return e
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(cfg)
}
