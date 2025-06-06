package wire

import (
	"testing"
)

func TestAnalyzer_FindCircularDependencies(t *testing.T) {
	// Create test components with circular dependency
	components := []Component{
		{
			Name:    "ServiceA",
			Type:    "ServiceA",
			Package: "test",
			Dependencies: []Dependency{
				{FieldName: "ServiceB", Type: "ServiceB", Qualifier: ""},
			},
		},
		{
			Name:    "ServiceB",
			Type:    "ServiceB",
			Package: "test",
			Dependencies: []Dependency{
				{FieldName: "ServiceA", Type: "ServiceA", Qualifier: ""},
			},
		},
	}

	analyzer := NewAnalyzer(components)
	cycles := analyzer.FindCircularDependencies()

	if len(cycles) == 0 {
		t.Error("Expected to find circular dependency, but none found")
	}

	if len(cycles) > 0 {
		cycle := cycles[0]
		if len(cycle.Path) < 2 {
			t.Errorf("Expected cycle path length >= 2, got %d", len(cycle.Path))
		}
	}
}

func TestAnalyzer_FindUnusedComponents(t *testing.T) {
	// Create test components where one is unused
	components := []Component{
		{
			Name:         "UsedService",
			Type:         "UsedService",
			Package:      "test",
			Dependencies: []Dependency{},
		},
		{
			Name:    "UsingService",
			Type:    "UsingService",
			Package: "test",
			Dependencies: []Dependency{
				{FieldName: "Used", Type: "UsedService", Qualifier: ""}, // Simple type reference
			},
		},
		{
			Name:    "RootService", // This represents a "root" component like a controller
			Type:    "RootService",
			Package: "test",
			Dependencies: []Dependency{
				{FieldName: "Service", Type: "UsingService", Qualifier: ""},
			},
		},
		{
			Name:         "UnusedService",
			Type:         "UnusedService",
			Package:      "test",
			Dependencies: []Dependency{},
		},
	}

	analyzer := NewAnalyzer(components)
	unused := analyzer.FindUnusedComponents()

	// Both RootService and UnusedService should be unused (not dependencies of other components)
	if len(unused) != 2 {
		t.Errorf("Expected 2 unused components, got %d", len(unused))
		for i, comp := range unused {
			t.Logf("Unused[%d]: %s.%s", i, comp.Package, comp.Type)
		}
	}

	// Check that both expected components are in the unused list
	unusedTypes := make(map[string]bool)
	for _, comp := range unused {
		unusedTypes[comp.Type] = true
	}

	if !unusedTypes["RootService"] {
		t.Error("Expected RootService to be unused")
	}
	if !unusedTypes["UnusedService"] {
		t.Error("Expected UnusedService to be unused")
	}
}

func TestAnalyzer_AnalyzeInterfaces(t *testing.T) {
	// Create test components with interface implementations
	components := []Component{
		{
			Name:       "ServiceA",
			Type:       "ServiceA",
			Package:    "test.ServiceA",
			Implements: []string{"test.Service"},
		},
		{
			Name:       "ServiceB",
			Type:       "ServiceB",
			Package:    "test.ServiceB",
			Implements: []string{"test.Service"},
		},
		{
			Name:    "Consumer",
			Type:    "Consumer",
			Package: "test.Consumer",
			Dependencies: []Dependency{
				{FieldName: "Service", Type: "test.Service", Qualifier: ""},
			},
		},
	}

	analyzer := NewAnalyzer(components)
	analysis := analyzer.AnalyzeInterfaces()

	if analysis.TotalInterfaces != 1 {
		t.Errorf("Expected 1 interface, got %d", analysis.TotalInterfaces)
	}

	if len(analysis.InterfacesWithMultiImpl) != 1 {
		t.Errorf("Expected 1 interface with multiple implementations, got %d", len(analysis.InterfacesWithMultiImpl))
	}

	if implementations, exists := analysis.InterfacesWithMultiImpl["test.Service"]; !exists {
		t.Error("Expected test.Service to have multiple implementations")
	} else if len(implementations) != 2 {
		t.Errorf("Expected 2 implementations for test.Service, got %d", len(implementations))
	}
}

