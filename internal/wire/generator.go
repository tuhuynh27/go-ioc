package wire

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"
)

// Generator handles the code generation process for dependency injection
type Generator struct {
	components []Component     // List of all components to be wired
	visited    map[string]bool // Tracks visited components during dependency resolution
	cyclicMap  map[string]bool // Tracks cyclic dependencies
}

// templateData holds the data needed for code template generation
type templateData struct {
	Imports    []string        // List of packages to import
	Components []componentInit // List of component initializations
}

// componentInit represents a single component's initialization data
type componentInit struct {
	VarName       string         // Variable name for the component instance
	Type          string         // Component type name
	Package       string         // Package path where component is defined
	Dependencies  []componentDep // List of dependencies for this component
	Interfaces    []interfaceReg // List of interfaces this component implements
	PostConstruct bool           // Whether component has PostConstruct method
	PreDestroy    bool           // Whether component has PreDestroy method
}

// componentDep represents a single dependency of a component
type componentDep struct {
	FieldName string // Name of the field in the struct
	VarName   string // Variable name of the dependency instance
}

// interfaceReg represents an interface implementation registration
type interfaceReg struct {
	Interface string // Name of the interface being implemented
	Qualifier string // Qualifier value for disambiguation
	VarName   string // Variable name of the implementing component
}

// NewGenerator creates a new Generator instance with the provided components
func NewGenerator(components []Component) *Generator {
	return &Generator{
		components: components,
		visited:    make(map[string]bool),
		cyclicMap:  make(map[string]bool),
	}
}

// Generate performs the code generation process for dependency injection
func (g *Generator) Generate(baseDir string) error {
	startTime := time.Now()

	// Validate that we have components to process
	if len(g.components) == 0 {
		return fmt.Errorf("no components found")
	}

	// Ensure the wire directory exists for output
	wireDir := filepath.Join(baseDir, "wire")
	if err := os.MkdirAll(wireDir, 0755); err != nil {
		return fmt.Errorf("failed to create wire directory: %w", err)
	}

	// Sort components based on their dependencies
	orderedComponents := g.topologicalSort()
	// Generate initialization code for each component
	inits := g.generateComponentInits(orderedComponents)

	// Define template helper functions
	funcMap := template.FuncMap{
		"base": func(pkg string) string {
			parts := strings.Split(pkg, "/")
			return parts[len(parts)-1]
		},
		"iterate": func(count int) []int {
			var result []int
			for i := count - 1; i >= 0; i-- {
				result = append(result, i)
			}
			return result
		},
	}

	// Create and parse the code generation template
	tmpl := template.New("wire").Funcs(funcMap)
	tmpl, err := tmpl.Parse(`// File: wire_gen.go
// Code generated by Go IoC. DO NOT EDIT.
//go:generate go run github.com/tuhuynh27/go-ioc/cmd/iocgen -dir=../
package wire

import ({{range .Imports}}
    "{{.}}"{{end}}
)

type Container struct {
    {{- range $comp := .Components}}
    {{$comp.Type}} *{{$comp.Package | base}}.{{$comp.Type}}
    {{- end}}
}

func Initialize() (*Container, func()) {
    container := &Container{}{{range $comp := .Components}}
    container.{{$comp.Type}} = &{{$comp.Package | base}}.{{$comp.Type}}{{if $comp.Dependencies}}{
        {{- range $dep := $comp.Dependencies}}
        {{$dep.FieldName}}: container.{{$dep.VarName}},{{- end}}
    }{{else}}{}{{end}}{{if $comp.PostConstruct}}
    container.{{$comp.Type}}.PostConstruct(){{- end}}{{end}}

    cleanup := func() {
        {{- range $i := len .Components | iterate}}
        {{- with index $.Components $i}}
        {{- if .PreDestroy}}
        container.{{.Type}}.PreDestroy()
        {{- end}}
        {{- end}}
        {{- end}}
    }

    return container, cleanup
}
`)
	if err != nil {
		return fmt.Errorf("template parsing failed: %w", err)
	}

	// Collect and deduplicate required imports
	imports := make(map[string]bool)
	for _, comp := range g.components {
		// Only include imports with valid package paths
		if strings.Contains(comp.Package, "/") && comp.Package != "github.com/tuhuynh27/go-ioc/ioc" {
			imports[comp.Package] = true
		}
		// Add imports for implemented interfaces
		for _, iface := range comp.Implements {
			if idx := strings.LastIndex(iface, "."); idx != -1 {
				pkgPath := iface[:idx]
				if strings.Contains(pkgPath, "/") && pkgPath != "github.com/tuhuynh27/go-ioc/ioc" {
					imports[pkgPath] = true
				}
			}
		}
	}

	// Convert imports map to sorted slice for consistent output
	var importSlice []string
	for imp := range imports {
		importSlice = append(importSlice, imp)
	}
	sort.Strings(importSlice)

	// Prepare data for template execution
	data := templateData{
		Imports:    importSlice,
		Components: inits,
	}

	// Generate the code using the template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	// Write the generated code to file
	outputPath := filepath.Join(wireDir, "wire_gen.go")
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write generated code: %w", err)
	}

	log.Printf("Generated wire_gen.go in %s (completed in %v)", wireDir, time.Since(startTime))

	return nil
}

// findPackageForType finds the full package path for a given package name
func (g *Generator) findPackageForType(pkgName string) string {
	for _, comp := range g.components {
		if strings.HasSuffix(comp.Package, "/"+pkgName) {
			return comp.Package
		}
	}
	return ""
}

