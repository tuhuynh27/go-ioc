package logger

import (
	"fmt"
)

type ConsoleLogger struct {
	Component  struct{}
	Qualifier  struct{} `value:"console"`
	Implements struct{} `implements:"Logger"`
}

func (l *ConsoleLogger) Info(message string) {
	fmt.Println("[INFO] ", message)
}

func (l *ConsoleLogger) Error(message string) {
	fmt.Println("[ERROR] ", message)
}
