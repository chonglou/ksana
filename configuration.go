package ksana

import (
	"encoding/json"
	"io/ioutil"
	"os"

	orm "github.com/chonglou/ksana/orm"
	redis "github.com/chonglou/ksana/redis"
	web "github.com/chonglou/ksana/web"
)

type configuration struct {
	file string

	Env      string      `json:"env"`
	Secret   []byte      `json:"secret"`
	Web      web.Config  `json:"web"`
	Database orm.Config  `json:"database"`
	Redis    redisConfig `json:"redis"`
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
