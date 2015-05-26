package ksana

import (
	"log"
	"testing"
)

const config_file = "tmp/config.json"

func TestConfiguration(t *testing.T) {
	log.Println("========== TEST CONFIGURATION ==========")

	cfg1 := configuration{
		Env:    "development",
		Secret: RandomBytes(512),
		Web: webConfig{
			Port:   8080,
			Cookie: RandomStr(8),
			Expire: 60 * 30},
		Redis: redisConfig{Host: "localhost", Port: 6379, Db: 0, Pool: 12},
		Database: dbConfig{
			Driver:   "postgres",
			Host:     "localhost",
			Port:     5432,
			Name:     "ksana",
			User:     "postgres",
			Password: "",
			Ssl:      "disable"},
	}

	err := writeConfig(&cfg1, config_file)
	if err != nil {
		t.Errorf("Write config error: %v", err)
	}

	cfg2 := configuration{}
	err = readConfig(&cfg2, config_file)
	if err != nil {
		t.Errorf("Read config error: %v", err)
	}

	if cfg1.Web.Port != cfg2.Web.Port {
		t.Errorf("Read not equal with write")
	}
}
