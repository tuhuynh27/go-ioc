package wire

import (
	"fmt"
	"sort"
	"strings"
)

// DependencyAnalyzer provides advanced analysis capabilities for IoC components
type DependencyAnalyzer struct {
	components []Component
}

// NewAnalyzer creates a new DependencyAnalyzer instance
func NewAnalyzer(components []Component) *DependencyAnalyzer {
	return &DependencyAnalyzer{components: components}
}

// AnalysisResult contains comprehensive analysis results
type AnalysisResult struct {
	TotalComponents       int
	TotalDependencies     int
	CircularDependencies  []CircularDependency
	UnusedComponents      []Component
	OrphanedComponents    []Component
	InterfaceAnalysis     InterfaceAnalysis
	QualifierConflicts    []QualifierConflict
	DependencyDepth       map[string]int
	ComponentsByPackage   map[string][]Component
}

// CircularDependency represents a detected circular dependency
type CircularDependency struct {
	Path        []string
	Description string
}

// InterfaceAnalysis provides insights about interface usage
type InterfaceAnalysis struct {
	TotalInterfaces         int
	InterfacesWithMultiImpl map[string][]string
	InterfacesWithNoImpl    []string
	ImplementationCount     map[string]int
}

// QualifierConflict represents a qualifier naming conflict
type QualifierConflict struct {
	Interface    string
	Qualifier    string
	Conflicting  []string
	Severity     string
}

// PerformComprehensiveAnalysis runs all analysis functions and returns results
func (a *DependencyAnalyzer) PerformComprehensiveAnalysis() *AnalysisResult {
	result := &AnalysisResult{
		TotalComponents:     len(a.components),
		TotalDependencies:   a.countTotalDependencies(),
		ComponentsByPackage: a.groupComponentsByPackage(),
	}

	result.CircularDependencies = a.FindCircularDependencies()
	result.UnusedComponents = a.FindUnusedComponents()
	result.OrphanedComponents = a.FindOrphanedComponents()
	result.InterfaceAnalysis = a.AnalyzeInterfaces()
	result.QualifierConflicts = a.FindQualifierConflicts()
	result.DependencyDepth = a.CalculateDependencyDepth()

	return result
}

// FindCircularDependencies detects circular dependency chains
func (a *DependencyAnalyzer) FindCircularDependencies() []CircularDependency {
	var cycles []CircularDependency
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	
	for _, comp := range a.components {
		key := comp.Package + "." + comp.Type
		if !visited[key] {
			if path := a.dfsCircular(comp, visited, recStack, []string{}); len(path) > 0 {
				cycles = append(cycles, CircularDependency{
					Path: path,
					Description: fmt.Sprintf("Circular dependency detected: %s", strings.Join(path, " → ")),
				})
			}
		}
	}
	
	return cycles
}

// dfsCircular performs depth-first search to detect circular dependencies
func (a *DependencyAnalyzer) dfsCircular(comp Component, visited, recStack map[string]bool, path []string) []string {
	key := comp.Package + "." + comp.Type
	visited[key] = true
	recStack[key] = true
	path = append(path, key)

	// Check dependencies
	for _, dep := range comp.Dependencies {
		depComponents := a.findDependencyComponents(dep)
		for _, depComp := range depComponents {
			depKey := depComp.Package + "." + depComp.Type
			
			if !visited[depKey] {
				if cyclePath := a.dfsCircular(depComp, visited, recStack, path); len(cyclePath) > 0 {
					return cyclePath
				}
			} else if recStack[depKey] {
				// Found cycle - return the cycle path
				cycleStart := -1
				for i, p := range path {
					if p == depKey {
						cycleStart = i
						break
					}
				}
				if cycleStart != -1 {
					cycle := append(path[cycleStart:], depKey)
					return cycle
				}
			}
		}
	}

	recStack[key] = false
	return nil
}

