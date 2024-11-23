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
			Qualifier: "stdout",
		},
		{
			Name:    "EmailService",
			Type:    "EmailService",
			Package: "github.com/tuhuynh27/go-ioc/examples/ioc-example-simple/message",
			Dependencies: []Dependency{
				{
					FieldName: "Logger",
					Type:      "logger.Logger",
					Qualifier: "stdout",
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

	// Check for imports
	expectedImports := []string{
		"github.com/tuhuynh27/go-ioc/examples/ioc-example-simple/logger",
		"github.com/tuhuynh27/go-ioc/examples/ioc-example-simple/message",
	}

	for _, imp := range expectedImports {
		if !contains(strings.Split(contentStr, "\n"), `"`+imp+`"`) {
			t.Errorf("Expected import %s not found", imp)
		}
	}

	// Check for interface registration
	expectedRegistrations := []string{
		`container.RegisterWithInterface("github.com/tuhuynh27/go-ioc/examples/ioc-example-simple/logger.Logger", "stdout", stdoutLogger)`,
	}

	for _, reg := range expectedRegistrations {
		if !strings.Contains(contentStr, reg) {
			t.Errorf("Expected interface registration not found: %s", reg)
		}
	}
}

func TestGenerator_GenerateComponentInits(t *testing.T) {
	components := []Component{
		{
			Name:       "StdoutLogger",
			Type:       "StdoutLogger",
			Package:    "logger",
			Implements: []string{"logger.Logger"},
			Qualifier:  "stdout",
		},
		{
			Name:    "EmailService",
			Type:    "EmailService",
			Package: "message",
			Dependencies: []Dependency{
				{
					FieldName: "Logger",
					Type:      "logger.Logger",
					Qualifier: "stdout",
				},
			},
		},
	}

	gen := NewGenerator(components)
	inits := gen.generateComponentInits(components)

	if len(inits) != 2 {
		t.Errorf("Expected 2 component inits, got %d", len(inits))
	}

	// Check first component (StdoutLogger)
	if inits[0].Type != "StdoutLogger" {
		t.Errorf("Expected StdoutLogger, got %s", inits[0].Type)
	}

	// Check interface registration for StdoutLogger
	if len(inits[0].Interfaces) != 1 {
		t.Errorf("Expected 1 interface registration for StdoutLogger, got %d", len(inits[0].Interfaces))
	}

	if inits[0].Interfaces[0].Interface != "logger.Logger" {
		t.Errorf("Expected logger.Logger interface, got %s", inits[0].Interfaces[0].Interface)
	}

	if inits[0].Interfaces[0].Qualifier != "stdout" {
		t.Errorf("Expected stdout qualifier, got %s", inits[0].Interfaces[0].Qualifier)
	}

	// Check second component (EmailService)
	if inits[1].Type != "EmailService" {
		t.Errorf("Expected EmailService, got %s", inits[1].Type)
	}

	// Check EmailService dependencies
	if len(inits[1].Dependencies) != 1 {
		t.Errorf("Expected 1 dependency for EmailService, got %d", len(inits[1].Dependencies))
	}

	if inits[1].Dependencies[0].FieldName != "Logger" {
		t.Errorf("Expected Logger field name, got %s", inits[1].Dependencies[0].FieldName)
	}
}
