package ioc

import (
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
	Qualifier struct{} `value:"email"`
}

func (s *EmailService) Send(message string) {}

type SMSService struct {
	Component
	Qualifier struct{} `value:"sms"`
}

func (s *SMSService) Send(message string) {}

func TestContainer_Get(t *testing.T) {
	container := NewContainer()
	logger := &StdoutLogger{}

	// Test with full path
	fullPath := "github.com/example/logger.StdoutLogger"
	container.Register(fullPath, logger)

	tests := []struct {
		name    string
		lookup  string
		wantNil bool
	}{
		{
			name:    "full path lookup",
			lookup:  fullPath,
			wantNil: false,
		},
		{
			name:    "short name lookup",
			lookup:  "StdoutLogger",
			wantNil: false,
		},
		{
			name:    "non-existent component",
			lookup:  "NonExistentLogger",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			component := container.Get(tt.lookup)

			if tt.wantNil {
				if component != nil {
					t.Error("expected nil component, got non-nil")
				}
				return
			}

			if component == nil {
				t.Error("unexpected nil component")
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
		name       string
		interface_ string
		qualifier  string
		wantNil    bool
		expected   MessageService // Changed to interface type
	}{
		{
			name:       "full path lookup",
			interface_: "github.com/example/message.MessageService",
			qualifier:  "email",
			wantNil:    false,
			expected:   emailService,
		},
		{
			name:       "short name lookup",
			interface_: "MessageService",
			qualifier:  "sms",
			wantNil:    false,
			expected:   smsService,
		},
		{
			name:       "non-existent interface",
			interface_: "NonExistentService",
			qualifier:  "email",
			wantNil:    true,
		},
		{
			name:       "invalid qualifier",
			interface_: "MessageService",
			qualifier:  "invalid",
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			component := container.GetQualified(tt.interface_, tt.qualifier)

			if tt.wantNil {
				if component != nil {
					t.Error("expected nil component, got non-nil")
				}
				return
			}

			if component == nil {
				t.Error("unexpected nil component")
			}

			if !tt.wantNil && component != tt.expected {
				t.Error("retrieved component does not match expected component")
			}
		})
	}
}
