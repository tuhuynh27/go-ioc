package logger

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tuhuynh27/go-ioc/ioc"
)

type LogLevel string

const (
	INFO  LogLevel = "INFO"
	ERROR LogLevel = "ERROR"
	DEBUG LogLevel = "DEBUG"
)

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     LogLevel  `json:"level"`
	Message   string    `json:"message"`
}

type Logger interface {
	Log(message string)
	LogWithLevel(level LogLevel, message string)
}

// ConsoleLogger implements simple console logging
type ConsoleLogger struct {
	ioc.Component `implements:"logger.Logger" value:"console"`
}

func (l *ConsoleLogger) Log(message string) {
	l.LogWithLevel(INFO, message)
}

func (l *ConsoleLogger) LogWithLevel(level LogLevel, message string) {
	fmt.Printf("[%s] %s: %s\n", time.Now().Format(time.RFC3339), level, message)
}

// JsonLogger implements JSON-formatted logging
type JsonLogger struct {
	ioc.Component `implements:"logger.Logger" value:"json"`
}

func (l *JsonLogger) Log(message string) {
	l.LogWithLevel(INFO, message)
}

func (l *JsonLogger) LogWithLevel(level LogLevel, message string) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
	}
	jsonData, _ := json.Marshal(entry)
	fmt.Println(string(jsonData))
}
