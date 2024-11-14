package logger

import (
	"log"
	"os"
)

var (
	debugLogger   *log.Logger
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
	fatalLogger   *log.Logger
	logLevel      int
)

const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
	FATAL
)

// Init initializes the loggers and sets the log level
func Init() {
	debugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime)
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	warningLogger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
	fatalLogger = log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime)
}

// ChangeLogLevel sets the log level based on the provided string
func ChangeLogLevel(level string) {
	switch level {
	case "DEBUG":
		logLevel = DEBUG
	case "INFO":
		logLevel = INFO
	case "WARNING":
		logLevel = WARNING
	case "ERROR":
		logLevel = ERROR
	case "FATAL":
		logLevel = FATAL
	default:
		logLevel = INFO // default to INFO if the provided level is unrecognized
	}
}

// Debug logs debugging messages
func Debug(v ...interface{}) {
	if logLevel <= DEBUG {
		debugLogger.Println(v...)
	}
}

// Debugf logs formatted debugging messages
func Debugf(format string, v ...interface{}) {
	if logLevel <= DEBUG {
		debugLogger.Printf(format, v...)
	}
}

// Info logs informational messages
func Info(v ...interface{}) {
	if logLevel <= INFO {
		infoLogger.Println(v...)
	}
}

// Infof logs formatted informational messages
func Infof(format string, v ...interface{}) {
	if logLevel <= INFO {
		infoLogger.Printf(format, v...)
	}
}

// Warning logs warning messages
func Warning(v ...interface{}) {
	if logLevel <= WARNING {
		warningLogger.Println(v...)
	}
}

// Warningf logs formatted warning messages
func Warningf(format string, v ...interface{}) {
	if logLevel <= WARNING {
		warningLogger.Printf(format, v...)
	}
}

// Error logs error messages
func Error(v ...interface{}) {
	if logLevel <= ERROR {
		errorLogger.Println(v...)
	}
}

// Errorf logs formatted error messages
func Errorf(format string, v ...interface{}) {
	if logLevel <= ERROR {
		errorLogger.Printf(format, v...)
	}
}

// Fatalf logs fatal errors and exits
func Fatalf(format string, v ...interface{}) {
	fatalLogger.Fatalf(format, v...)
}
