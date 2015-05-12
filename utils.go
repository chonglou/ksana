package ksana

import (
	"crypto/rand"
	"fmt"
	"log"
)

func UUID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("error on generate uuid: %v", err)
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func RandomStr(size int) string {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("error on generate random string: %v", err)
	}
	return fmt.Sprintf("%x", b)
}
