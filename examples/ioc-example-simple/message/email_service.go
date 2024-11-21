package message

import (
	"fmt"

	"github.com/tuhuynh27/go-ioc/examples/ioc-example-simple/logger"
	"github.com/tuhuynh27/go-ioc/ioc"
)

type EmailService struct {
	Component  ioc.Component
	Qualifier  struct{}      `value:"email"`
	Implements struct{}      `implements:"MessageService"`
	Logger     logger.Logger `autowired:"true"`
}

func (s *EmailService) SendMessage(msg string) string {
	s.Logger.Log("Sending an email message")
	return fmt.Sprintf("Email: %s", msg)
}
