package main

import (
	"github.com/tuhuynh27/go-ioc/examples/ioc-example-simple/notification"
	"github.com/tuhuynh27/go-ioc/examples/ioc-example-simple/wire"
)

func main() {
	container, err := wire.InitializeContainer()
	if err != nil {
		panic(err)
	}

	notificationServiceType, err := container.Get("NotificationService")
	if err != nil {
		panic(err)
	}

	notificationService := notificationServiceType.(*notification.NotificationService)
	notificationService.SendNotifications("Hello, World!")
}
