package ksana

import (
	"log"
	"testing"
)

func TestUuid(t *testing.T) {
	log.Printf("UUID: %s\t%s", UUID(), UUID())
	log.Printf("Random string: %s\t%s", RandomStr(16), RandomStr(16))
}