// FindUnusedComponents identifies components that are not dependencies of any other component
func (a *DependencyAnalyzer) FindUnusedComponents() []Component {
	used := make(map[string]bool)
	
	// Mark components that are dependencies of others
	for _, comp := range a.components {
		for _, dep := range comp.Dependencies {
			// Find all components that satisfy this dependency
			satisfyingComponents := a.findDependencyComponents(dep)
			for _, satisfying := range satisfyingComponents {
				compKey := satisfying.Package + "." + satisfying.Type
				used[compKey] = true
			}
		}
	}
	
	var unused []Component
	for _, comp := range a.components {
		compKey := comp.Package + "." + comp.Type
		if !used[compKey] {
			unused = append(unused, comp)
		}
	}
	
	return unused
}

// FindOrphanedComponents identifies components with dependencies that cannot be satisfied
func (a *DependencyAnalyzer) FindOrphanedComponents() []Component {
	var orphaned []Component
	
	for _, comp := range a.components {
		hasUnsatisfiedDep := false
		for _, dep := range comp.Dependencies {
			if len(a.findDependencyComponents(dep)) == 0 {
				hasUnsatisfiedDep = true
				break
			}
		}
		if hasUnsatisfiedDep {
			orphaned = append(orphaned, comp)
		}
	}
	
	return orphaned
}

// AnalyzeInterfaces provides comprehensive interface usage analysis
func (a *DependencyAnalyzer) AnalyzeInterfaces() InterfaceAnalysis {
	analysis := InterfaceAnalysis{
		InterfacesWithMultiImpl: make(map[string][]string),
		ImplementationCount:     make(map[string]int),
	}
	
	interfaceImpls := make(map[string][]string)
	usedInterfaces := make(map[string]bool)
	
	// Collect all interface implementations
	for _, comp := range a.components {
		for _, iface := range comp.Implements {
			interfaceImpls[iface] = append(interfaceImpls[iface], comp.Package+"."+comp.Type)
			analysis.ImplementationCount[iface]++
		}
	}
	
	// Mark interfaces that are actually used as dependencies
	for _, comp := range a.components {
		for _, dep := range comp.Dependencies {
			if strings.Contains(dep.Type, ".") {
				usedInterfaces[dep.Type] = true
			}
		}
	}
	
	analysis.TotalInterfaces = len(interfaceImpls)
	
	// Find interfaces with multiple implementations
	for iface, impls := range interfaceImpls {
		if len(impls) > 1 {
			analysis.InterfacesWithMultiImpl[iface] = impls
		}
	}
	
	// Find interfaces with no implementations (used but not implemented)
	for iface := range usedInterfaces {
		if _, exists := interfaceImpls[iface]; !exists {
			analysis.InterfacesWithNoImpl = append(analysis.InterfacesWithNoImpl, iface)
		}
	}
	
	return analysis
}

// FindQualifierConflicts identifies naming conflicts in qualifiers
func (a *DependencyAnalyzer) FindQualifierConflicts() []QualifierConflict {
	var conflicts []QualifierConflict
	qualifierMap := make(map[string]map[string][]string) // interface -> qualifier -> components
	
	// Build qualifier mapping
	for _, comp := range a.components {
		for _, iface := range comp.Implements {
			if qualifierMap[iface] == nil {
				qualifierMap[iface] = make(map[string][]string)
			}
			qualifierMap[iface][comp.Qualifier] = append(
				qualifierMap[iface][comp.Qualifier], 
				comp.Package+"."+comp.Type,
			)
		}
	}
	
	// Find conflicts
	for iface, qualifiers := range qualifierMap {
		for qualifier, components := range qualifiers {
			if len(components) > 1 {
				severity := "ERROR" // Multiple components with same qualifier
				if qualifier == "" {
					severity = "WARNING" // Empty qualifier with multiple implementations
				}
				
				conflicts = append(conflicts, QualifierConflict{
					Interface:   iface,
					Qualifier:   qualifier,
					Conflicting: components,
					Severity:    severity,
				})
			}
		}
	}
	
	return conflicts
}

