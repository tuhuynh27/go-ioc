package service

import (
	"github.com/tuhuynh27/go-ioc/internal/wire/testdata/config"
	"github.com/tuhuynh27/go-ioc/internal/wire/testdata/logger"
)

type SMSMessageService struct {
	Component  struct{}           `implements:"MessageService"`
	Qualifier  struct{}           `value:"sms"`
	ConfigData *config.ConfigData `autowired:"true"`
	Logger     logger.Logger      `autowired:"true" qualifier:"console"`
}

func (s *SMSMessageService) SendNotification(message string) error {
	url := s.ConfigData.GetConfig().SMSUrl
	s.Logger.Info("Sending SMS message " + message + " to " + url)
	return nil
}