// generateComponentInits creates initialization data for all components
func (g *Generator) generateComponentInits(components []Component) []componentInit {
	var inits []componentInit
	varNames := make(map[string]string)
	nameCount := make(map[string]int)

	// First pass: generate unique variable names for all components
	for _, comp := range components {
		baseName := comp.Type

		count := nameCount[baseName]
		nameCount[baseName]++

		varName := baseName
		if count > 0 {
			varName = fmt.Sprintf("%s%d", baseName, count+1)
		}

		varNames[comp.Package+"."+comp.Type] = varName
	}

	// Second pass: create component initializations with dependencies and interface registrations
	for _, comp := range components {
		init := componentInit{
			VarName:       varNames[comp.Package+"."+comp.Type],
			Type:          comp.Type,
			Package:       comp.Package,
			Dependencies:  []componentDep{},
			Interfaces:    []interfaceReg{},
			PostConstruct: comp.PostConstruct,
			PreDestroy:    comp.PreDestroy,
		}

		// Process each dependency for the component
		for _, dep := range comp.Dependencies {
			depVarName := ""
			depType := dep.Type

			if strings.Contains(depType, ".") {
				// Handle interface dependencies
				parts := strings.Split(depType, ".")
				pkgName := parts[0]
				interfaceName := parts[1]

				// Find matching component that implements the interface
				for _, c := range components {
					if strings.HasSuffix(c.Package, pkgName) {
						// Check interface implementations
						for _, impl := range c.Implements {
							if strings.HasSuffix(impl, interfaceName) && c.Qualifier == dep.Qualifier {
								depVarName = varNames[c.Package+"."+c.Type]
								break
							}
						}
						// Check direct type match
						if c.Type == interfaceName && c.Qualifier == dep.Qualifier {
							depVarName = varNames[c.Package+"."+c.Type]
							break
						}
					}
				}
			} else {
				// Handle direct type dependencies
				for _, c := range components {
					if c.Type == depType {
						depVarName = varNames[c.Package+"."+c.Type]
						break
					}
				}
			}

			// Add dependency if found
			if depVarName != "" {
				init.Dependencies = append(init.Dependencies, componentDep{
					FieldName: dep.FieldName,
					VarName:   depVarName,
				})
			} else {
				panic(fmt.Sprintf("Error: Could not resolve dependency '%s' with qualifier '%s' for component '%s.%s' in field '%s'.\n"+
					"Please check:\n"+
					"1. The dependency type exists and is marked with Component{}\n"+
					"2. For interface dependencies, ensure an implementation struct tag is declared\n"+
					"3. If using qualifiers, verify the qualifier values match\n"+
					"4. The package containing the dependency is included in component scanning\n"+
					"If you want more debugging information, re run the command with -verbose flag",
					dep.Type, dep.Qualifier, init.Package, init.Type, dep.FieldName))
			}
		}

		// Register interfaces implemented by this component
		for _, iface := range comp.Implements {
			init.Interfaces = append(init.Interfaces, interfaceReg{
				Interface: iface,
				Qualifier: comp.Qualifier,
				VarName:   init.VarName,
			})
		}

		inits = append(inits, init)
	}

	return inits
}

// contains checks if a slice contains a string (with suffix matching)
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if strings.HasSuffix(s, str) {
			return true
		}
	}
	return false
}

// topologicalSort sorts components based on their dependencies
func (g *Generator) topologicalSort() []Component {
	var ordered []Component

	// First process components with no dependencies
	for _, comp := range g.components {
		if len(comp.Dependencies) == 0 && !g.visited[comp.Package+"."+comp.Type] {
			g.dfs(comp, &ordered)
		}
	}

	// Then process remaining components
	for _, comp := range g.components {
		if !g.visited[comp.Package+"."+comp.Type] {
			g.dfs(comp, &ordered)
		}
	}

	return ordered
}

// dfs performs a depth-first search to order components by dependencies
func (g *Generator) dfs(comp Component, ordered *[]Component) {
	componentKey := comp.Package + "." + comp.Type

	if g.cyclicMap[componentKey] {
		// Build the dependency cycle path
		var cyclePath []string
		for k := range g.cyclicMap {
			cyclePath = append(cyclePath, k)
		}
		cyclePath = append(cyclePath, componentKey)

		panic(fmt.Sprintf("Cyclic dependency detected!\n"+
			"Component %s.%s depends on a component that eventually depends back on it.\n"+
			"Dependency path: %s\n"+
			"Please check these components and their dependencies to resolve the cycle.",
			comp.Package, comp.Type,
			strings.Join(cyclePath, " -> ")))
	}

	g.cyclicMap[componentKey] = true
	defer delete(g.cyclicMap, componentKey)

	if g.visited[componentKey] {
		return
	}

	g.visited[componentKey] = true

	// Process dependencies before the component itself
	for _, dep := range comp.Dependencies {
		if strings.Contains(dep.Type, ".") {
			parts := strings.Split(dep.Type, ".")
			pkgName := parts[0]
			typeName := parts[1]

			// Find and visit each dependency
			for _, other := range g.components {
				if strings.HasSuffix(other.Package, pkgName) {
					// Check both direct type matches and interface implementations
					if other.Type == typeName || contains(other.Implements, typeName) {
						g.dfs(other, ordered)
					}
				}
			}
		} else {
			// Handle direct type dependencies
			for _, other := range g.components {
				if other.Type == dep.Type {
					g.dfs(other, ordered)
				}
			}
		}
	}

	*ordered = append(*ordered, comp)
}
