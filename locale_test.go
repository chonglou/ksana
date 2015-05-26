package ksana

import (
	"log"
	"testing"
	"time"
)

func TestI18n(t *testing.T) {
	Load("tmp")
	for _, l := range []string{"en", "zh_CN"} {
		log.Printf("%s: %v", l, T(l, "hello", time.Now()))
	}
}
