package service

import (
	"github.com/tuhuynh27/go-ioc/ioc"
	"github.com/tuhuynh27/go-ioc/ioc/testdata/components/logger"
)

type MessageService interface {
	Send(to, message string) error
}

type EmailService struct {
	ioc.Component `implements:"service.MessageService" value:"email"`
	Logger        logger.Logger `autowired:""`
}

func (s *EmailService) Send(to, message string) error {
	s.Logger.Log("Sending email to " + to + ": " + message)
	return nil
}
