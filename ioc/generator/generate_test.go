package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerator_Generate(t *testing.T) {
	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "ioc-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test components
	components := []Component{
		{
			Name:    "StdoutLogger",
			Type:    "StdoutLogger",
			Package: "github.com/tuhuynh27/go-ioc/examples/ioc-example-simple/logger",
			Implements: []string{
				"github.com/tuhuynh27/go-ioc/examples/ioc-example-simple/logger.Logger",
			},
		},
		{
			Name:    "EmailService",
			Type:    "EmailService",
			Package: "github.com/tuhuynh27/go-ioc/examples/ioc-example-simple/message",
			Dependencies: []Dependency{
				{
					FieldName: "Logger",
					Type:      "logger.Logger",
				},
			},
		},
		{
			Name:    "NotificationService",
			Type:    "NotificationService",
			Package: "github.com/tuhuynh27/go-ioc/examples/ioc-example-simple/notification",
			Dependencies: []Dependency{
				{
					FieldName: "EmailSender",
					Type:      "message.MessageService",
					Qualifier: "email",
				},
				{
					FieldName: "SmsSender",
					Type:      "message.MessageService",
					Qualifier: "sms",
				},
			},
		},
	}

	// Create generator
	gen := NewGenerator(components)

	// Generate code
	err = gen.Generate(tmpDir)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Check if wire directory was created
	wireDir := filepath.Join(tmpDir, "wire")
	if _, err := os.Stat(wireDir); os.IsNotExist(err) {
		t.Error("wire directory was not created")
	}

	// Check if wire_gen.go was created
	wirePath := filepath.Join(wireDir, "wire_gen.go")
	if _, err := os.Stat(wirePath); os.IsNotExist(err) {
		t.Error("wire_gen.go was not created")
	}

	// Read generated file
	content, err := os.ReadFile(wirePath)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	// Check content
	contentStr := string(content)
	expectedImports := []string{
		"github.com/tuhuynh27/go-ioc/examples/ioc-example-simple/logger",
		"github.com/tuhuynh27/go-ioc/examples/ioc-example-simple/message",
		"github.com/tuhuynh27/go-ioc/examples/ioc-example-simple/notification",
	}

	for _, imp := range expectedImports {
		if !contains(strings.Split(contentStr, "\n"), `"`+imp+`"`) {
			t.Errorf("Expected import %s not found", imp)
		}
	}
}

func TestGenerator_GenerateComponentInits(t *testing.T) {
	components := []Component{
		{
			Name:    "StdoutLogger",
			Type:    "StdoutLogger",
			Package: "logger",
		},
		{
			Name:    "EmailService",
			Type:    "EmailService",
			Package: "message",
			Dependencies: []Dependency{
				{
					FieldName: "Logger",
					Type:      "logger.Logger",
				},
			},
		},
	}

	gen := NewGenerator(components)
	inits := gen.generateComponentInits(components)

	if len(inits) != 2 {
		t.Errorf("Expected 2 component inits, got %d", len(inits))
	}

	// Check first component
	if inits[0].Type != "StdoutLogger" {
		t.Errorf("Expected StdoutLogger, got %s", inits[0].Type)
	}

	// Check second component
	if inits[1].Type != "EmailService" {
		t.Errorf("Expected EmailService, got %s", inits[1].Type)
	}
}
