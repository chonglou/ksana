package ksana

import (
	"log"
	"testing"
	"time"
)

var dbCfg = databaseConfig{
	Driver:   "postgres",
	Host:     "localhost",
	Port:     5432,
	Name:     "ksana_test",
	User:     "postgres",
	Password: "",
	Ssl:      "disable"}

type TestBean1 struct {
	Name1 string `sql:"size=155;unique=true;long=true;null=false;index=Created,Updated;default=aaa"`
	Name2 string `sql:"size=155;unique=true;fix=true;null=false;index=Created,Updated;default=aaa"`

	Id       string
	Time     time.Time `sql:"type=time"`
	Date     time.Time `sql:"type=date"`
	Datetime time.Time `sql:"type=datetime"`

	Created time.Time
	Updated time.Time
}

type TestBean2 struct {
	Id int32

	//Int8    int8
	Int32 int32
	Int64 int64
	//Uint8 uint8
	// Uint32    uint32
	// Uint64    uint64

	Rune rune
	//Byte    byte
	String  string `sql:"size=255"`
	Boolean bool
	Float   float32
	Double  float64
	Bytes   []byte

	Versions []string `sql:"-"`
}

func TestModel(t *testing.T) {
	log.Printf("==================MODEl=============================")
	var m Model
	m = &model{sql: &Sql{dialect: &pgDialect{config: &dbCfg}}}

	for _, b := range []interface{}{TestBean1{}, TestBean2{}} {

		c, d, err := m.Table(b)
		if err == nil {
			log.Printf("UP: %s", c)
			log.Printf("DOWN: %s", d)
		} else {
			t.Errorf("Error on register: %v", err)
		}
	}
}
