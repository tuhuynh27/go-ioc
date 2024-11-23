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

	// Update expected component count
	expectedComponents := 6 // Updated for: EmailMessageService, SMSMessageService, NotificationService, MessageService interface, ConsoleLogger, JsonLogger
	if len(components) != expectedComponents {
		t.Errorf("Expected %d components, got %d", expectedComponents, len(components))
		for _, comp := range components {
			t.Logf("Found component: %s (package: %s)", comp.Type, comp.Package)
		}
	}

	// Verify components were found
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

	// Verify Message Services
	var emailService, smsService bool
	for _, comp := range components {
		switch comp.Type {
		case "EmailMessageService":
			emailService = true
			if comp.Qualifier != "email" {
				t.Errorf("Expected email qualifier for EmailMessageService, got %s", comp.Qualifier)
			}
			if !contains(comp.Implements, "MessageService") {
				t.Error("EmailMessageService should implement MessageService interface")
			}
			// Verify dependencies
			if len(comp.Dependencies) != 2 { // ConfigData and Logger
				t.Errorf("Expected 2 dependencies for EmailMessageService, got %d", len(comp.Dependencies))
			}
		case "SMSMessageService":
			smsService = true
			if comp.Qualifier != "sms" {
				t.Errorf("Expected sms qualifier for SMSMessageService, got %s", comp.Qualifier)
			}
			if !contains(comp.Implements, "MessageService") {
				t.Error("SMSMessageService should implement MessageService interface")
			}
			// Verify dependencies
			if len(comp.Dependencies) != 2 { // ConfigData and Logger
				t.Errorf("Expected 2 dependencies for SMSMessageService, got %d", len(comp.Dependencies))
			}
		}
	}

	if !emailService {
		t.Error("EmailMessageService component not found")
	}
	if !smsService {
		t.Error("SMSMessageService component not found")
	}

	// Verify NotificationService
	var notificationService bool
	for _, comp := range components {
		if comp.Type == "NotificationService" {
			notificationService = true
			if len(comp.Dependencies) != 3 { // EmailSender, SmsSender, and Logger
				t.Errorf("Expected 3 dependencies for NotificationService, got %d", len(comp.Dependencies))
			}

			// Verify specific dependencies
			var hasEmailSender, hasSmsSender, hasLogger bool
			for _, dep := range comp.Dependencies {
				switch {
				case dep.Type == "service.MessageService" && dep.Qualifier == "email":
					hasEmailSender = true
				case dep.Type == "service.MessageService" && dep.Qualifier == "sms":
					hasSmsSender = true
				case dep.Type == "logger.Logger" && dep.Qualifier == "json":
					hasLogger = true
				}
			}

			if !hasEmailSender {
				t.Error("NotificationService missing EmailMessageService dependency")
			}
			if !hasSmsSender {
				t.Error("NotificationService missing SMSMessageService dependency")
			}
			if !hasLogger {
				t.Error("NotificationService missing json logger dependency")
			}
		}
	}

	if !notificationService {
		t.Error("NotificationService component not found")
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
	// Create package directories
	dirs := []string{
		filepath.Join(dst, "logger"),
		filepath.Join(dst, "service"),
		filepath.Join(dst, "config"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Get absolute path of source directory
	srcAbs, err := filepath.Abs(".")
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Define files to copy for each package
	filesToCopy := map[string][]string{
		"logger": {"logger.go", "console_logger.go", "json_logger.go"},
		"service": {
			"message_service.go",
			"notification_service.go",
			"email_message_service.go",
			"sms_message_service.go",
		},
		"config": {"config.go"},
	}

	// Copy files
	for pkg, files := range filesToCopy {
		for _, file := range files {
			srcPath := filepath.Join(srcAbs, pkg, file)
			dstPath := filepath.Join(dst, pkg, file)
			if err := copyFile(srcPath, dstPath); err != nil {
				return fmt.Errorf("failed to copy %s file %s: %w", pkg, file, err)
			}
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
