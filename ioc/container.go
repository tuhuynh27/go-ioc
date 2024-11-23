// Package ioc provides a dependency injection container implementation for Go applications.
package ioc

import (
	"fmt"
	"strings"
)

// Container manages component registration and retrieval
type Container struct {
	components map[string]interface{}            // Stores instantiated components by name
	interfaces map[string]map[string]interface{} // Maps interface names to qualified implementations
}

// NewContainer creates and initializes a new dependency injection container
func NewContainer() *Container {
	return &Container{
		components: make(map[string]interface{}),
		interfaces: make(map[string]map[string]interface{}),
	}
}

// Register adds a component to the container
func (c *Container) Register(name string, component interface{}) {
	c.components[name] = component
}

// RegisterWithInterface registers a component as an implementation of an interface
func (c *Container) RegisterWithInterface(interfaceName string, qualifier string, component interface{}) {
	if c.interfaces[interfaceName] == nil {
		c.interfaces[interfaceName] = make(map[string]interface{})
	}
	c.interfaces[interfaceName][qualifier] = component
}

// Get retrieves a component by name from the container
// If multiple components match the short name, returns an error
func (c *Container) Get(name string) (interface{}, error) {
	// If the name contains full path, use direct lookup
	if component, exists := c.components[name]; exists {
		return component, nil
	}

	// Search for components matching the short name
	var matches []string
	var matchedComponent interface{}

	for fullPath, component := range c.components {
		parts := strings.Split(fullPath, ".")
		shortName := parts[len(parts)-1]
		if shortName == name {
			matches = append(matches, fullPath)
			matchedComponent = component
		}
	}

	switch len(matches) {
	case 0:
		return nil, fmt.Errorf("no component found with name: %s", name)
	case 1:
		return matchedComponent, nil
	default:
		return nil, fmt.Errorf("multiple components found with name %s: %v", name, matches)
	}
}

// GetQualified retrieves a specific implementation of an interface using a qualifier
// If multiple interfaces match the short name, returns an error
func (c *Container) GetQualified(interfaceName, qualifier string) (interface{}, error) {
	// If the interface name contains full path, use direct lookup
	if impls, exists := c.interfaces[interfaceName]; exists {
		if component, exists := impls[qualifier]; exists {
			return component, nil
		}
		return nil, fmt.Errorf("no component found for interface %s with qualifier %s", interfaceName, qualifier)
	}

	// Search for interfaces matching the short name
	var matches []string
	var matchedComponent interface{}

	for fullPath, impls := range c.interfaces {
		parts := strings.Split(fullPath, ".")
		shortName := parts[len(parts)-1]
		if shortName == interfaceName {
			matches = append(matches, fullPath)
			if component, exists := impls[qualifier]; exists {
				matchedComponent = component
			}
		}
	}

	switch len(matches) {
	case 0:
		return nil, fmt.Errorf("no interface found with name: %s", interfaceName)
	case 1:
		if matchedComponent == nil {
			return nil, fmt.Errorf("no component found for interface %s with qualifier %s", matches[0], qualifier)
		}
		return matchedComponent, nil
	default:
		return nil, fmt.Errorf("multiple interfaces found with name %s: %v", interfaceName, matches)
	}
}
