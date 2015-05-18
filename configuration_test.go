package ksana

import (
	"log"
	"testing"
)

func TestConfiguration(t *testing.T) {
	log.Println("========== TEST CONFIGURATION ==========")

	fn := "tmp/config.json"

	cfg1 := configuration{
		Port:     8080,
		Env:      "development",
		Home:     "public",
		Key:      RandomBytes(32),
		Password: RandomBytes(32),
		Session:  sessionC{Name: "_ksana", Secret: RandomBytes(32)},
		Redis:    redisC{Url: "localhost:6379", Db: 0, Pool: 12},
		Database: databaseC{
			Driver: "postgres",
			Url:    "postgres://postgres@localhost/ksana?sslmode=disable"},
	}

	err := writeConfig(&cfg1, fn)
	if err != nil {
		t.Errorf("Write config error: %v", err)
	}

	cfg2 := configuration{}
	err = readConfig(&cfg2, fn)
	if err != nil {
		t.Errorf("Read config error: %v", err)
	}

	if cfg1.Port != cfg2.Port {
		t.Errorf("Read not equal with write")
	}
}
