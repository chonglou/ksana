package ksana_orm

import (
	"testing"
	"time"
)

type TestBean1 struct {
	Name1 string `sql:"size=155;unique=true;long=true;null=false;index=Created,Updated;default=aaa"`
	Name2 string `sql:"size=155;unique=true;fix=true;null=false;index=Created,Updated;default=aaa"`

	Uid     string    `sql:"type=uuid"`
	Time    time.Time `sql:"type=time"`
	Date    time.Time `sql:"type=date"`
	Created time.Time `sql:"type=created"`
	Updated time.Time `sql:"type=updated"`
}

type TestBean2 struct {
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
	db := Database{}

	err := db.Open(path, &cfg)
	if err != nil {
		t.Errorf("Error on open: %v", err)
	}

	// for _, b := range []interface{}{TestBean1{}, TestBean2{}} {
	// 	m := Model{bean: &b}
	// 	err = m.Register(&db)
	// 	if err != nil {
	// 		t.Errorf("Error on register: %v", err)
	// 	}
	// }
}
