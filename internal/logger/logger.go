package logger

import (
	"log"
	"os"
)

type LogLevel int

const (
	LogLevelQuiet LogLevel = iota
	LogLevelNormal
	LogLevelVerbose
	LogLevelDebug
)

var currentLevel LogLevel = LogLevelNormal

func SetLogLevel(level LogLevel) {
	currentLevel = level
}

func Info(format string, v ...interface{}) {
	if currentLevel >= LogLevelNormal {
		log.Printf("[INFO] "+format, v...)
	}
}

func Verbose(format string, v ...interface{}) {
	if currentLevel >= LogLevelVerbose {
		log.Printf("[VERBOSE] "+format, v...)
	}
}

func Debug(format string, v ...interface{}) {
	if currentLevel >= LogLevelDebug {
		log.Printf("[DEBUG] "+format, v...)
	}
}

func Error(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}

func Fatal(format string, v ...interface{}) {
	log.Fatalf("[FATAL] "+format, v...)
}

func Warn(format string, v ...interface{}) {
	if currentLevel >= LogLevelNormal {
		log.Printf("[WARN] "+format, v...)
	}
}

func Init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime)
}

func GetLogLevelString() string {
	switch currentLevel {
	case LogLevelQuiet:
		return "quiet"
	case LogLevelNormal:
		return "normal"
	case LogLevelVerbose:
		return "verbose"
	case LogLevelDebug:
		return "debug"
	default:
		return "unknown"
	}
}
