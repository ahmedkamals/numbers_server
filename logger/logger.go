package logger

import "log"

type Log struct {
	log log.Logger
}

const (
	EMERGENCY = "EMERGENCY"
	ALERT     = "ALERT"
	CRITICAL  = "CRITICAL"
	ERROR     = "ERROR"
	WARNING   = "WARNING"
	NOTICE    = "NOTICE"
	INFO      = "INFO"
	DEBUG     = "DEBUG"
	TRACE     = "TRACE"
)

func NewLogger() *Log{

	return &Log{
		log.Logger{},
	}
}

func (self *Log) Debug(message string, context interface{}) {
	self.log.Println(format(prepareMessage(message, context), DEBUG))
}

func (self *Log) Emergency(message string, context interface{}) {
	self.log.Println(format(prepareMessage(message, context), EMERGENCY))
}

func (self *Log) Alert(message string, context interface{}) {
	self.log.Println(format(prepareMessage(message, context), ALERT))
}

func (self *Log) Critical(message string, context interface{}) {
	self.log.Println(format(prepareMessage(message, context), CRITICAL))
}

func (self *Log) Error(message string, context interface{}) {
	self.log.Println(format(prepareMessage(message, context), ERROR))
}

func (self *Log) Warning(message string, context interface{}) {
	self.log.Println(format(prepareMessage(message, context), WARNING))
}

func (self *Log) Notice(message string, context interface{}) {
	self.log.Println(format(prepareMessage(message, context), NOTICE))
}

func (self *Log) Info(message string, context interface{}) {
	self.log.Println(format(prepareMessage(message, context), INFO))
}

func (self *Log) Trace(message string, context interface{}) {
	self.log.Println(format(prepareMessage(message, context), TRACE))
}
