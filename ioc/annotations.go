package ioc

// Component is a marker struct used to identify injectable components
type Component struct{}

// Qualifier provides a way to disambiguate between multiple implementations
type Qualifier struct {
	Value string
}
