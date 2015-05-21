package ksana_i18n

import (
	"log"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	Load("tmp")
	for _, l := range []string{"en", "zh_CN"} {
		log.Printf("%s: %v", l, T(l, "hello", time.Now()))
	}
}
