package configs

import (
	"log"
	"os"

	"github.com/op/go-logging"
)

// GetLogger return an instance of logging.Logger
func GetLogger(config *Config) *logging.Logger {
	l := logging.MustGetLogger("logger")
	logging.SetFormatter(logging.MustStringFormatter("[%{level:.3s}] %{message}"))
	// Setup stderr backend
	logBackend := logging.NewLogBackend(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	if config.Debug {
		logBackend.Color = true
	}
	// Setup syslog backend
	syslogBackend, err := logging.NewSyslogBackend("glive-consumer: ")
	if err != nil {
		l.Fatal(err)
	}
	// Combine them both into one logging backend.
	logging.SetBackend(logBackend, syslogBackend)
	logging.SetLevel(logging.DEBUG, "logger")
	return l
}
