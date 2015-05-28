package ksana

import (
	"log/syslog"
	"os"
)

func IsProduction() bool {
	return os.Getenv("KSANA_ENV") == "production"
}

func OpenLogger(tag string) (*syslog.Writer, error) {
	var level syslog.Priority
	if IsProduction() {
		level = syslog.LOG_INFO
	} else {
		level = syslog.LOG_DEBUG
	}
	return syslog.New(level, tag)
}
