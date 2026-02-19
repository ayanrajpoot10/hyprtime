package logger

import (
	"log"
	"os"
)

// LogLevel represents the logging level
type LogLevel int

const (
	LogLevelQuiet LogLevel = iota
	LogLevelNormal
	LogLevelVerbose
	LogLevelDebug
)

var currentLevel LogLevel = LogLevelNormal

// SetLogLevel sets the global log level
func SetLogLevel(level LogLevel) {
	currentLevel = level
}

// Info logs informational messages (always shown unless quiet)
func Info(format string, v ...interface{}) {
	if currentLevel >= LogLevelNormal {
		log.Printf("[INFO] "+format, v...)
	}
}

// Verbose logs detailed operational messages
func Verbose(format string, v ...interface{}) {
	if currentLevel >= LogLevelVerbose {
		log.Printf("[VERBOSE] "+format, v...)
	}
}

// Debug logs debug messages for troubleshooting
func Debug(format string, v ...interface{}) {
	if currentLevel >= LogLevelDebug {
		log.Printf("[DEBUG] "+format, v...)
	}
}

// Error logs error messages (always shown)
func Error(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}

// Fatal logs fatal errors and exits
func Fatal(format string, v ...interface{}) {
	log.Fatalf("[FATAL] "+format, v...)
}

// Warn logs warning messages
func Warn(format string, v ...interface{}) {
	if currentLevel >= LogLevelNormal {
		log.Printf("[WARN] "+format, v...)
	}
}

// Init initializes the logger
func Init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime)
}

// GetLogLevelString returns a string representation of the log level
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
