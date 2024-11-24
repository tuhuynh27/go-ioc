package wire

import (
	"bufio"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Component represents a parsed IoC component with its metadata
type Component struct {
	Name          string       // Name of the component (from name tag or type name)
	Type          string       // The Go type name
	Package       string       // Full package path
	Qualifier     string       // Qualifier value for disambiguation
	Implements    []string     // Interfaces implemented by this component
	Dependencies  []Dependency // List of autowired dependencies
	PostConstruct bool         // Whether component has PostConstruct method
	PreDestroy    bool         // Whether component has PreDestroy method
}

// Dependency represents an autowired dependency field in a component
type Dependency struct {
	FieldName string // Name of the struct field
	Type      string // Type of the dependency
	Qualifier string // Qualifier for selecting specific implementation
}

// getModulePath finds and returns the Go module path from go.mod file
// by walking up directories from rootDir until go.mod is found
func getModulePath(rootDir string) (string, error) {
	dir := rootDir
	for {
		modPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(modPath); err == nil {
			// Found go.mod, read module path
			file, err := os.Open(modPath)
			if err != nil {
				return "", err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "module ") {
					return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
				}
			}
			return "", scanner.Err()
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

// ParseComponents scans the given root directory for Go files containing IoC components
// and returns a slice of parsed Component structs
func ParseComponents(rootDir string) ([]Component, error) {
	startTime := time.Now()
	var components []Component
	fset := token.NewFileSet() // Used for parsing Go source files

	// First get the module path which is needed for constructing full package paths
	modulePath, err := getModulePath(rootDir)
	if err != nil {
		return nil, err
	}

	// Walk through all directories under rootDir
	err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-directories and hidden directories (starting with .)
		if !info.IsDir() {
			return nil
		}
		if strings.HasPrefix(filepath.Base(path), ".") {
			return filepath.SkipDir
		}

		// Parse all Go files in the current directory
		pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
		if err != nil {
			log.Printf("Error parsing directory %s: %v", path, err)
			return nil
		}

		// Process each package
		for _, pkg := range pkgs {
			// Process each file in the package
			for _, file := range pkg.Files {
				// Inspect the AST of each file
				ast.Inspect(file, func(n ast.Node) bool {
					// Look for type declarations
					typeSpec, ok := n.(*ast.TypeSpec)
					if !ok {
						return true
					}

					// Check if it's a struct type
					structType, ok := typeSpec.Type.(*ast.StructType)
					if !ok {
						return true
					}

					// Get relative package path and construct full package path
					relPath, err := filepath.Rel(rootDir, path)
					if err != nil {
						log.Printf("Error getting relative path: %v", err)
						return true
					}
					fullPkgPath := filepath.Join(modulePath, relPath)
					fullPkgPath = strings.ReplaceAll(fullPkgPath, string(filepath.Separator), "/")

					// Initialize component with basic info
					comp := Component{
						Name:    typeSpec.Name.Name,
						Type:    typeSpec.Name.Name,
						Package: fullPkgPath,
					}

					// Analyze struct fields for component markers and metadata
					hasComponent := false
					for _, field := range structType.Fields.List {
						// Check if field is the IoC Component marker
						if len(field.Names) > 0 && field.Names[0].Name == "Component" {
							if _, ok := field.Type.(*ast.StructType); ok {
								hasComponent = true
								// Check for name override in tag
								if field.Tag != nil {
									tag := parseStructTag(field.Tag.Value)
									if name, ok := tag["name"]; ok {
										comp.Name = name
									}
								}
							}
						}

						// Process struct tags if present
						if field.Tag != nil {
							tag := parseStructTag(field.Tag.Value)

							// Check if this field has a "value" tag and is a Qualifier field
							value, hasValue := tag["value"]
							isQualifierField := len(field.Names) > 0 && field.Names[0].Name == "Qualifier"
							isEmptyStruct := false
							if _, ok := field.Type.(*ast.StructType); ok {
								isEmptyStruct = true
							}

							if hasValue && isQualifierField && isEmptyStruct {
								comp.Qualifier = value
							}

							// Check for implements declarations
							if impl, ok := tag["implements"]; ok {
								// Convert relative interface paths to absolute
								if !strings.Contains(impl, ".") {
									impl = filepath.Join(fullPkgPath, impl)
									impl = strings.ReplaceAll(impl, string(filepath.Separator), "/")
								}
								comp.Implements = append(comp.Implements, impl)
							}

							// Process autowired dependencies
							if _, ok := tag["autowired"]; ok {
								dep := Dependency{
									FieldName: field.Names[0].Name,
									Qualifier: tag["qualifier"],
								}

								// Extract type information based on AST node type
								var typ string
								switch t := field.Type.(type) {
								case *ast.Ident: // Simple type
									typ = t.Name
								case *ast.StarExpr: // Pointer type
									switch x := t.X.(type) {
									case *ast.Ident:
										typ = x.Name
									case *ast.SelectorExpr:
										if y, ok := x.X.(*ast.Ident); ok {
											typ = y.Name + "." + x.Sel.Name
										}
									}
								case *ast.SelectorExpr: // Qualified type
									if x, ok := t.X.(*ast.Ident); ok {
										typ = x.Name + "." + t.Sel.Name
									}
								}
								dep.Type = typ

								comp.Dependencies = append(comp.Dependencies, dep)
							}
						}
					}

					// Only add if it's a valid component
					if hasComponent {
						// Check for PostConstruct and PreDestroy methods
						for _, decl := range file.Decls {
							if funcDecl, ok := decl.(*ast.FuncDecl); ok {
								// Check if it's a method on our component type
								if funcDecl.Recv != nil && len(funcDecl.Recv.List) == 1 {
									recvType := funcDecl.Recv.List[0].Type
									if starExpr, ok := recvType.(*ast.StarExpr); ok {
										if ident, ok := starExpr.X.(*ast.Ident); ok {
											if ident.Name == comp.Type {
												// Check for PostConstruct method
												if funcDecl.Name.Name == "PostConstruct" {
													// Check if method has no parameters and no return values
													if funcDecl.Type.Params.NumFields() == 0 &&
														(funcDecl.Type.Results == nil || funcDecl.Type.Results.NumFields() == 0) {
														comp.PostConstruct = true
													}
												}
												// Check for PreDestroy method
												if funcDecl.Name.Name == "PreDestroy" {
													// Check if method has no parameters and no return values
													if funcDecl.Type.Params.NumFields() == 0 &&
														(funcDecl.Type.Results == nil || funcDecl.Type.Results.NumFields() == 0) {
														comp.PreDestroy = true
													}
												}
											}
										}
									}
								}
							}
						}
						components = append(components, comp)
					}

					return true
				})
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	log.Printf("Found %d components (scan completed in %v)", len(components), time.Since(startTime))

	return components, nil
}

// parseStructTag parses a Go struct tag string into a map of key-value pairs
func parseStructTag(tag string) map[string]string {
	tag = strings.Trim(tag, "`")
	tags := make(map[string]string)

	// Split tag into space-separated sections
	for _, section := range strings.Split(tag, " ") {
		if section == "" {
			continue
		}

		// Split each section into key:value
		parts := strings.Split(section, ":")
		if len(parts) != 2 {
			continue
		}

		// Clean up quotes and store in map
		key := strings.Trim(parts[0], "\"")
		value := strings.Trim(parts[1], "\"")
		tags[key] = value
	}

	return tags
}
