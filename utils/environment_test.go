package ksana

import (
	"testing"
)

func TestEnvironment(t *testing.T) {
	port := 8080
	envs := make(map[string]Environment)
	envs["development"] = Environment{
		Port: 1234,
	}
	envs["test"] = Environment{}
	envs["production"] = Environment{
		Port: port,
		Database: DatabaseCfg{
			Port: 3306,
		},
	}
	StoreEnvironment(envs)
	if e2 := LoadEnvironment("production"); e2.Port != port {
		t.Errorf("port == %i, want %i", e2.Port, port)
	}
}
