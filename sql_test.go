package ksana

import (
	"log"
	"testing"
)

func TestSql(t *testing.T) {
	log.Printf("==================SQL=============================")
	sql := "INSERT INTO t1(aaa, bbb, ccc) VALUES($1, $2, $3);"
	log.Printf("%s\n%s", sql, re_sql_script.ReplaceAllString(sql, "?"))
}
