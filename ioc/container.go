// Package ioc provides a dependency injection container implementation for Go applications.
package ioc

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"
)

// Container manages component registration, dependency resolution and injection.
// It maintains maps to track components, their definitions, interface implementations,
// and visited status during dependency resolution.
type Container struct {
	components  map[string]interface{}            // Stores instantiated components by name
	definitions map[string]*componentDefinition   // Stores component definitions by name
	interfaces  map[string]map[string]interface{} // Maps interface names to qualified implementations
	visited     map[string]bool                   // Tracks visited components during dependency resolution
}

// componentDefinition stores metadata about a registered component including its type,
// instance, dependencies, implemented interfaces and qualifier.
type componentDefinition struct {
	typ          reflect.Type // The concrete type of the component
	dependencies []dependency // List of dependencies to be injected
	implements   []string     // List of interface names this component implements
	qualifier    string       // Qualifier value for disambiguation
}

// dependency represents a single dependency to be injected into a component.
type dependency struct {
	fieldName string // Name of the field to inject into
	typeName  string // Type name of the dependency (full package path)
	qualifier string // Optional qualifier to specify which implementation to inject
}

// getFullComponentName returns the full component name using package path and type name
func getFullComponentName(t reflect.Type, name string) string {
	return strings.ToLower(t.PkgPath() + "." + name)
}

// getFullInterfaceName returns the full interface name using package path
// If no package path is provided, uses the current package path
func getFullInterfaceName(t reflect.Type, interfaceName string) string {
	if strings.Contains(interfaceName, ".") {
		return interfaceName // Already has package path
	}
	return t.PkgPath() + "." + interfaceName
}

// NewContainer creates and initializes a new dependency injection container.
// It initializes all required internal maps for tracking components and their relationships.
func NewContainer() *Container {
	return &Container{
		components:  make(map[string]interface{}),
		definitions: make(map[string]*componentDefinition),
		interfaces:  make(map[string]map[string]interface{}),
		visited:     make(map[string]bool),
	}
}

