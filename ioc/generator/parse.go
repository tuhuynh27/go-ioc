package generator

import (
	"bufio"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Component struct {
	Name         string
	Type         string
	Package      string
	Qualifier    string
	Implements   []string
	Dependencies []Dependency
}

type Dependency struct {
	FieldName string
	Type      string
	Qualifier string
}

func getModulePath(rootDir string) (string, error) {
	// Find go.mod file by walking up directories
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

func ParseComponents(rootDir string) ([]Component, error) {
	var components []Component
	fset := token.NewFileSet()

	// Get module path from go.mod
	modulePath, err := getModulePath(rootDir)
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() || strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		log.Printf("Scanning directory: %s", path)

		pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
		if err != nil {
			log.Printf("Error parsing directory %s: %v", path, err)
			return nil
		}

		log.Printf("Found %d packages in %s", len(pkgs), path)

		for pkgName, pkg := range pkgs {
			log.Printf("Processing package: %s", pkgName)

			for fileName, file := range pkg.Files {
				log.Printf("Processing file: %s", fileName)

				ast.Inspect(file, func(n ast.Node) bool {
					typeSpec, ok := n.(*ast.TypeSpec)
					if !ok {
						return true
					}

					structType, ok := typeSpec.Type.(*ast.StructType)
					if !ok {
						return true
					}

					log.Printf("Found struct: %s", typeSpec.Name.Name)

					// Get relative package path from root directory
					relPath, err := filepath.Rel(rootDir, path)
					if err != nil {
						log.Printf("Error getting relative path: %v", err)
						return true
					}

					// Construct full package path using module path
					fullPkgPath := filepath.Join(modulePath, relPath)
					fullPkgPath = strings.ReplaceAll(fullPkgPath, string(filepath.Separator), "/")

					comp := Component{
						Name:    typeSpec.Name.Name,
						Type:    typeSpec.Name.Name,
						Package: fullPkgPath,
					}

					// Check if this struct has a Component field
					hasComponent := false
					for _, field := range structType.Fields.List {
						if ident, ok := field.Type.(*ast.Ident); ok {
							if ident.Name == "Component" {
								hasComponent = true
								log.Printf("Found Component field in %s", comp.Name)
								if field.Tag != nil {
									tag := parseStructTag(field.Tag.Value)
									if name, ok := tag["name"]; ok {
										comp.Name = name
									}
								}
							} else if ident.Name == "Qualifier" {
								if field.Tag != nil {
									tag := parseStructTag(field.Tag.Value)
									if value, ok := tag["value"]; ok {
										comp.Qualifier = value
										log.Printf("Found Qualifier: %s", value)
									}
								}
							}
						} else if sel, ok := field.Type.(*ast.SelectorExpr); ok {
							// Handle Component field from imported package
							if sel.Sel.Name == "Component" {
								hasComponent = true
								log.Printf("Found Component field in %s", comp.Name)
								if field.Tag != nil {
									tag := parseStructTag(field.Tag.Value)
									if name, ok := tag["name"]; ok {
										comp.Name = name
									}
								}
							} else if sel.Sel.Name == "Qualifier" {
								if field.Tag != nil {
									tag := parseStructTag(field.Tag.Value)
									if value, ok := tag["value"]; ok {
										comp.Qualifier = value
										log.Printf("Found Qualifier: %s", value)
									}
								}
							}
						}

						// Check for implements tag
						if field.Tag != nil {
							tag := parseStructTag(field.Tag.Value)
							if impl, ok := tag["implements"]; ok {
								if !strings.Contains(impl, ".") {
									impl = filepath.Join(fullPkgPath, impl)
									impl = strings.ReplaceAll(impl, string(filepath.Separator), "/")
								}
								comp.Implements = append(comp.Implements, impl)
								log.Printf("Found implements: %s", impl)
							}
						}

						// Check for autowired fields
						if field.Tag != nil {
							tag := parseStructTag(field.Tag.Value)
							if _, ok := tag["autowired"]; ok {
								dep := Dependency{
									FieldName: field.Names[0].Name,
									Qualifier: tag["qualifier"],
								}

								// Get the type name and package
								switch t := field.Type.(type) {
								case *ast.Ident:
									dep.Type = t.Name
								case *ast.SelectorExpr:
									if x, ok := t.X.(*ast.Ident); ok {
										dep.Type = x.Name + "." + t.Sel.Name
									}
								}

								log.Printf("Found dependency: %+v", dep)
								comp.Dependencies = append(comp.Dependencies, dep)
							}
						}
					}

					if hasComponent {
						log.Printf("Adding component: %+v", comp)
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

	log.Printf("Total components found: %d", len(components))
	for _, comp := range components {
		log.Printf("Component: %s (package: %s)", comp.Name, comp.Package)
		log.Printf("- Qualifier: %s", comp.Qualifier)
		log.Printf("- Implements: %v", comp.Implements)
		for _, dep := range comp.Dependencies {
			log.Printf("- Dependency: %s (%s) qualifier: %s", dep.FieldName, dep.Type, dep.Qualifier)
		}
	}

	return components, nil
}

func parseStructTag(tag string) map[string]string {
	tag = strings.Trim(tag, "`")
	tags := make(map[string]string)

	for _, section := range strings.Split(tag, " ") {
		if section == "" {
			continue
		}

		parts := strings.Split(section, ":")
		if len(parts) != 2 {
			continue
		}

		key := strings.Trim(parts[0], "\"")
		value := strings.Trim(parts[1], "\"")
		tags[key] = value
	}

	return tags
}
