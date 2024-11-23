package service

import (
	"github.com/tuhuynh27/go-ioc/ioc"
	"github.com/tuhuynh27/go-ioc/ioc/testdata/components/logger"
)

type SmsService struct {
	ioc.Component `implements:"service.MessageService" value:"sms"`
	Logger        logger.Logger `autowired:""`
}

func (s *SmsService) Send(to, message string) error {
	s.Logger.Log("Sending SMS to " + to + ": " + message)
	return nil
}
