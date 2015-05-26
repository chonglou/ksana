package ksana

import (
	utils "github.com/chonglou/ksana/utils"
	_ "github.com/lib/pq"
	"testing"
)

var cfg = Config{
	Driver:   "postgres",
	Host:     "localhost",
	Port:     5432,
	Name:     "ksana_test",
	User:     "postgres",
	Password: "",
	Ssl:      "disable"}
var path = "/tmp/migrate"

func TestOrm(t *testing.T) {

	db := Connection{}

	err := db.Open(path, &cfg)
	if err != nil {
		t.Errorf("Error on open: %v", err)
	}

	err = db.Generate("t" + utils.RandomStr(5))
	if err != nil {
		t.Errorf("Error on generate: %v", err)
	}

	err = db.Migrate()
	if err != nil {
		t.Errorf("Error on migrate: %v", err)
	}

	err = db.Rollback()
	if err != nil {
		t.Errorf("Error on rollback: %v", err)
	}

}
