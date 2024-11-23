package notification

import (
	"github.com/tuhuynh27/go-ioc/ioc"
	"github.com/tuhuynh27/go-ioc/ioc/testdata/components/service"
)

type NotificationService struct {
	ioc.Component
	EmailSender service.MessageService `autowired:"" qualifier:"email"`
	SmsSender   service.MessageService `autowired:"" qualifier:"sms"`
}

func (s *NotificationService) NotifyUser(to, message string) error {
	if err := s.EmailSender.Send(to, message); err != nil {
		return err
	}
	return s.SmsSender.Send(to, message)
}
