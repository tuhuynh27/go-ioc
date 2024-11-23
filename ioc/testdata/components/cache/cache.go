package cache

import (
	"sync"
	"time"

	"github.com/tuhuynh27/go-ioc/ioc"
	"github.com/tuhuynh27/go-ioc/ioc/testdata/components/logger"
	"github.com/tuhuynh27/go-ioc/ioc/testdata/components/metrics"
)

type Cache interface {
	Set(key string, value interface{}, expiration time.Duration)
	Get(key string) (interface{}, bool)
}

type InMemoryCache struct {
	ioc.Component `implements:"cache.Cache"`
	Logger        logger.Logger            `autowired:"" qualifier:"console"`
	Metrics       metrics.MetricsCollector `autowired:""`

	items map[string]cacheItem
	mu    sync.RWMutex
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

func (c *InMemoryCache) Set(key string, value interface{}, expiration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.items == nil {
		c.items = make(map[string]cacheItem)
	}

	c.items[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(expiration),
	}

	c.Logger.Log("Cache item set: " + key)
	c.Metrics.RecordMetric("cache.items.count", float64(len(c.items)))
}

func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.expiration) {
		c.Logger.LogWithLevel(logger.DEBUG, "Cache item expired: "+key)
		return nil, false
	}

	return item.value, true
}
