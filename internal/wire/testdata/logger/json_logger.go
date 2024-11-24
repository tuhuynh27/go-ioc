package logger

import (
	"fmt"
)

type JsonLogger struct {
	Component struct{} `implements:"Logger"`
	Qualifier struct{} `value:"json"`
}

func (l *JsonLogger) Info(message string) {
	fmt.Printf("{\"level\":\"info\",\"message\":\"%s\"}\n", message)
}

func (l *JsonLogger) Error(message string) {
	fmt.Printf("{\"level\":\"error\",\"message\":\"%s\"}\n", message)
}

// New PostConstruct method
func (l *JsonLogger) PostConstruct() {
	fmt.Println("JsonLogger initialized")
}

// New PreDestroy method
func (l *JsonLogger) PreDestroy() {
	fmt.Println("JsonLogger is being destroyed")
}
