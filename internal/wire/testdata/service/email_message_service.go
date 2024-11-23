package service

import (
	"github.com/tuhuynh27/go-ioc/examples/ioc-example-simple/logger"
	"github.com/tuhuynh27/go-ioc/internal/wire/testdata/config"
)

type EmailMessageService struct {
	Component  struct{}
	Implements struct{}           `implements:"MessageService"`
	Qualifier  struct{}           `value:"email"`
	ConfigData *config.ConfigData `autowired:"true"`
	Logger     logger.Logger      `autowired:"true" qualifier:"console"`
}

func (s *EmailMessageService) SendMessage(message string) error {
	url := s.ConfigData.GetConfig().APIUrl
	s.Logger.Info("Sending Email message " + message + " to " + url)
	return nil
}
