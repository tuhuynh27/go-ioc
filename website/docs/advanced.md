---
sidebar_position: 3
---

# Advanced Features

Explore Go IoC's powerful analysis and debugging capabilities.

## Component Analysis

Go IoC provides comprehensive analysis tools to understand and optimize your component configurations.

### Component Listing

List all discovered components with detailed information:

```bash
iocgen --list
```

Example output:
```
📦 notification.NotificationService [src/notification/service.go:15]
   🔗 Dependencies:
     - EmailSender: message.MessageService (qualifier: email)
     - SmsSender: message.MessageService (qualifier: sms)
     - Logger: logger.Logger (qualifier: json)

📦 message.EmailService (qualifier: email) [src/message/email.go:8]
   📋 Implements: message.MessageService
   🔗 Dependencies:
     - Config: config.Config
     - Logger: logger.Logger (qualifier: console)
   🚀 Has PostConstruct method
```

### Comprehensive Analysis

Perform in-depth analysis of your component architecture:

```bash
iocgen --analyze
```

The analysis report includes:

#### 🔍 Component Overview
- Total components and dependencies
- Package distribution
- Component hierarchy depth

#### 🔄 Circular Dependency Detection
- Identifies dependency cycles
- Shows complete dependency paths
- Suggests resolution strategies

#### 🗑️ Unused Component Detection
- Finds components that aren't used as dependencies
- Helps identify potential cleanup opportunities
- Useful for removing dead code

#### 🔗 Orphaned Component Detection
- Components with unsatisfied dependencies
- Missing interface implementations
- Configuration issues

#### 🔌 Interface Analysis
- Multiple implementations per interface
- Interfaces with no implementations
- Qualifier conflicts and ambiguity

#### 📏 Dependency Depth Analysis
- Calculates component initialization order
- Shows dependency hierarchy levels
- Identifies deeply nested dependencies

## Dependency Graph Visualization

Visualize component relationships:

```bash
iocgen --graph
```

Shows a tree-like structure of all components and their dependencies, making it easy to understand the overall architecture.

Example output:
```
📦 Components Dependency Graph:
├── notification.NotificationService
│   ├── 🔗 EmailSender: message.MessageService (qualifier: email)
│   │   └── → message.EmailService
│   ├── 🔗 SmsSender: message.MessageService (qualifier: sms)
│   │   └── → message.SmsService
│   └── 🔗 Logger: logger.Logger (qualifier: json)
│       └── → logger.JsonLogger
├── message.EmailService
│   ├── 🔗 Config: config.Config
│   │   └── → config.AppConfig
│   └── 🔗 Logger: logger.Logger (qualifier: console)
│       └── → logger.ConsoleLogger
```

## Validation and Debugging

### Dry Run Validation

Validate your component configuration without generating files:

```bash
# Quick validation
iocgen --dry-run

# Detailed validation with warnings
iocgen --dry-run --verbose
```

### Enhanced Error Messages

Go IoC provides detailed error messages with source location information:

```
Error: Failed to resolve dependency
  Component: notification.NotificationService
  File: src/notification/service.go:15
  Field: EmailSender (message.MessageService)
  Qualifier: email
  
  Suggestion: Ensure you have a component implementing MessageService with qualifier "email"
```

### Common Issues and Solutions

#### Missing Interface Implementation

**Error**: Interface has no implementations
```go
// Problem: No component implements UserService
type UserService interface {
    GetUser(id string) (*User, error)
}
```

**Solution**: Add a component that implements the interface
```go
type UserServiceImpl struct {
    Component struct{}
    Implements struct{} `implements:"UserService"`
}
```

#### Qualifier Conflicts

**Error**: Multiple components with same qualifier
```go
// Problem: Both components have qualifier "primary"
type ServiceA struct {
    Component struct{}
    Qualifier struct{} `value:"primary"`
}

type ServiceB struct {
    Component struct{}
    Qualifier struct{} `value:"primary"`
}
```

**Solution**: Use unique qualifiers
```go
type ServiceA struct {
    Component struct{}
    Qualifier struct{} `value:"service-a"`
}

type ServiceB struct {
    Component struct{}
    Qualifier struct{} `value:"service-b"`
}
```

#### Circular Dependencies

**Error**: Circular dependency detected
```go
// Problem: A depends on B, B depends on A
type ServiceA struct {
    Component struct{}
    ServiceB ServiceBInterface `autowired:"true"`
}

type ServiceB struct {
    Component struct{}
    ServiceA ServiceAInterface `autowired:"true"`
}
```

**Solution**: Refactor to remove circular dependency
```go
// Extract common functionality to a shared service
type SharedService struct {
    Component struct{}
}

type ServiceA struct {
    Component struct{}
    Shared SharedServiceInterface `autowired:"true"`
}

type ServiceB struct {
    Component struct{}
    Shared SharedServiceInterface `autowired:"true"`
}
```

## Code Generation Options

### Custom Output Directory

Generate wire files in a specific directory:

```bash
iocgen --dir=./internal/wire
```

### Custom Output File

Specify a custom output file name:

```bash
iocgen --output=custom_wire.go
```

### Verbose Output

Get detailed information during generation:

```bash
iocgen --verbose
```

## Integration with Build Tools

### Go Generate

Add a go:generate directive to automatically run iocgen:

```go
//go:generate go run github.com/tuhuynh27/go-ioc/cmd/iocgen
```

Then run:
```bash
go generate ./...
```

### Makefile Integration

```makefile
.PHONY: generate
generate:
	iocgen

.PHONY: validate
validate:
	iocgen --dry-run --verbose

.PHONY: analyze
analyze:
	iocgen --analyze

.PHONY: graph
graph:
	iocgen --graph
```

### CI/CD Integration

Add validation to your CI pipeline:

```yaml
# GitHub Actions example
- name: Validate IoC Configuration
  run: |
    go install github.com/tuhuynh27/go-ioc/cmd/iocgen@latest
    iocgen --dry-run --verbose
```

## Performance Considerations

### Generated Code Optimization

Go IoC generates optimized initialization code:

```go
// Generated code is direct and efficient
func Initialize() (*Container, func()) {
    container := &Container{}
    
    // Direct instantiation - no reflection
    container.ConfigService = &config.ConfigService{}
    container.LoggerService = logger.NewLoggerService()
    
    // Direct field assignment - compile-time safe
    container.UserService = &user.UserService{
        Config: container.ConfigService,
        Logger: container.LoggerService,
    }
    
    // Lifecycle methods called directly
    container.ConfigService.PostConstruct()
    container.LoggerService.PostConstruct()
    
    cleanup := func() {
        container.UserService.PreDestroy()
        container.LoggerService.PreDestroy()
        container.ConfigService.PreDestroy()
    }
    
    return container, cleanup
}
```

### Build Time Impact

- Component scanning: Minimal impact (AST parsing)
- Code generation: Fast (direct Go code generation)
- No runtime overhead: Pure compile-time solution

## Comparison with Other Solutions

| Feature | Go IoC | Google Wire | Uber Dig |
|---------|--------|-------------|----------|
| Runtime Overhead | None | None | Reflection-based |
| Configuration Style | Annotations | Provider functions | Constructor injection |
| Interface Binding | Built-in | Manual setup | Manual setup |
| Qualifier Support | Yes | No | Via annotations |
| Component Scanning | Automatic | Manual | Manual |
| Lifecycle Hooks | Yes | No | Limited |
| Dependency Analysis | Advanced | Basic | Basic |
| Circular Detection | Automatic | Build-time | Runtime |
| Graph Visualization | Built-in | Manual | Manual |