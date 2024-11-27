package config

type Config struct {
	APIUrl string
	SMSUrl string
}

type ConfigData struct {
	Component struct{}
}

func NewConfigData() *ConfigData {
	return &ConfigData{}
}

func (c *ConfigData) GetConfig() *Config {
	return &Config{
		APIUrl: "https://api.notification.com",
		SMSUrl: "https://sms.notification.com",
	}
}
