package ksana

import (
	"crypto/rand"
	"fmt"
	"log"
)

func Uuid() string {
	b := RandomBytes(16)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func RandomStr(size int) string {
	b := RandomBytes(size)
	const dictionary = "0123456789abcdefghijklmnopqrstuvwxyz"
	for k, v := range b {
		b[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(b)

}

func RandomBytes(size int) []byte {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("error on generate random string: %v", err)
	}
	return b
}
