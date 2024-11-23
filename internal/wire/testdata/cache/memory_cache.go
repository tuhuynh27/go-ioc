package cache

import (
	"sync"

	"github.com/tuhuynh27/go-ioc/internal/wire/testdata/logger"
)

type MemoryCache struct {
	Component  struct{}
	Qualifier  struct{}      `value:"memory"`
	Implements struct{}      `implements:"Cache"`
	Logger     logger.Logger `autowired:"true" qualifier:"console"`

	mu    sync.RWMutex
	store map[string]string
}

func (c *MemoryCache) Get(key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	c.Logger.Info("Getting key from memory: " + key)
	return c.store[key], nil
}

func (c *MemoryCache) Set(key string, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Logger.Info("Setting key in memory: " + key)
	if c.store == nil {
		c.store = make(map[string]string)
	}
	c.store[key] = value
	return nil
}
