package ioc

import (
	"reflect"
	"strings"
)

// Component is a marker struct used to identify injectable components
type Component struct{}

// Qualifier provides a way to disambiguate between multiple implementations
type Qualifier struct {
	Value string
}

// hasComponentAnnotation checks if a type has the Component annotation and returns
// the component name and true if found, empty string and false otherwise.
// The component name is either from the "name" tag or the lowercase type name.
func hasComponentAnnotation(t reflect.Type) (string, bool) {
	field, ok := t.FieldByName("Component")
	if !ok {
		return "", false
	}
	if name := field.Tag.Get("name"); name != "" {
		return strings.ToLower(name), true
	}
	return strings.ToLower(t.Name()), true
}

// getQualifier extracts the qualifier value from a type's Qualifier field tag.
// Returns empty string if no Qualifier field or value tag is found.
func getQualifier(t reflect.Type) string {
	field, ok := t.FieldByName("Qualifier")
	if !ok {
		return ""
	}
	return field.Tag.Get("value")
}

// getImplementedInterfaces scans a type's fields for "implements" tags and returns
// a slice of interface names that the type implements.
func getImplementedInterfaces(t reflect.Type) []string {
	var interfaces []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if impl := field.Tag.Get("implements"); impl != "" {
			interfaces = append(interfaces, impl)
		}
	}
	return interfaces
}

// findDependencies scans a type's fields for "autowired" tags and returns
// a slice of dependencies that need to be injected. Each dependency includes
// the field name, type name and optional qualifier.
func findDependencies(t reflect.Type) []dependency {
	var deps []dependency
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if _, ok := field.Tag.Lookup("autowired"); ok {
			dep := dependency{
				fieldName: field.Name,
				typeName:  field.Type.Name(),
				qualifier: field.Tag.Get("qualifier"),
			}
			deps = append(deps, dep)
		}
	}
	return deps
}
