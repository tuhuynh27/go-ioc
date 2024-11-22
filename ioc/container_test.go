package ioc

import (
	"strings"
	"testing"
)

// Mock interfaces and structs for testing
type Logger interface {
	Log(message string)
}

type MessageService interface {
	Send(message string)
}

type StdoutLogger struct {
	Component
}

func (l *StdoutLogger) Log(message string) {}

type EmailService struct {
	Component
	Qualifier `value:"email"`
}

func (s *EmailService) Send(message string) {}

type SMSService struct {
	Component
	Qualifier `value:"sms"`
}

func (s *SMSService) Send(message string) {}

func TestContainer_Get(t *testing.T) {
	container := NewContainer()
	logger := &StdoutLogger{}

	// Test with full path
	fullPath := "github.com/example/logger.StdoutLogger"
	container.Register(fullPath, logger)

	tests := []struct {
		name        string
		lookup      string
		wantErr     bool
		errContains string
	}{
		{
			name:    "full path lookup",
			lookup:  fullPath,
			wantErr: false,
		},
		{
			name:    "short name lookup",
			lookup:  "StdoutLogger",
			wantErr: false,
		},
		{
			name:        "non-existent component",
			lookup:      "NonExistentLogger",
			wantErr:     true,
			errContains: "no component found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			component, err := container.Get(tt.lookup)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error %q should contain %q", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if component != logger {
				t.Error("retrieved component does not match registered component")
			}
		})
	}
}

func TestContainer_GetQualified(t *testing.T) {
	container := NewContainer()
	emailService := &EmailService{}
	smsService := &SMSService{}

	// Register with full paths
	container.RegisterWithInterface("github.com/example/message.MessageService", "email", emailService)
	container.RegisterWithInterface("github.com/example/message.MessageService", "sms", smsService)

	tests := []struct {
		name        string
		interface_  string
		qualifier   string
		wantErr     bool
		errContains string
		expected    MessageService // Changed to interface type
	}{
		{
			name:       "full path lookup",
			interface_: "github.com/example/message.MessageService",
			qualifier:  "email",
			wantErr:    false,
			expected:   emailService,
		},
		{
			name:       "short name lookup",
			interface_: "MessageService",
			qualifier:  "sms",
			wantErr:    false,
			expected:   smsService,
		},
		{
			name:        "non-existent interface",
			interface_:  "NonExistentService",
			qualifier:   "email",
			wantErr:     true,
			errContains: "no interface found",
		},
		{
			name:        "invalid qualifier",
			interface_:  "MessageService",
			qualifier:   "invalid",
			wantErr:     true,
			errContains: "no component found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			component, err := container.GetQualified(tt.interface_, tt.qualifier)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error %q should contain %q", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.wantErr && component != tt.expected {
				t.Error("retrieved component does not match expected component")
			}
		})
	}
}

func TestContainer_MultipleMatches(t *testing.T) {
	container := NewContainer()

	// Register components with same short name but different paths
	container.Register("github.com/example1/logger.StdoutLogger", &StdoutLogger{})
	container.Register("github.com/example2/logger.StdoutLogger", &StdoutLogger{})

	// Test Get with multiple matches
	t.Run("multiple component matches", func(t *testing.T) {
		_, err := container.Get("StdoutLogger")
		if err == nil {
			t.Error("expected error for multiple matches")
		}
		if !strings.Contains(err.Error(), "multiple components found") {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	// Register interfaces with same short name but different paths
	container.RegisterWithInterface("github.com/example1/message.MessageService", "email", &EmailService{})
	container.RegisterWithInterface("github.com/example2/message.MessageService", "email", &EmailService{})

	// Test GetQualified with multiple matches
	t.Run("multiple interface matches", func(t *testing.T) {
		_, err := container.GetQualified("MessageService", "email")
		if err == nil {
			t.Error("expected error for multiple matches")
		}
		if !strings.Contains(err.Error(), "multiple interfaces found") {
			t.Errorf("unexpected error message: %v", err)
		}
	})
}
