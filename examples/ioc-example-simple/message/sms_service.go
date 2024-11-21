package message

import (
	"fmt"

	"github.com/tuhuynh27/go-ioc/examples/ioc-example-simple/logger"
	"github.com/tuhuynh27/go-ioc/ioc"
)

type SMSService struct {
	Component  ioc.Component
	Qualifier  struct{}      `value:"sms"`
	Implements struct{}      `implements:"MessageService"`
	Logger     logger.Logger `autowired:"true"`
}

func (s *SMSService) SendMessage(msg string) string {
	s.Logger.Log("Sending an SMS message")
	return fmt.Sprintf("SMS: %s", msg)
}