// RegisterComponents registers multiple components with the container and resolves their dependencies.
// It processes each component through registration phase and then initiates the dependency resolution process.
// Returns an error if registration or resolution fails for any component.
func (c *Container) RegisterComponents(components ...interface{}) error {
	startTime := time.Now()
	log.Printf("Starting IoC container initialization with %d components...", len(components))

	for _, component := range components {
		if err := c.register(component); err != nil {
			return fmt.Errorf("failed to register component %T: %w", component, err)
		}
	}

	err := c.resolve()
	if err != nil {
		return fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	duration := time.Since(startTime)
	log.Printf("IoC container initialization completed in %v", duration)

	return nil
}

// register processes a single component for registration in the container.
// It extracts component metadata, creates a definition, and registers interface implementations.
// Returns an error if the component is not properly annotated or registration fails.
func (c *Container) register(component interface{}) error {
	t := reflect.TypeOf(component)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	name, ok := hasComponentAnnotation(t)
	if !ok {
		return fmt.Errorf("type %s is not marked as @Component - ensure Component struct field is present", t.Name())
	}

	fullName := getFullComponentName(t, name)
	qualifier := getQualifier(t)
	implements := getImplementedInterfaces(t)

	// Convert interface names to full paths
	fullImplements := make([]string, len(implements))
	for i, iface := range implements {
		fullImplements[i] = getFullInterfaceName(t, iface)
	}

	def := &componentDefinition{
		typ:          t,
		dependencies: findDependencies(t),
		implements:   fullImplements,
		qualifier:    qualifier,
	}
	c.definitions[fullName] = def

	// register interface implementations
	for _, _interface := range fullImplements {
		if c.interfaces[_interface] == nil {
			c.interfaces[_interface] = make(map[string]interface{})
		}
	}

	return nil
}

// resolve performs the two-pass dependency resolution process.
// First pass creates and registers all component instances.
// Second pass resolves and injects dependencies for all components.
// Returns an error if resolution fails for any component.
func (c *Container) resolve() error {
	// Reset all tracking maps
	c.visited = make(map[string]bool)
	c.components = make(map[string]interface{})
	c.interfaces = make(map[string]map[string]interface{})

	// First pass: register all interface implementations
	for name, def := range c.definitions {
		instance := reflect.New(def.typ).Interface()
		c.components[name] = instance

		// register interfaces immediately
		for _, _interface := range def.implements {
			if c.interfaces[_interface] == nil {
				c.interfaces[_interface] = make(map[string]interface{})
			}
			c.interfaces[_interface][def.qualifier] = instance
		}
	}

	// Second pass: Resolve dependencies
	for name := range c.definitions {
		if !c.visited[name] {
			if err := c.dfsResolve(name); err != nil {
				return fmt.Errorf("failed to resolve component '%s': %w", name, err)
			}
		}
	}
	return nil
}

// dfsResolve performs depth-first resolution of a component's dependencies.
// It ensures all dependencies are resolved before injecting them into the component.
// Returns an error if resolution fails for any dependency.
func (c *Container) dfsResolve(name string) error {
	// Check if this component has already been resolved to avoid cycles
	if c.visited[name] {
		return nil
	}

	// Look up the component definition by name
	def := c.definitions[name]
	if def == nil {
		return fmt.Errorf("component definition not found for '%s' - ensure component is properly registered", name)
	}

	// Get the component instance that was created in the first pass
	instance := c.components[name]

	// First resolve all dependencies recursively using DFS
	// This ensures dependencies are resolved before the current component
	for _, dep := range def.dependencies {
		field, ok := def.typ.FieldByName(dep.fieldName)
		if !ok {
			return fmt.Errorf("field '%s' not found in component type", dep.fieldName)
		}

		depType := field.Type
		if depType.Kind() == reflect.Ptr {
			depType = depType.Elem()
		}

		// Create full dependency name using package path
		depName := getFullComponentName(depType, dep.typeName)

		// Only resolve if dependency is a registered component
		if _, exists := c.definitions[depName]; exists {
			// Recursively resolve the dependency
			if err := c.dfsResolve(depName); err != nil {
				return fmt.Errorf("failed to resolve dependency '%s' for component '%s': %w", depName, name, err)
			}
		}
	}

	// After all dependencies are resolved, inject them into the current component
	if err := c.injectDependencies(instance, def.dependencies); err != nil {
		return fmt.Errorf("failed to inject dependencies for component '%s': %w", name, err)
	}

	// Mark this component as resolved to prevent cycles and duplicate resolution
	c.visited[name] = true
	return nil
}

// injectDependencies injects all dependencies into a component instance.
// It handles both interface and direct component dependencies, with support for qualifiers.
// Returns an error if any dependency cannot be resolved or injected.
func (c *Container) injectDependencies(instance interface{}, deps []dependency) error {
	val := reflect.ValueOf(instance).Elem()

	for _, dep := range deps {
		field := val.FieldByName(dep.fieldName)
		if !field.IsValid() {
			return fmt.Errorf("field '%s' not found in component type %T", dep.fieldName, instance)
		}

		var dependency interface{}

		// Get full interface name using the field's type
		fullTypeName := getFullInterfaceName(field.Type(), dep.typeName)

		// First try to find the dependency in interfaces if it has a qualifier
		if dep.qualifier != "" {
			if impls := c.interfaces[fullTypeName]; impls != nil {
				dependency = impls[dep.qualifier]
			}
		} else {
			// Try to find in interfaces without qualifier
			if impls := c.interfaces[fullTypeName]; impls != nil {
				if len(impls) == 1 {
					// If only one implementation, use it
					for _, v := range impls {
						dependency = v
						break
					}
				} else {
					// If multiple implementations and no qualified specified:
					availableQualifiers := make([]string, 0, len(impls))
					for q := range impls {
						availableQualifiers = append(availableQualifiers, q)
					}
					return fmt.Errorf("multiple implementations found for interface '%s' but no qualifier specified. Available qualifiers: %v", fullTypeName, availableQualifiers)
				}
			}
		}

		// Fallback to direct component if not found in interfaces
		if dependency == nil {
			// Get the full package path for the dependency type
			depType := field.Type()
			if depType.Kind() == reflect.Ptr {
				depType = depType.Elem()
			}
			// Create full dependency name using package path
			depName := getFullComponentName(depType, dep.typeName)
			dependency = c.components[depName]
		}

		if dependency == nil {
			return fmt.Errorf("dependency not found: type='%s' qualifier='%s' - ensure dependency is properly registered as a component", fullTypeName, dep.qualifier)
		}

		// Adjust the dependency type to match the field type
		dependencyValue := reflect.ValueOf(dependency)
		if field.Type() != dependencyValue.Type() && field.Type() == dependencyValue.Elem().Type() {
			dependencyValue = dependencyValue.Elem()
		}

		field.Set(dependencyValue)
	}
	return nil
}

// Get retrieves a component by name from the container.
// The name is converted to lowercase before lookup.
func (c *Container) Get(name string) interface{} {
	// Try to find component by name directly first
	if component := c.components[strings.ToLower(name)]; component != nil {
		return component
	}

	// If not found, try with package path prefix if available
	for fullName, component := range c.components {
		if strings.HasSuffix(fullName, "."+strings.ToLower(name)) {
			return component
		}
	}
	return nil
}

// GetQualified retrieves a specific implementation of an interface using a qualifier.
// Returns nil if the interface or qualified implementation is not found.
func (c *Container) GetQualified(interfaceName, qualifier string) interface{} {
	impls := c.interfaces[interfaceName]
	if impls == nil {
		return nil
	}
	return impls[qualifier]
}
