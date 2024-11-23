package service

import (
	"github.com/tuhuynh27/go-ioc/internal/wire/testdata/logger"
	"github.com/tuhuynh27/go-ioc/internal/wire/testdata/service"
)

type NotificationService struct {
	Component   struct{}
	EmailSender service.MessageService `autowired:"true" qualifier:"email"`
	SmsSender   service.MessageService `autowired:"true" qualifier:"sms"`
	Logger      logger.Logger          `autowired:"true" qualifier:"json"`
}

func (s *NotificationService) SendMessage(message string) error {
	s.EmailSender.SendMessage(message)
	s.SmsSender.SendNotification(message)
	s.Logger.Info("Sent message!")
	return nil
}
