package logger

import (
	"log"
	"strings"
)

type LogLevel int

const (
	FatalLevel LogLevel = iota
	ErrorLevel
	InfoLevel
	DebugLevel
)

type Logger struct {
	level LogLevel
	*log.Logger
}

var logger = Logger{FatalLevel, log.Default()}

func SetLevel(l LogLevel) {
	logger.level = l
}

func ParseLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "fatal":
		return FatalLevel
	case "error":
		return ErrorLevel
	case "info":
		return InfoLevel
	case "debug":
		return DebugLevel
	default:
		log.Fatalf("unknown log level '%s'", level)
	}
	return FatalLevel
}

func SetLevelFromString(level string) {
	SetLevel(ParseLevel(level))
}

func Debug(v ...interface{}) {
	if logger.level >= DebugLevel {
		logger.Print("DEBUG ", v)
	}
}

func Debugf(format string, v ...interface{}) {
	if logger.level >= DebugLevel {
		logger.Printf("DEBUG "+format, v...)
	}
}

func Info(v ...interface{}) {
	if logger.level >= InfoLevel {
		logger.Print("INFO ", v)
	}
}

func Infof(format string, v ...interface{}) {
	if logger.level >= InfoLevel {
		logger.Printf("INFO "+format, v...)
	}
}

func Error(v ...interface{}) {
	if logger.level >= ErrorLevel {
		logger.Print("ERROR ", v)
	}
}

func Errorf(format string, v ...interface{}) {
	if logger.level >= ErrorLevel {
		logger.Printf("ERROR "+format, v...)
	}
}

func Fatal(v ...interface{}) {
	logger.Fatal("FATAL ", v)
}

func Fatalf(format string, v ...interface{}) {
	logger.Fatalf("FATAL "+format, v...)
}
