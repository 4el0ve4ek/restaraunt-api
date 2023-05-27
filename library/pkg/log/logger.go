package log

import (
	"log"
	"os"
)

type Logger interface {
	Info(err error)

	Warn(err error)

	Error(err error)

	Fatal(err error)
}

func NewLogger() Logger {
	return &logger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

type logger struct {
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
}

func (l *logger) Info(err error) {
	l.sendMessage(levelInfo, err)
}

func (l *logger) Warn(err error) {
	l.sendMessage(levelWarn, err)
}

func (l *logger) Error(err error) {
	l.sendMessage(levelError, err)
}

func (l *logger) Fatal(err error) {
	l.sendMessage(levelFatal, err)
}

type level int

const (
	levelInfo level = iota
	levelWarn
	levelError
	levelFatal
)

func (l *logger) sendMessage(level level, err error) {
	switch level {
	case levelInfo:
		l.infoLogger.Print(err)
	case levelWarn:
		l.warnLogger.Print(err)
	case levelError:
		l.errorLogger.Print(err)
	case levelFatal:
		l.errorLogger.Fatal(err)

	default:
		panic("unknown message level")
	}
}
