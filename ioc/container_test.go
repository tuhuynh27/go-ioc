package ioc

import (
	"testing"
)

type ComponentA struct {
	Component Component
}

type ComponentB struct {
	Component Component
	A         ComponentA `autowired:"true"`
}

func TestNewContainer(t *testing.T) {
	container := NewContainer()
	if container == nil {
		t.Fatal("Expected a new container instance, got nil")
	}
}

func TestRegisterComponent(t *testing.T) {
	container := NewContainer()

	err := container.RegisterComponents(
		&ComponentA{}, &ComponentB{},
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if container.Get("ComponentA") == nil {
		t.Fatal("Expected component to be registered")
	}
}

// Tests for multiple implementations of the same interface but no qualifier specified
type TestInterface interface {
	TestMethod()
}

type ComponentC struct {
	Component
	Qualifier  struct{} `value:"c"`
	Implements struct{} `implements:"TestInterface"`
}

type ComponentD struct {
	Component
	Qualifier  struct{} `value:"d"`
	Implements struct{} `implements:"TestInterface"`
}

func (c *ComponentC) TestMethod() {}
func (d *ComponentD) TestMethod() {}

type ComponentE struct {
	Component
	Dep TestInterface `autowired:"true"`
}

func TestMultipleImplementationsNoQualifier(t *testing.T) {
	container := NewContainer()

	err := container.RegisterComponents(
		&ComponentC{}, &ComponentD{}, &ComponentE{},
	)

	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
}
