package log

import (
	"io"
	"time"

	"../trace"
)

type LogLevel = uint8

const (
	OFF   LogLevel = 0
	FATAL LogLevel = 1
	ERROR LogLevel = 2
	WARN  LogLevel = 3
	INFO  LogLevel = 4
	DEBUG LogLevel = 5
	TRACE LogLevel = 6
)

type ILogger interface {
	Log(message string, level LogLevel)
	LogTrace(etrace trace.ITrace)
}

type Logger struct {
	Level LogLevel

	ErrOut   io.Writer
	WarnOut  io.Writer
	InfoOut  io.Writer
	TraceOut io.Writer
}

func format(label, message string) string {
	return "[" + label + "] " + time.Now().String() + "\n" + message + "\n\n"
}

func (logger *Logger) dispatch(message string, level LogLevel) {
	switch level {
	case FATAL, ERROR:
		logger.InfoOut.Write([]byte(format("Error", message)))
	case WARN:
		logger.InfoOut.Write([]byte(format("Warning", message)))
	case INFO, TRACE:
		logger.InfoOut.Write([]byte(format("Info", message)))
	}
}

func (logger *Logger) Log(message string, level LogLevel) {
	if level <= logger.Level {
		logger.dispatch(message, level)
	}
}

func (logger *Logger) LogTrace(etrace trace.ITrace) {
	logger.Log(etrace.String(), etrace.SafetyLevel())
}
