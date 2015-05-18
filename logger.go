package ksana

import (
	"log/syslog"
	"os"
)

func openLogger(env string, tag string) (*syslog.Writer, error) {
	var level syslog.Priority
	if env == "production" {
		level = syslog.LOG_INFO
	} else {
		level = syslog.LOG_DEBUG
	}
	return syslog.New(level, tag)

}

var logger, _ = openLogger(os.Getenv("KSANA_ENVIRONMENT"), "ksana")
