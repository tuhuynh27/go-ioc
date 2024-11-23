package integration

import (
	"os"
	"testing"
	"time"
)

func TestComplexIoCIntegration(t *testing.T) {
	// Setup test environment
	tmpDir, err := os.MkdirTemp("", "ioc-integration-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize application
	app := Initialize()

	// Get UserService
	userService := app.GetUserService()
	if userService == nil {
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
	metricsCollector := app.GetMetricsCollector()
	if metricsCollector == nil {
		t.Fatal("Failed to get MetricsCollector")
	}

	cacheHits := metricsCollector.GetMetric("cache.hits")
	if cacheHits != 1 {
		t.Errorf("Expected 1 cache hit, got %f", cacheHits)
	}
}
