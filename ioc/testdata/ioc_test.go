package integration

import (
	"os"
	"testing"
	"time"

	"github.com/tuhuynh27/go-ioc/ioc"
	"github.com/tuhuynh27/go-ioc/ioc/testdata/components/cache"
	"github.com/tuhuynh27/go-ioc/ioc/testdata/components/logger"
	"github.com/tuhuynh27/go-ioc/ioc/testdata/components/metrics"
	"github.com/tuhuynh27/go-ioc/ioc/testdata/components/service"
)

func TestComplexIoCIntegration(t *testing.T) {
	// Setup test environment
	tmpDir, err := os.MkdirTemp("", "ioc-integration-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize container
	container := setupTestContainer(t)

	// Get UserService
	userService, ok := container.Get("service.UserService").(*service.UserService)
	if !ok {
		t.Fatal("Failed to get UserService")
	}

	// Test user creation and retrieval
	testUser := service.User{
		ID:    "123",
		Name:  "Test User",
		Email: "test@example.com",
	}

	// Create user
	err = userService.CreateUser(testUser)
	if err != nil {
		t.Errorf("Failed to create user: %v", err)
	}

	// Wait for async operations
	time.Sleep(100 * time.Millisecond)

	// Retrieve user
	retrievedUser, err := userService.GetUser("123")
	if err != nil {
		t.Errorf("Failed to get user: %v", err)
	}

	if retrievedUser.ID != testUser.ID || retrievedUser.Name != testUser.Name {
		t.Errorf("Retrieved user doesn't match created user")
	}

	// Test metrics
	metricsCollector, ok := container.Get("metrics.MetricsCollector").(metrics.MetricsCollector)
	if !ok {
		t.Fatal("Failed to get MetricsCollector")
	}

	cacheHits := metricsCollector.GetMetric("cache.hits")
	if cacheHits != 1 {
		t.Errorf("Expected 1 cache hit, got %f", cacheHits)
	}
}

func setupTestContainer(t *testing.T) *ioc.Container {
	container := ioc.NewContainer()

	// Initialize loggers
	consoleLogger := &logger.ConsoleLogger{}
	jsonLogger := &logger.JsonLogger{}
	container.Register("logger.ConsoleLogger", consoleLogger)
	container.Register("logger.JsonLogger", jsonLogger)
	container.RegisterWithInterface("logger.Logger", "console", consoleLogger)
	container.RegisterWithInterface("logger.Logger", "json", jsonLogger)

	// Initialize metrics
	metricsCollector := &metrics.InMemoryMetrics{
		Logger: jsonLogger,
	}
	container.Register("metrics.MetricsCollector", metricsCollector)
	container.RegisterWithInterface("metrics.MetricsCollector", "", metricsCollector)

	// Initialize cache
	cacheService := &cache.InMemoryCache{
		Logger:  consoleLogger,
		Metrics: metricsCollector,
	}
	container.Register("cache.Cache", cacheService)
	container.RegisterWithInterface("cache.Cache", "", cacheService)

	// Initialize message services
	emailService := &service.EmailService{
		Logger: jsonLogger,
	}
	smsService := &service.SmsService{
		Logger: jsonLogger,
	}
	container.Register("service.EmailService", emailService)
	container.Register("service.SmsService", smsService)
	container.RegisterWithInterface("service.MessageService", "email", emailService)
	container.RegisterWithInterface("service.MessageService", "sms", smsService)

	// Initialize user service
	userService := &service.UserService{
		Logger:      jsonLogger,
		Cache:       cacheService,
		Metrics:     metricsCollector,
		EmailSender: emailService,
		SmsSender:   smsService,
	}
	container.Register("service.UserService", userService)

	return container
}
