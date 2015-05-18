package ksana

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"testing"
	"time"
)

type TestBean1 struct {
	Bean

	Name1 string `sql:"size=155;unique=true;long=true;null=false;index=Created,Updated;default=aaa"`
	Name2 string `sql:"size=155;unique=true;fix=true;null=false;index=Created,Updated;default=aaa"`

	Uid     string    `sql:"type=uuid"`
	Time    time.Time `sql:"type=time"`
	Date    time.Time `sql:"type=date"`
	Created time.Time `sql:"type=created"`
	Updated time.Time `sql:"type=updated"`
}

type TestBean2 struct {
	Bean

	Id int `sql:"type=serial"`

	Int    int
	Int8   int8
	Int32  int32
	Int64  int64
	Uint   uint
	Rune   rune
	Byte   byte
	String string `sql:"size=255"`
	Enable bool
	Float  float32
	Double float64
}

func TestModel(t *testing.T) {
	log.Println("============== TEST MODEL ======================")

	db, err := sql.Open("postgres", "postgres://postgres@localhost/ksana_t?sslmode=disable")

	if err != nil {
		t.Errorf("Open database: %v", err)
	}
	defer db.Close()
	var m Model
	m = &model{path: "/tmp/migrate", db: db}

	for _, b := range []Bean{TestBean1{}, TestBean2{}} {
		err := m.Register(b)
		if err != nil {
			t.Errorf("Error on add bean: %v", err)
		}

	}

	err = m.Migrate()
	if err != nil {
		t.Errorf("Error on add bean: %v", err)
	}
	err = m.Rollback()
	if err != nil {
		t.Errorf("Error on add bean: %v", err)
	}
}
