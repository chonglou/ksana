package ksana

import (
	// "database/sql"
	// _ "github.com/lib/pq"
	"log"
	"testing"
	"time"
)

type TestBean struct {
	Bean
	Name string `sql:"size:155;unique:true;long:true;null:false;index:true;default:aaa"`

	Uid string `sql:"type:uuid"`
	Id  int    `sql:"type:serial"`

	Time time.Time

	Int    int
	Int8   int8
	Int32  int32
	Int64  int64
	Uint   uint
	Rune   rune
	Byte   byte
	String string
	Enable bool
	Float  float32
	Double float64
}

func TestMigration(t *testing.T) {
	log.Println("============== TEST MIGRATION ======================")

	var m Migration
	m = &migration{path: "tmp/migrate"}

	err := m.Add(TestBean{})
	if err != nil {
		t.Errorf("Error on add bean: %v", err)
	}

}
