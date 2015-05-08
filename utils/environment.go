package ksana

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type DatabaseCfg struct {
	Adapter  string `json:"adapter"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	Pool     int    `json:"pool"`
	Timeout  int    `json:"timeout"`
}

type RedisCfg struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	Db   int    `json:"db"`
	Pool int    `json:"pool"`
}
type SecretCfg struct {
	KeyBase string `json:"key_base"`
}
type Environment struct {
	Port     int         `json:"port"`
	Database DatabaseCfg `json:"database"`
	Redis    RedisCfg    `json:"redis"`
	Secret   SecretCfg   `json:"secret"`
}

func StoreEnvironment(envs map[string]Environment) {
	cj, err := json.MarshalIndent(envs, "", "\t")
	if err != nil {
		log.Fatalf("Error on generate json: %v", err)
	}
	err = ioutil.WriteFile("config/settings.json", cj, 0600)
	if err != nil {
		log.Fatalf("Error on write config file: %v", err)
	}
}

func LoadEnvironment(mode string) Environment {
	f, err := os.Open("config/settings.json")
	if err != nil {
		log.Fatalf("Error on open config file: %v", err)
	}
	defer f.Close()

	envs := make(map[string]Environment)
	if err = json.NewDecoder(f).Decode(&envs); err != nil {
		log.Fatalf("Error on parse config file: %v", err)
	}

	return envs[mode]
}
