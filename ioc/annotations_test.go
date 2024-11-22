package ioc

import (
	"reflect"
	"testing"
)

func TestComponent_Marker(t *testing.T) {
	// Test that Component can be embedded in a struct
	type TestComponent struct {
		Component
	}

	comp := TestComponent{}

	// Verify that the struct has Component embedded
	t.Run("component embedding", func(t *testing.T) {
		typ := reflect.TypeOf(comp)
		field, exists := typ.FieldByName("Component")
		if !exists {
			t.Error("Component field not found in struct")
		}
		if field.Type != reflect.TypeOf(Component{}) {
			t.Error("Embedded field is not of type Component")
		}
	})
}

func TestQualifier_Value(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"empty qualifier", "", ""},
		{"email qualifier", "email", "email"},
		{"sms qualifier", "sms", "sms"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := Qualifier{Value: tt.value}
			if q.Value != tt.expected {
				t.Errorf("Qualifier.Value = %v, want %v", q.Value, tt.expected)
			}
		})
	}
}

func TestQualifier_Embedding(t *testing.T) {
	// Test that Qualifier can be embedded in a struct with a value
	type TestService struct {
		Component
		Qualifier `value:"test"`
	}

	service := TestService{}

	t.Run("qualifier embedding", func(t *testing.T) {
		typ := reflect.TypeOf(service)
		field, exists := typ.FieldByName("Qualifier")
		if !exists {
			t.Error("Qualifier field not found in struct")
		}
		if field.Type != reflect.TypeOf(Qualifier{}) {
			t.Error("Embedded field is not of type Qualifier")
		}

		// Check if the struct tag is present
		tag := field.Tag.Get("value")
		if tag != "test" {
			t.Errorf("Qualifier tag value = %v, want %v", tag, "test")
		}
	})
}

func TestComponent_And_Qualifier_Together(t *testing.T) {
	// Test both Component and Qualifier together
	type TestService struct {
		Component
		Qualifier `value:"test"`
		Name      string
	}

	service := TestService{Name: "test-service"}

	t.Run("combined annotations", func(t *testing.T) {
		typ := reflect.TypeOf(service)

		// Check Component
		_, hasComponent := typ.FieldByName("Component")
		if !hasComponent {
			t.Error("Component field not found")
		}

		// Check Qualifier
		qualifierField, hasQualifier := typ.FieldByName("Qualifier")
		if !hasQualifier {
			t.Error("Qualifier field not found")
		}

		// Check Qualifier tag
		tag := qualifierField.Tag.Get("value")
		if tag != "test" {
			t.Errorf("Qualifier tag = %v, want %v", tag, "test")
		}

		// Check that regular fields are still accessible
		if service.Name != "test-service" {
			t.Errorf("Name = %v, want %v", service.Name, "test-service")
		}
	})
}
