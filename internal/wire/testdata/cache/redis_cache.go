package cache

import (
	"github.com/tuhuynh27/go-ioc/internal/wire/testdata/logger"
)

type RedisCache struct {
	Component  struct{}      `name:"redisCache"`
	Qualifier  struct{}      `value:"redis"`
	Implements struct{}      `implements:"Cache"`
	Logger     logger.Logger `autowired:"true" qualifier:"console"`
}

func (c *RedisCache) Get(key string) (string, error) {
	c.Logger.Info("Getting key from Redis: " + key)
	return "", nil
}

func (c *RedisCache) Set(key string, value string) error {
	c.Logger.Info("Setting key in Redis: " + key)
	return nil
}
