package xlog

import (
	"os"
	"sync"
)

type Logger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})
}

type StdLogger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})

	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
}

var (
	fatalLogger   Logger = &NullLogger{}
	errorLogger   Logger = &NullLogger{}
	warningLogger Logger = &NullLogger{}
	infoLogger    Logger = &NullLogger{}
	debugLogger   Logger = &NullLogger{}
	traceLogger   Logger = &NullLogger{}

	loggersMu sync.RWMutex
)

type NullLogger struct {
}

func (l *NullLogger) Print(...interface{}) {
	// it is a null logger
}

func (l *NullLogger) Printf(string, ...interface{}) {
	// it is a null logger
}

func (l *NullLogger) Println(...interface{}) {
	// it is a null logger
}

func SetLoggers(loggers ...Logger) {
	loggersMu.Lock()
	defer loggersMu.Unlock()

	var defLogger, lastLogger Logger
	for i, j := 0, len(loggers); i < 6; i++ {
		var logger Logger
		if defLogger != nil {
			logger = defLogger
		} else {
			if i < j {
				currLogger := loggers[i]
				if currLogger == nil && lastLogger != nil {
					defLogger = lastLogger
					logger = defLogger
				} else {
					logger = currLogger
				}
				lastLogger = currLogger
			}
		}
		if logger == nil {
			logger = &NullLogger{}
		}
		switch i {
		case 0:
			fatalLogger = logger
		case 1:
			errorLogger = logger
		case 2:
			warningLogger = logger
		case 3:
			infoLogger = logger
		case 4:
			debugLogger = logger
		case 5:
			traceLogger = logger
		}
	}
}

func Fatal(args ...interface{}) {
	defer os.Exit(1)
	loggersMu.RLock()
	defer loggersMu.RUnlock()
	fatalLogger.Print(args)
}

func Fatalf(format string, args ...interface{}) {
	defer os.Exit(1)
	loggersMu.RLock()
	defer loggersMu.RUnlock()
	fatalLogger.Printf(format, args)
}

func Fatalln(args ...interface{}) {
	defer os.Exit(1)
	loggersMu.RLock()
	defer loggersMu.RUnlock()
	fatalLogger.Println(args)
}

func Error(args ...interface{}) {
	loggersMu.RLock()
	defer loggersMu.RUnlock()
	errorLogger.Print(args)
}

func Errorf(format string, args ...interface{}) {
	loggersMu.RLock()
	defer loggersMu.RUnlock()
	errorLogger.Printf(format, args)
}

func Errorln(args ...interface{}) {
	loggersMu.RLock()
	defer loggersMu.RUnlock()
	errorLogger.Println(args)
}