// CalculateDependencyDepth calculates the maximum dependency depth for each component
func (a *DependencyAnalyzer) CalculateDependencyDepth() map[string]int {
	depths := make(map[string]int)
	calculating := make(map[string]bool)
	
	var calculateDepth func(comp Component) int
	calculateDepth = func(comp Component) int {
		key := comp.Package + "." + comp.Type
		
		if depth, exists := depths[key]; exists {
			return depth
		}
		
		if calculating[key] {
			return -1 // Circular dependency
		}
		
		calculating[key] = true
		defer delete(calculating, key)
		
		maxDepth := 0
		for _, dep := range comp.Dependencies {
			depComponents := a.findDependencyComponents(dep)
			for _, depComp := range depComponents {
				depDepth := calculateDepth(depComp)
				if depDepth == -1 {
					depths[key] = -1
					return -1
				}
				if depDepth >= maxDepth {
					maxDepth = depDepth + 1
				}
			}
		}
		
		depths[key] = maxDepth
		return maxDepth
	}
	
	for _, comp := range a.components {
		calculateDepth(comp)
	}
	
	return depths
}

// groupComponentsByPackage groups components by their package
func (a *DependencyAnalyzer) groupComponentsByPackage() map[string][]Component {
	groups := make(map[string][]Component)
	
	for _, comp := range a.components {
		groups[comp.Package] = append(groups[comp.Package], comp)
	}
	
	return groups
}

// countTotalDependencies counts all dependencies across components
func (a *DependencyAnalyzer) countTotalDependencies() int {
	total := 0
	for _, comp := range a.components {
		total += len(comp.Dependencies)
	}
	return total
}

// findDependencyComponents finds all components that satisfy a given dependency
func (a *DependencyAnalyzer) findDependencyComponents(dep Dependency) []Component {
	var matches []Component
	
	for _, comp := range a.components {
		// Direct type match - check both full package.Type and just Type
		compKey := comp.Package + "." + comp.Type
		if (comp.Type == dep.Type || compKey == dep.Type) && comp.Qualifier == dep.Qualifier {
			matches = append(matches, comp)
			continue
		}
		
		// Interface implementation match
		if strings.Contains(dep.Type, ".") {
			depParts := strings.Split(dep.Type, ".")
			depInterface := depParts[len(depParts)-1] // Get the interface name
			
			for _, iface := range comp.Implements {
				ifaceParts := strings.Split(iface, "/")
				ifaceName := ifaceParts[len(ifaceParts)-1]
				if ifaceName == depInterface && comp.Qualifier == dep.Qualifier {
					matches = append(matches, comp)
					break
				}
			}
		}
	}
	
	return matches
}

// findImplementationsForDependency finds components that implement an interface dependency
func (a *DependencyAnalyzer) findImplementationsForDependency(dep Dependency) []Component {
	var implementations []Component
	
	if !strings.Contains(dep.Type, ".") {
		return implementations
	}
	
	interfaceName := strings.Split(dep.Type, ".")[1]
	
	for _, comp := range a.components {
		for _, iface := range comp.Implements {
			if strings.HasSuffix(iface, interfaceName) && comp.Qualifier == dep.Qualifier {
				implementations = append(implementations, comp)
			}
		}
	}
	
	return implementations
}

