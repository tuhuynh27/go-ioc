package wire

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
		if !strings.Contains(contentStr, `"`+imp+`"`) {
			t.Errorf("Expected import %s not found", imp)
		}
	}

	// Check for Container struct
	if !strings.Contains(contentStr, "type Container struct {") {
		t.Error("Container struct not found in generated code")
	}

	// Check for Initialize function
	if !strings.Contains(contentStr, "func Initialize() (*Container, func()) {") {
		t.Error("Initialize function not found in generated code")
	}

	// Check for exported fields
	expectedFields := []string{
		"StdoutLogger *logger.StdoutLogger",
		"EmailService *message.EmailService",
	}

	for _, field := range expectedFields {
		if !strings.Contains(contentStr, field) {
			t.Errorf("Expected field not found: %s", field)
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

func TestGenerator_GenerateWithLifecycleMethods(t *testing.T) {
	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "ioc-test-lifecycle-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test components with lifecycle methods
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
			Name:          "LifecycleService",
			Type:          "LifecycleService",
			Package:       "github.com/tuhuynh27/go-ioc/examples/ioc-example-simple/lifecycle",
			PostConstruct: true,
			PreDestroy:    true,
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
		"github.com/tuhuynh27/go-ioc/examples/ioc-example-simple/lifecycle",
		"github.com/tuhuynh27/go-ioc/examples/ioc-example-simple/message",
		"github.com/tuhuynh27/go-ioc/examples/ioc-example-simple/logger",
	}

	for _, imp := range expectedImports {
		if !strings.Contains(contentStr, `"`+imp+`"`) {
			t.Errorf("Expected import %s not found", imp)
		}
	}

	// Check for Container struct
	if !strings.Contains(contentStr, "type Container struct {") {
		t.Error("Container struct not found in generated code")
	}

	// Check for Initialize function with updated signature
	if !strings.Contains(contentStr, "func Initialize() (*Container, func()) {") {
		t.Error("Initialize function not found in generated code")
	}

	// Check for exported fields
	expectedFields := []string{
		"LifecycleService *lifecycle.LifecycleService",
		"EmailService *message.EmailService",
		"StdoutLogger *logger.StdoutLogger",
	}

	for _, field := range expectedFields {
		if !strings.Contains(contentStr, field) {
			t.Errorf("Expected field not found: %s", field)
		}
	}

	// Check for PostConstruct call
	if !strings.Contains(contentStr, "container.LifecycleService.PostConstruct()") {
		t.Error("PostConstruct method call not found in generated code")
	}

	// Check for PreDestroy call
	if !strings.Contains(contentStr, "container.LifecycleService.PreDestroy()") {
		t.Error("PreDestroy method call not found in generated code")
	}
}

func TestGenerator_GenerateCyclicDependencies(t *testing.T) {
	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "ioc-test-cyclic-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test components with cyclic dependencies
	components := []Component{
		{
			Name:    "ServiceA",
			Type:    "ServiceA",
			Package: "github.com/tuhuynh27/go-ioc/examples/ioc-example-simple/service",
			Dependencies: []Dependency{
				{
					FieldName: "B",
					Type:      "ServiceB",
				},
			},
		},
		{
			Name:    "ServiceB",
			Type:    "ServiceB",
			Package: "github.com/tuhuynh27/go-ioc/examples/ioc-example-simple/service",
			Dependencies: []Dependency{
				{
					FieldName: "A",
					Type:      "ServiceA",
				},
			},
		},
	}

	// Create generator
	gen := NewGenerator(components)

	// Expect panic due to cyclic dependencies
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic due to cyclic dependencies, but did not panic")
		}
	}()

	// Attempt to generate code
	gen.Generate(tmpDir) // This should trigger the panic
}
