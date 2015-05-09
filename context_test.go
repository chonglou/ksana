package ksana

import (
	_ "github.com/lib/pq"
	"testing"
)

func TestRedis(t *testing.T) {
	c := Context{}
	c.Init()
}
