package config

import (
	"fmt"
)

type DatabaseConfigData struct {
	Component struct{}
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func (c *DatabaseConfigData) GetConfig() *DatabaseConfig {
	// In real application, this could load from environment variables or config files
	return &DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		Password: "secret",
	}
}

func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/mydb",
		c.Username, c.Password, c.Host, c.Port)
}
