package integration

import (
	"os"
	"path/filepath"
	"testing"
)

// ensureTestStructure ensures all required test directories and files exist
func ensureTestStructure(t *testing.T) string {
	t.Helper()

	// Get the current test directory
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Create components directory if it doesn't exist
	componentsDir := filepath.Join(currentDir, "components")
	dirs := []string{
		filepath.Join(componentsDir, "logger"),
		filepath.Join(componentsDir, "service"),
		filepath.Join(componentsDir, "notification"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	return currentDir
}

func TestMain(m *testing.M) {
	// Run setup before tests
	code := m.Run()
	// Run cleanup after tests
	os.Exit(code)
}
