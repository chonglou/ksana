package ksana_utils

import (
	"log/syslog"
	"os"
)

func OpenLogger(tag string) (*syslog.Writer, error) {
	var level syslog.Priority
	if os.Getenv("KSANA_ENV") == "production" {
		level = syslog.LOG_INFO
	} else {
		level = syslog.LOG_DEBUG
	}
	return syslog.New(level, tag)
}
