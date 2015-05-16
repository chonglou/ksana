package ksana

import (
	"testing"
	"log"
)

func TestContext(t *testing.T) {
	cfg := configuration{}
	log.Printf("====================TEST CONTEXT======================")
	err := loadConfiguration("examples/context.xml", &cfg)
	if err != nil {
		t.Errorf("Load config error: %v", err)
	}
}
