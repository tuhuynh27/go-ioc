// Package ioc provides a dependency injection container implementation for Go applications.
package ioc

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

// RegisterInterface registers a component as an implementation of an interface
func (c *Container) RegisterWithInterface(interfaceName string, qualifier string, component interface{}) {
	if c.interfaces[interfaceName] == nil {
		c.interfaces[interfaceName] = make(map[string]interface{})
	}
	c.interfaces[interfaceName][qualifier] = component
}

// Get retrieves a component by name from the container
func (c *Container) Get(name string) interface{} {
	return c.components[name]
}

// GetQualified retrieves a specific implementation of an interface using a qualifier
func (c *Container) GetQualified(interfaceName, qualifier string) interface{} {
	if impls := c.interfaces[interfaceName]; impls != nil {
		return impls[qualifier]
	}
	return nil
}
