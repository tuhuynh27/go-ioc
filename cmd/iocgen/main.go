package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tuhuynh27/go-ioc/internal/wire"
)

func main() {
	printBanner()

	// Check for help flag
	if len(os.Args) > 1 && (os.Args[1] == "-help" || os.Args[1] == "help") {
		printHelp()
		return
	}

	var (
		dir     = flag.String("dir", ".", "Directory to scan for components")
		output  = flag.String("output", "wire/wire_gen.go", "Output file for generated code")
		verbose = flag.Bool("verbose", false, "Enable verbose logging")
		help    = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	if *help {
		printHelp()
		return
	}

	// Convert to absolute path
	absDir, err := filepath.Abs(*dir)
	if err != nil {
		log.Fatalf("Error getting absolute path: %v", err)
	}

	log.Printf("Scanning directory: %s", absDir)

	// Parse components
	components, err := wire.ParseComponents(absDir)

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
	gen := wire.NewGenerator(components)
	if err := gen.Generate(absDir); err != nil {
		log.Fatalf("Error generating code: %v", err)
	}

	log.Printf("Successfully generated wire file: %s/%s", absDir, *output)
}

func printHelp() {
	fmt.Println("Go IoC - Dependency Injection Code Generator")
	fmt.Println("\nUsage:")
	fmt.Println("  ioc-generate [flags]")
	fmt.Println("\nFlags:")
	fmt.Println("  -dir string")
	fmt.Println("        Directory to scan for components (default \".\")")
	fmt.Println("  -output string")
	fmt.Println("        Output file for generated code (default \"wire/wire_gen.go\")")
	fmt.Println("  -verbose")
	fmt.Println("        Enable verbose logging")
	fmt.Println("  -help")
	fmt.Println("        Show this help message")
	fmt.Println("\nExample:")
	fmt.Println("  ioc-generate -dir=./src -output=wire/generated.go -verbose")
	os.Exit(0)
}

func printBanner() {
	fmt.Println(`
   ______      _____ ____  ______
  / ____/___  /  _/ / __ \/ ____/
 / / __/ __ \ / // / / / / /     
/ /_/ / /_/ // // / /_/ / /___   
\____/\____/___/_/\____/\____/   
                                 
Inversion of Control for Go
Version 0.0.0`)
}
