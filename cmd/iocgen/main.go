package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tuhuynh27/go-ioc/internal/wire"
)

var (
	rootCmd = &cobra.Command{
		Use:   "ioc-generate",
		Short: "Go IoC - Dependency Injection Code Generator",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	dir, output   string
	verbose, help bool
)

func init() {
	rootCmd.InitDefaultHelpCmd()
}

func main() {
	printBanner()

	rootCmd.PersistentFlags().StringVar(&dir, "dir", ".", "Directory to scan for components")
	rootCmd.PersistentFlags().StringVar(&output, "output", "wire/wire_gen.go", "Output file for generated code")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	rootCmd.PersistentFlags().BoolVar(&help, "help", false, "Show help message")

	rootCmd.Execute()

	// Convert to absolute path
	absDir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalf("Error getting absolute path: %v", err)
	}

	log.Printf("Scanning directory: %s", absDir)

	// Parse components
	components, err := wire.ParseComponents(absDir)

	if err != nil {
		log.Fatalf("Error parsing components: %v", err)
	}

	if verbose {
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

	log.Printf("Successfully generated wire file: %s/%s", absDir, output)
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
