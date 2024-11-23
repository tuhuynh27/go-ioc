package metrics

import (
	"sync"

	"github.com/tuhuynh27/go-ioc/ioc"
	"github.com/tuhuynh27/go-ioc/ioc/testdata/components/logger"
)

type MetricsCollector interface {
	RecordMetric(name string, value float64)
	GetMetric(name string) float64
}

type InMemoryMetrics struct {
	ioc.Component `implements:"metrics.MetricsCollector"`
	Logger        logger.Logger `autowired:"" qualifier:"json"` // Using JSON logger for metrics
	metrics       map[string]float64
	mu            sync.RWMutex
}

func (m *InMemoryMetrics) RecordMetric(name string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.metrics == nil {
		m.metrics = make(map[string]float64)
	}

	m.metrics[name] = value
	m.Logger.LogWithLevel(logger.DEBUG, "Recorded metric: "+name)
}

func (m *InMemoryMetrics) GetMetric(name string) float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.metrics[name]
}
