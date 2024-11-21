package logger

import (
	"fmt"

	"github.com/tuhuynh27/go-ioc/ioc"
)

type StdoutLogger struct {
	Component  ioc.Component
	Implements struct{} `implements:"github.com/tuhuynh27/go-ioc/examples/ioc-example-simple/logger.Logger"`
}

func (s *StdoutLogger) Log(msg string) {
	fmt.Println("Logger: " + msg)
}
