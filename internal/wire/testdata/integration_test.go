package testdata

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tuhuynh27/go-ioc/internal/wire"
)

func TestIntegration(t *testing.T) {
	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "ioc-integration-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create go.mod file
	goModContent := `module github.com/tuhuynh27/go-ioc/internal/wire/testdata
go 1.20

require github.com/tuhuynh27/go-ioc v0.0.0
replace github.com/tuhuynh27/go-ioc => ../../../..
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	// Copy test files to temp directory
	testdataDir := "testdata"
	if err := copyTestFiles(testdataDir, tmpDir); err != nil {
		t.Fatalf("Failed to copy test files: %v", err)
	}

	// Parse components
	components, err := wire.ParseComponents(tmpDir)
	if err != nil {
		t.Fatalf("Failed to parse components: %v", err)
	}

	// Verify components were found
	expectedComponents := 10
	if len(components) != expectedComponents {
		t.Errorf("Expected %d components, got %d", expectedComponents, len(components))
		// Print found components for debugging
		for _, comp := range components {
			t.Logf("Found component: %s (package: %s)", comp.Type, comp.Package)
		}
	}

	// Verify loggers
	var consoleLogger, jsonLogger bool
	for _, comp := range components {
		switch comp.Type {
		case "ConsoleLogger":
			consoleLogger = true
			if comp.Qualifier != "console" {
				t.Errorf("Expected console qualifier for ConsoleLogger, got %s", comp.Qualifier)
			}
			if !contains(comp.Implements, "Logger") {
				t.Error("ConsoleLogger should implement Logger interface")
			}
		case "JsonLogger":
			jsonLogger = true
			if comp.Qualifier != "json" {
				t.Errorf("Expected json qualifier for JsonLogger, got %s", comp.Qualifier)
			}
			if !contains(comp.Implements, "Logger") {
				t.Error("JsonLogger should implement Logger interface")
			}
		}
	}

	if !consoleLogger {
		t.Error("ConsoleLogger component not found")
	}
	if !jsonLogger {
		t.Error("JsonLogger component not found")
	}

	// Verify NotificationService dependencies
	var notificationService bool
	for _, comp := range components {
		if comp.Type == "NotificationService" {
			notificationService = true
			if len(comp.Dependencies) != 3 {
				t.Errorf("Expected 3 dependencies for NotificationService, got %d", len(comp.Dependencies))
			}

			var hasConsole, hasJson, hasConfig bool
			for _, dep := range comp.Dependencies {
				switch {
				case dep.Type == "logger.Logger" && dep.Qualifier == "console":
					hasConsole = true
				case dep.Type == "logger.Logger" && dep.Qualifier == "json":
					hasJson = true
				case dep.Type == "config.NotificationConfigData":
					hasConfig = true
				}
			}

			if !hasConsole {
				t.Error("NotificationService missing console logger dependency")
			}
			if !hasJson {
				t.Error("NotificationService missing json logger dependency")
			}
			if !hasConfig {
				t.Error("NotificationService missing config dependency")
			}
		}
	}

	if !notificationService {
		t.Error("NotificationService component not found")
	}

	// Verify UserService nested dependency
	var userService bool
	for _, comp := range components {
		if comp.Type == "UserService" {
			userService = true
			if len(comp.Dependencies) != 1 {
				t.Errorf("Expected 1 dependency for UserService, got %d", len(comp.Dependencies))
			}
			if comp.Dependencies[0].Type != "NotificationService" {
				t.Errorf("Expected NotificationService dependency, got %s", comp.Dependencies[0].Type)
			}
		}
	}

	if !userService {
		t.Error("UserService component not found")
	}

	// Verify config components
	var dbConfig, notifyConfig bool
	for _, comp := range components {
		switch comp.Type {
		case "DatabaseConfigData":
			dbConfig = true
		case "NotificationConfigData":
			notifyConfig = true
		}
	}

	if !dbConfig {
		t.Error("DatabaseConfigData component not found")
	}
	if !notifyConfig {
		t.Error("NotificationConfigData component not found")
	}

	// Verify components that depend on configs
	for _, comp := range components {
		switch comp.Type {
		case "UserRepository":
			found := false
			for _, dep := range comp.Dependencies {
				if dep.Type == "config.DatabaseConfigData" {
					found = true
					break
				}
			}
			if !found {
				t.Error("UserRepository missing DatabaseConfigData dependency")
			}
		case "NotificationService":
			found := false
			for _, dep := range comp.Dependencies {
				if dep.Type == "config.NotificationConfigData" {
					found = true
					break
				}
			}
			if !found {
				t.Error("NotificationService missing NotificationConfigData dependency")
			}
		}
	}

	// Generate wire code
	gen := wire.NewGenerator(components)
	if err := gen.Generate(tmpDir); err != nil {
		t.Fatalf("Failed to generate wire code: %v", err)
	}

	// Verify wire file was generated
	wireFile := filepath.Join(tmpDir, "wire", "wire_gen.go")
	if _, err := os.Stat(wireFile); os.IsNotExist(err) {
		t.Error("wire_gen.go was not generated")
	}
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if strings.HasSuffix(s, str) {
			return true
		}
	}
	return false
}

func copyTestFiles(src, dst string) error {
	// Create logger package directory
	loggerDir := filepath.Join(dst, "logger")
	if err := os.MkdirAll(loggerDir, 0755); err != nil {
		return fmt.Errorf("failed to create logger directory: %w", err)
	}

	// Create service package directory
	serviceDir := filepath.Join(dst, "service")
	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		return fmt.Errorf("failed to create service directory: %w", err)
	}

	// Create additional package directories
	configDir := filepath.Join(dst, "config")
	cacheDir := filepath.Join(dst, "cache")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Create repository package directory
	repoDir := filepath.Join(dst, "repository")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		return fmt.Errorf("failed to create repository directory: %w", err)
	}

	// Get the absolute path of the source directory
	srcAbs, err := filepath.Abs(".")
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Copy logger package files
	loggerFiles := []string{"logger.go", "console_logger.go", "json_logger.go"}
	for _, file := range loggerFiles {
		srcPath := filepath.Join(srcAbs, "logger", file)
		dstPath := filepath.Join(loggerDir, file)
		if err := copyFile(srcPath, dstPath); err != nil {
			return fmt.Errorf("failed to copy logger file %s: %w", file, err)
		}
	}

	// Copy service package files
	serviceFiles := []string{"notification_service.go", "user_service.go", "user_service_enhanced.go"}
	for _, file := range serviceFiles {
		srcPath := filepath.Join(srcAbs, "service", file)
		dstPath := filepath.Join(serviceDir, file)
		if err := copyFile(srcPath, dstPath); err != nil {
			return fmt.Errorf("failed to copy service file %s: %w", file, err)
		}
	}

	// Copy config files
	configFiles := []string{"database_config.go", "notification_config.go"}
	for _, file := range configFiles {
		srcPath := filepath.Join(srcAbs, "config", file)
		dstPath := filepath.Join(configDir, file)
		if err := copyFile(srcPath, dstPath); err != nil {
			return fmt.Errorf("failed to copy config file %s: %w", file, err)
		}
	}

	// Copy cache files
	cacheFiles := []string{"redis_cache.go", "memory_cache.go"}
	for _, file := range cacheFiles {
		srcPath := filepath.Join(srcAbs, "cache", file)
		dstPath := filepath.Join(cacheDir, file)
		if err := copyFile(srcPath, dstPath); err != nil {
			return fmt.Errorf("failed to copy cache file %s: %w", file, err)
		}
	}

	// Copy repository files
	repoFiles := []string{"user_repository.go"}
	for _, file := range repoFiles {
		srcPath := filepath.Join(srcAbs, "repository", file)
		dstPath := filepath.Join(repoDir, file)
		if err := copyFile(srcPath, dstPath); err != nil {
			return fmt.Errorf("failed to copy repository file %s: %w", file, err)
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	// Read source file
	content, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	// Write to destination file
	if err := os.WriteFile(dst, content, 0644); err != nil {
		return fmt.Errorf("failed to write destination file: %w", err)
	}

	return nil
}
