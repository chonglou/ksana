package ksana

import (
	"log"
	"log/syslog"
	"os"
)

func OpenLogger(tag string) *syslog.Writer {
	var level syslog.Priority
	if os.Getenv("KSANA_ENVIRONMENT") == "production" {
		level = syslog.LOG_INFO
	} else {
		level = syslog.LOG_DEBUG
	}
	logger, err := syslog.New(level, tag)
	if err != nil {
		log.Fatalf("error on open syslog: %v", err)
	}
	return logger
}

var logger = OpenLogger("ksana")
