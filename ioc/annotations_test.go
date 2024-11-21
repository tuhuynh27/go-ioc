package ioc

import (
	"reflect"
	"testing"
)

type TestStruct struct {
	Component `name:"testComponent"`
	Qualifier `value:"testValue"`
	Field1    string `autowired:"true" qualifier:"testQualifier"`
	Field2    int    `implements:"TestInterface"`
}

func TestHasComponentAnnotation(t *testing.T) {
	tType := reflect.TypeOf(TestStruct{})
	name, ok := hasComponentAnnotation(tType)
	if !ok || name != "testcomponent" {
		t.Errorf("Expected component name 'testcomponent', got '%s'", name)
	}
}

func TestGetQualifier(t *testing.T) {
	tType := reflect.TypeOf(TestStruct{})
	value := getQualifier(tType)
	if value != "testValue" {
		t.Errorf("Expected qualifier value 'testValue', got '%s'", value)
	}
}

func TestGetImplementedInterfaces(t *testing.T) {
	tType := reflect.TypeOf(TestStruct{})
	interfaces := getImplementedInterfaces(tType)
	expected := []string{"TestInterface"}
	if len(interfaces) != 1 || interfaces[0] != expected[0] {
		t.Errorf("Expected interfaces %v, got %v", expected, interfaces)
	}
}

func TestFindDependencies(t *testing.T) {
	tType := reflect.TypeOf(TestStruct{})
	deps := findDependencies(tType)
	if len(deps) != 1 {
		t.Errorf("Expected 1 dependency, got %d", len(deps))
	}
	if deps[0].fieldName != "Field1" || deps[0].typeName != "string" || deps[0].qualifier != "testQualifier" {
		t.Errorf("Dependency mismatch: %+v", deps[0])
	}
}