func TestAnalyzer_FindQualifierConflicts(t *testing.T) {
	// Create test components with qualifier conflicts
	components := []Component{
		{
			Name:       "ServiceA",
			Type:       "ServiceA",
			Package:    "test.ServiceA",
			Qualifier:  "default",
			Implements: []string{"test.Service"},
		},
		{
			Name:       "ServiceB",
			Type:       "ServiceB",
			Package:    "test.ServiceB",
			Qualifier:  "default", // Same qualifier - conflict!
			Implements: []string{"test.Service"},
		},
	}

	analyzer := NewAnalyzer(components)
	conflicts := analyzer.FindQualifierConflicts()

	if len(conflicts) != 1 {
		t.Errorf("Expected 1 qualifier conflict, got %d", len(conflicts))
	}

	if len(conflicts) > 0 {
		conflict := conflicts[0]
		if conflict.Interface != "test.Service" {
			t.Errorf("Expected conflict for test.Service, got %s", conflict.Interface)
		}
		if conflict.Qualifier != "default" {
			t.Errorf("Expected conflict for qualifier 'default', got '%s'", conflict.Qualifier)
		}
		if len(conflict.Conflicting) != 2 {
			t.Errorf("Expected 2 conflicting components, got %d", len(conflict.Conflicting))
		}
	}
}

func TestAnalyzer_CalculateDependencyDepth(t *testing.T) {
	// Create test components with different dependency depths
	components := []Component{
		{
			Name:         "Level0",
			Type:         "Level0",
			Package:      "test",
			Dependencies: []Dependency{}, // No dependencies - depth 0
		},
		{
			Name:    "Level1",
			Type:    "Level1",
			Package: "test",
			Dependencies: []Dependency{
				{FieldName: "Dep", Type: "Level0", Qualifier: ""},
			}, // Depends on Level0 - depth 1
		},
		{
			Name:    "Level2",
			Type:    "Level2",
			Package: "test",
			Dependencies: []Dependency{
				{FieldName: "Dep", Type: "Level1", Qualifier: ""},
			}, // Depends on Level1 - depth 2
		},
	}

	analyzer := NewAnalyzer(components)
	depths := analyzer.CalculateDependencyDepth()

	if depths["test.Level0"] != 0 {
		t.Errorf("Expected Level0 depth 0, got %d", depths["test.Level0"])
	}
	if depths["test.Level1"] != 1 {
		t.Errorf("Expected Level1 depth 1, got %d", depths["test.Level1"])
	}
	if depths["test.Level2"] != 2 {
		t.Errorf("Expected Level2 depth 2, got %d", depths["test.Level2"])
	}
}

func TestAnalyzer_PerformComprehensiveAnalysis(t *testing.T) {
	// Create a comprehensive test case
	components := []Component{
		{
			Name:         "ConfigService",
			Type:         "ConfigService",
			Package:      "test",
			Dependencies: []Dependency{},
		},
		{
			Name:    "DatabaseService",
			Type:    "DatabaseService",
			Package: "test",
			Dependencies: []Dependency{
				{FieldName: "Config", Type: "ConfigService", Qualifier: ""},
			},
		},
		{
			Name:         "UnusedService",
			Type:         "UnusedService",
			Package:      "test",
			Dependencies: []Dependency{},
		},
	}

	analyzer := NewAnalyzer(components)
	result := analyzer.PerformComprehensiveAnalysis()

	if result.TotalComponents != 3 {
		t.Errorf("Expected 3 total components, got %d", result.TotalComponents)
	}

	if result.TotalDependencies != 1 {
		t.Errorf("Expected 1 total dependency, got %d", result.TotalDependencies)
	}

	// Both DatabaseService and UnusedService should be unused (not dependencies of other components)
	if len(result.UnusedComponents) != 2 {
		t.Errorf("Expected 2 unused components, got %d", len(result.UnusedComponents))
	}

	if len(result.CircularDependencies) != 0 {
		t.Errorf("Expected 0 circular dependencies, got %d", len(result.CircularDependencies))
	}

	if len(result.ComponentsByPackage) != 1 {
		t.Errorf("Expected 1 package, got %d", len(result.ComponentsByPackage))
	}
}