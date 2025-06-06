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
		Run: func(cmd *cobra.Command, args []string) {
			if help {
				cmd.Help()
				return
			}

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
					log.Printf("- Source: %s:%d", comp.SourceFile, comp.LineNumber)
					log.Printf("- Qualifier: %s", comp.Qualifier)
					log.Printf("- Implements: %v", comp.Implements)
					log.Printf("- Dependencies: %v", comp.Dependencies)
				}
			}

			// Create generator
			gen := wire.NewGenerator(components)

			// Handle special modes
			if showGraph {
				gen.PrintDependencyGraph()
				return
			}

			if dryRun {
				if err := gen.ValidateOnly(); err != nil {
					log.Fatalf("Validation failed: %v", err)
				}
				return
			}

			if listComponents {
				gen.ListComponents()
				return
			}

			if analyzeComponents {
				analyzer := wire.NewAnalyzer(components)
				analyzer.PrintAnalysisReport()
				return
			}

			// Generate code
			if err := gen.Generate(absDir); err != nil {
				log.Fatalf("Error generating code: %v", err)
			}

			log.Printf("Successfully generated wire file: %s/%s", absDir, output)
		},
	}

	dir, output                           string
	verbose, help                         bool
	showGraph, dryRun                     bool
	listComponents, analyzeComponents     bool
)

func main() {
	rootCmd.PersistentFlags().StringVarP(&dir, "dir", "d", ".", "Directory to scan for components")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "wire/wire_gen.go", "Output file for generated code")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().BoolVarP(&help, "help", "h", false, "Show help message")
	rootCmd.PersistentFlags().BoolVar(&showGraph, "graph", false, "Show dependency graph visualization")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Validate components without generating files")
	rootCmd.PersistentFlags().BoolVar(&listComponents, "list", false, "List all discovered components")
	rootCmd.PersistentFlags().BoolVar(&analyzeComponents, "analyze", false, "Perform comprehensive component analysis")

	printBanner()
	rootCmd.Execute()
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
