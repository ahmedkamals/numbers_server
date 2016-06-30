package logger

import (
	baseLog "log"
	"os"
)

type Log struct {
	baseLogger *baseLog.Logger
}

type Context interface{}

type Level uint

const (
	EMERGENCY Level = 8 - iota
	ALERT
	CRITICAL
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
	TRACE
)

func NewLogger() *Log{

	return &Log{
		baseLogger:baseLog.New(os.Stderr, "", baseLog.LstdFlags),
	}
}

func (self *Log) Emergency(message string, context ...Context) {
	self.baseLogger.Fatal(format(prepareMessage(message, context), EMERGENCY))
}

func (self *Log) Alert(message string, context ...Context) {
	self.baseLogger.Println(format(prepareMessage(message, context), ALERT))
}

func (self *Log) Critical(message string, context ...Context) {
	self.baseLogger.Println(format(prepareMessage(message, context), CRITICAL))
}

func (self *Log) Error(message string, context ...Context) {
	self.baseLogger.Println(format(prepareMessage(message, context), ERROR))
}

func (self *Log) Warning(message string, context ...Context) {
	self.baseLogger.Println(format(prepareMessage(message, context), WARNING))
}

func (self *Log) Notice(message string, context ...Context) {
	self.baseLogger.Println(format(prepareMessage(message, context), NOTICE))
}

func (self *Log) Info(message string, context ...Context) {
	self.baseLogger.Println(format(prepareMessage(message, context), INFO))
}

func (self *Log) Debug(message string, context ...Context) {
	self.baseLogger.Println(format(prepareMessage(message, context), DEBUG))
}

func (self *Log) Trace(message string, context ...Context) {
	self.baseLogger.Println(format(prepareMessage(message, context), TRACE))
}


func levelAsString(l Level) string {
	switch l {
	case EMERGENCY:
		return "EMERGENCY"
	case ALERT:
		return "ALERT"
	case CRITICAL:
		return "CRITICAL"
	case ERROR:
		return "ERROR"
	case WARNING:
		return "WARNING"
	case NOTICE:
		return "NOTICE"
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	case TRACE:
		return "TRACE"
	}

	return "TRACE"
}
