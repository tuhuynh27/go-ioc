package main

import (
	"github.com/tuhuynh27/go-ioc/examples/ioc-example-simple/logger"
	"github.com/tuhuynh27/go-ioc/examples/ioc-example-simple/message"
	"github.com/tuhuynh27/go-ioc/examples/ioc-example-simple/notification"
	"github.com/tuhuynh27/go-ioc/ioc"
)

func main() {
	container := ioc.NewContainer()

	err := container.RegisterComponents(
		&logger.StdoutLogger{},
		&message.EmailService{},
		&message.SMSService{},
		&notification.NotificationService{},
	)

	if err != nil {
		panic(err)
	}

	notificationService := container.Get("NotificationService").(*notification.NotificationService)
	notificationService.SendNotifications("Hello, World!")
}
