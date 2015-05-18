package ksana

import (
	_ "github.com/lib/pq"
	"testing"
)

func TestContext(t *testing.T) {
	if err := ctx.Load(config_file); err != nil {
		t.Errorf("Error on init context %v", err)
	}
}
