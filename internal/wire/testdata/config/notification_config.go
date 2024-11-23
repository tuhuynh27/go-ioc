package config

type NotificationConfigData struct {
	Component struct{}
}

type NotificationConfig struct {
	APIUrl string
}

func (c *NotificationConfigData) GetConfig() *NotificationConfig {
	return &NotificationConfig{
		APIUrl: "http://api.notification.com",
	}
}
