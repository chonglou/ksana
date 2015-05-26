package ksana

import (
	_ "github.com/lib/pq"
	"log"
	"testing"
)

func TestMigrator(t *testing.T) {
	log.Printf("==================MIGRATION=============================")
	var path = "tmp/migrate"

	m, err := NewMigrator(path, &dbCfg)
	if err != nil {
		t.Errorf("Error on open database: %v", err)
	}

	err = m.Generate("t" + RandomStr(5))
	if err != nil {
		t.Errorf("Error on generate: %v", err)
	}

	err = m.Migrate()
	if err != nil {
		t.Errorf("Error on migrate: %v", err)
	}

	err = m.Rollback()
	if err != nil {
		t.Errorf("Error on rollback: %v", err)
	}

}
