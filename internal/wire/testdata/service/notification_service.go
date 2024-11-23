package service

import (
	"github.com/tuhuynh27/go-ioc/internal/wire/testdata/config"
	"github.com/tuhuynh27/go-ioc/internal/wire/testdata/logger"
)

type NotificationService struct {
	Component              struct{}
	ConsoleLogger          logger.Logger                  `autowired:"true" qualifier:"console"`
	JsonLogger             logger.Logger                  `autowired:"true" qualifier:"json"`
	NotificationConfigData *config.NotificationConfigData `autowired:"true"`
}

func (s *NotificationService) SendNotification(message string) {
	s.ConsoleLogger.Info("Sending notification via console: " + message + " to " + s.NotificationConfigData.GetConfig().APIUrl)
	s.JsonLogger.Info("Sending notification via JSON: " + message + " to " + s.NotificationConfigData.GetConfig().APIUrl)
}
