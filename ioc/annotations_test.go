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
