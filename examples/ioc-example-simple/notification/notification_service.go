package notification

import (
	"fmt"

	"github.com/tuhuynh27/go-ioc/examples/ioc-example-simple/message"
	"github.com/tuhuynh27/go-ioc/ioc"
)

type NotificationService struct {
	Component   ioc.Component
	EmailSender message.MessageService `autowired:"true" qualifier:"email"`
	SmsSender   message.MessageService `autowired:"true" qualifier:"sms"`
}

func (s *NotificationService) SendNotifications(msg string) {
	fmt.Println(s.EmailSender.SendMessage(msg))
	fmt.Println(s.SmsSender.SendMessage(msg))
}
