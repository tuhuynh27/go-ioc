package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseComponents(t *testing.T) {
	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "ioc-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create go.mod file
	goModContent := `module example.com/test
go 1.20
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	// Create test file structure
	err = createTestFiles(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create test files: %v", err)
	}

	// Parse components
	components, err := ParseComponents(tmpDir)
	if err != nil {
		t.Fatalf("ParseComponents failed: %v", err)
	}

	// Verify components
	if len(components) != 2 {
		t.Errorf("Expected 2 components, got %d", len(components))
	}

	// Update import paths to match the module
	expectedLoggerInterface := "example.com/test/logger.Logger"

	// Check logger component
	loggerFound := false
	for _, comp := range components {
		if comp.Type == "StdoutLogger" {
			loggerFound = true
			if len(comp.Implements) != 1 {
				t.Errorf("Expected 1 interface implementation, got %d", len(comp.Implements))
			}
			if comp.Implements[0] != expectedLoggerInterface {
				t.Errorf("Expected %s interface, got %s", expectedLoggerInterface, comp.Implements[0])
			}
			if comp.Package != "example.com/test/logger" {
				t.Errorf("Expected package example.com/test/logger, got %s", comp.Package)
			}
		}
	}
	if !loggerFound {
		t.Error("StdoutLogger component not found")
	}

	// Check service component
	serviceFound := false
	for _, comp := range components {
		if comp.Type == "EmailService" {
			serviceFound = true
			if len(comp.Dependencies) != 1 {
				t.Errorf("Expected 1 dependency, got %d", len(comp.Dependencies))
			}
			if comp.Dependencies[0].Type != "logger.Logger" {
				t.Errorf("Expected logger.Logger dependency, got %s", comp.Dependencies[0].Type)
			}
			if comp.Package != "example.com/test/message" {
				t.Errorf("Expected package example.com/test/message, got %s", comp.Package)
			}
		}
	}
	if !serviceFound {
		t.Error("EmailService component not found")
	}
}

func createTestFiles(dir string) error {
	// Create logger package
	loggerDir := filepath.Join(dir, "logger")
	if err := os.MkdirAll(loggerDir, 0755); err != nil {
		return err
	}

	// Create logger.go with fully qualified interface path
	loggerContent := `
package logger

import "example.com/test/ioc"

type Logger interface {
    Log(message string)
}

type StdoutLogger struct {
    ioc.Component ` + "`implements:\"example.com/test/logger.Logger\"`" + `
}

func (l *StdoutLogger) Log(message string) {
    println(message)
}
`
	if err := os.WriteFile(filepath.Join(loggerDir, "logger.go"), []byte(loggerContent), 0644); err != nil {
		return err
	}

	// Create message package
	messageDir := filepath.Join(dir, "message")
	if err := os.MkdirAll(messageDir, 0755); err != nil {
		return err
	}

	// Create email_service.go with fully qualified import
	emailContent := `
package message

import (
    "example.com/test/ioc"
    "example.com/test/logger"
)

type EmailService struct {
    ioc.Component
    Logger logger.Logger ` + "`autowired:\"\"`" + `
}

func (s *EmailService) SendEmail(to, message string) error {
    s.Logger.Log("Sending email: " + message)
    return nil
}
`
	if err := os.WriteFile(filepath.Join(messageDir, "email_service.go"), []byte(emailContent), 0644); err != nil {
		return err
	}

	// Create ioc package directory and file
	iocDir := filepath.Join(dir, "ioc")
	if err := os.MkdirAll(iocDir, 0755); err != nil {
		return err
	}

	// Create component.go in ioc package
	componentContent := `
package ioc

type Component struct{}
`
	if err := os.WriteFile(filepath.Join(iocDir, "component.go"), []byte(componentContent), 0644); err != nil {
		return err
	}

	return nil
}

func TestParseStructTag(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		expected map[string]string
	}{
		{
			name: "autowired tag",
			tag:  "`autowired:\"\" qualifier:\"email\"`",
			expected: map[string]string{
				"autowired": "",
				"qualifier": "email",
			},
		},
		{
			name: "implements tag",
			tag:  "`implements:\"logger.Logger\"`",
			expected: map[string]string{
				"implements": "logger.Logger",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseStructTag(tt.tag)
			for k, v := range tt.expected {
				if result[k] != v {
					t.Errorf("Expected %s=%s, got %s=%s", k, v, k, result[k])
				}
			}
		})
	}
}
