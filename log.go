package main

import "log"

// Logger is a simple interface for logger.
type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type stdLogger struct{}

func (stdLogger) Infof(format string, args ...interface{}) {
	log.Printf("[INFO]: "+format, args...)
}

func (stdLogger) Errorf(format string, args ...interface{}) {
	log.Printf("[ERROR]: "+format, args...)
}
