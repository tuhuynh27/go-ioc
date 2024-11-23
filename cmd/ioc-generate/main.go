package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/tuhuynh27/go-ioc/ioc/generator"
)

func main() {
	var (
		dir     = flag.String("dir", ".", "Directory to scan for components")
		output  = flag.String("output", "wire/wire_gen.go", "Output file for generated code")
		verbose = flag.Bool("verbose", false, "Enable verbose logging")
	)
	flag.Parse()

	// Convert to absolute path
	absDir, err := filepath.Abs(*dir)
	if err != nil {
		log.Fatalf("Error getting absolute path: %v", err)
	}

	log.Printf("Scanning directory: %s", absDir)

	// Parse components
	components, err := generator.ParseComponents(absDir)

	if err != nil {
		log.Fatalf("Error parsing components: %v", err)
	}

	if *verbose {
		log.Printf("Found %d components", len(components))
		for _, comp := range components {
			log.Printf("Component: %s (package: %s)", comp.Name, comp.Package)
			log.Printf("- Qualifier: %s", comp.Qualifier)
			log.Printf("- Implements: %v", comp.Implements)
			log.Printf("- Dependencies: %v", comp.Dependencies)
		}
	}

	// Generate code
	gen := generator.NewGenerator(components)
	if err := gen.Generate(absDir); err != nil {
		log.Fatalf("Error generating code: %v", err)
	}

	log.Printf("Successfully generated wire file: %s/%s", absDir, *output)
}
