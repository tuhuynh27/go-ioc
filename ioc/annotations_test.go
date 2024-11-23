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

func TestQualifier_Tag(t *testing.T) {
	// Test that Qualifier can be used as a field with struct tag
	type TestService struct {
		Component
		Qualifier struct{} `value:"test"`
	}

	service := TestService{}

	t.Run("qualifier tag", func(t *testing.T) {
		typ := reflect.TypeOf(service)
		field, exists := typ.FieldByName("Qualifier")
		if !exists {
			t.Error("Qualifier field not found in struct")
		}
		if field.Type != reflect.TypeOf(struct{}{}) {
			t.Error("Qualifier field is not of type struct{}")
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
		Qualifier  struct{} `value:"test"`
		Implements struct{} `implements:"Service"`
		Name       string
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

		// Check Implements tag
		implementsField, hasImplements := typ.FieldByName("Implements")
		if !hasImplements {
			t.Error("Implements field not found")
		}
		implementsTag := implementsField.Tag.Get("implements")
		if implementsTag != "Service" {
			t.Errorf("Implements tag = %v, want %v", implementsTag, "Service")
		}

		// Check that regular fields are still accessible
		if service.Name != "test-service" {
			t.Errorf("Name = %v, want %v", service.Name, "test-service")
		}
	})
}