// PrintAnalysisReport prints a comprehensive analysis report
func (a *DependencyAnalyzer) PrintAnalysisReport() {
	analysis := a.PerformComprehensiveAnalysis()
	
	fmt.Println("📊 Component Analysis Report")
	fmt.Println("============================")
	
	// Overview
	fmt.Printf("\n📋 Overview:\n")
	fmt.Printf("  Total Components: %d\n", analysis.TotalComponents)
	fmt.Printf("  Total Dependencies: %d\n", analysis.TotalDependencies)
	fmt.Printf("  Packages: %d\n", len(analysis.ComponentsByPackage))
	
	// Package breakdown
	fmt.Printf("\n📦 Components by Package:\n")
	var packageNames []string
	for pkg := range analysis.ComponentsByPackage {
		packageNames = append(packageNames, pkg)
	}
	sort.Strings(packageNames)
	
	for _, pkg := range packageNames {
		components := analysis.ComponentsByPackage[pkg]
		fmt.Printf("  %s: %d components\n", pkg, len(components))
	}
	
	// Circular dependencies
	if len(analysis.CircularDependencies) > 0 {
		fmt.Printf("\n🔄 Circular Dependencies (%d found):\n", len(analysis.CircularDependencies))
		for i, cycle := range analysis.CircularDependencies {
			fmt.Printf("  %d. %s\n", i+1, cycle.Description)
		}
	} else {
		fmt.Printf("\n✅ No circular dependencies found\n")
	}
	
	// Unused components
	if len(analysis.UnusedComponents) > 0 {
		fmt.Printf("\n🗑️  Unused Components (%d found):\n", len(analysis.UnusedComponents))
		for _, comp := range analysis.UnusedComponents {
			fmt.Printf("  - %s.%s [%s:%d]\n", comp.Package, comp.Type, comp.SourceFile, comp.LineNumber)
		}
	} else {
		fmt.Printf("\n✅ All components are used\n")
	}
	
	// Orphaned components
	if len(analysis.OrphanedComponents) > 0 {
		fmt.Printf("\n🔗 Orphaned Components (%d found):\n", len(analysis.OrphanedComponents))
		for _, comp := range analysis.OrphanedComponents {
			fmt.Printf("  - %s.%s [%s:%d]\n", comp.Package, comp.Type, comp.SourceFile, comp.LineNumber)
			fmt.Printf("    Has unsatisfied dependencies\n")
		}
	} else {
		fmt.Printf("\n✅ All component dependencies can be satisfied\n")
	}
	
	// Interface analysis
	fmt.Printf("\n🔌 Interface Analysis:\n")
	fmt.Printf("  Total Interfaces: %d\n", analysis.InterfaceAnalysis.TotalInterfaces)
	
	if len(analysis.InterfaceAnalysis.InterfacesWithMultiImpl) > 0 {
		fmt.Printf("  Interfaces with Multiple Implementations:\n")
		for iface, impls := range analysis.InterfaceAnalysis.InterfacesWithMultiImpl {
			fmt.Printf("    %s: %s\n", iface, strings.Join(impls, ", "))
		}
	}
	
	if len(analysis.InterfaceAnalysis.InterfacesWithNoImpl) > 0 {
		fmt.Printf("  Interfaces with No Implementation:\n")
		for _, iface := range analysis.InterfaceAnalysis.InterfacesWithNoImpl {
			fmt.Printf("    - %s\n", iface)
		}
	}
	
	// Qualifier conflicts
	if len(analysis.QualifierConflicts) > 0 {
		fmt.Printf("\n⚠️  Qualifier Conflicts (%d found):\n", len(analysis.QualifierConflicts))
		for _, conflict := range analysis.QualifierConflicts {
			fmt.Printf("  %s - Interface: %s, Qualifier: '%s'\n", 
				conflict.Severity, conflict.Interface, conflict.Qualifier)
			fmt.Printf("    Conflicting: %s\n", strings.Join(conflict.Conflicting, ", "))
		}
	} else {
		fmt.Printf("\n✅ No qualifier conflicts found\n")
	}
	
	// Dependency depth
	fmt.Printf("\n📏 Dependency Depth Analysis:\n")
	maxDepth := 0
	for _, depth := range analysis.DependencyDepth {
		if depth > maxDepth {
			maxDepth = depth
		}
	}
	fmt.Printf("  Maximum Dependency Depth: %d\n", maxDepth)
	
	// Components by depth
	depthGroups := make(map[int][]string)
	for comp, depth := range analysis.DependencyDepth {
		if depth >= 0 {
			depthGroups[depth] = append(depthGroups[depth], comp)
		}
	}
	
	for depth := 0; depth <= maxDepth; depth++ {
		if components, exists := depthGroups[depth]; exists {
			fmt.Printf("  Depth %d: %d components\n", depth, len(components))
		}
	}
}