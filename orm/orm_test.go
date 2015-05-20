package ksana_orm

import (
	utils "github.com/chonglou/ksana/utils"
	"testing"
)

func TestOrm(t *testing.T) {

	db := Database{}

	err := db.Open("/tmp/migrate", &Config{
		Driver:   "postgres",
		Host:     "localhost",
		Port:     5432,
		Name:     "ksana_test",
		User:     "postgres",
		Password: "",
		Ssl:      "disable"})
	if err != nil {
		t.Errorf("Error on open: %v", err)
	}

	err = db.Generate(utils.RandomStr(5))
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
