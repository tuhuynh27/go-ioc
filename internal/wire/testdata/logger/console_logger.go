package logger

import (
	"fmt"
)

type ConsoleLogger struct {
	Component struct{} `implements:"Logger"`
	Qualifier struct{} `value:"console"`
}

func (l *ConsoleLogger) Info(message string) {
	fmt.Println("[INFO] ", message)
}

func (l *ConsoleLogger) Error(message string) {
	fmt.Println("[ERROR] ", message)
}
