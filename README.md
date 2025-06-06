# Go IoC: Bring "@Autowired" to Go!

Go IoC brings Spring-style autowiring to Go, offering a compile-time Inversion of Control (IoC) solution. Unlike Spring's runtime reflection-based `@Autowired`, Go IoC uses code generation to make dependencies explicit and type-safe at compile time, while keeping the familiar and intuitive struct tag syntax for dependency injection.

## Hasn't this been done already?

While other DI solutions exist in Go ([Google's Wire](https://github.com/google/wire/tree/main), [Uber's Dig](https://github.com/uber-go/dig) or [Facebook's Inject](https://github.com/facebookarchive/inject)), Go IoC takes a unique approach by:

- Providing a familiar Spring-like API that Java developers will recognize
- Using code generation for compile-time safety and zero runtime overhead (unlike Spring's runtime IoC)
- Using struct tags and marker structs as clean "annotations" 
- Supporting interface implementations and qualifiers elegantly
- Enabling automatic component scanning via struct marker
- Supporting lifecycle hooks via PostConstruct and PreDestroy struct methods
- **Advanced component analysis and debugging tools** including dependency graph visualization, circular dependency detection, and unused component identification

## But why bring @Autowired to Go?

For teams transitioning from Spring/Java to Go, especially those with significant Spring experience, dependency injection via `@Autowired` is more than just a familiar pattern — it's a proven productivity booster. Here's why I built Go IoC:

- **Smoother Java-to-Go Migration**: Many teams, including my team, come from a Spring background. Having a similar dependency injection pattern significantly reduces the learning curve and accelerates the transition to Go.

- **Proven Developer Experience**: Spring's autowiring has demonstrated for years that automatic dependency scanning and injection leads to:
  - Less boilerplate code for wiring dependencies
  - Cleaner, more maintainable code layers
  - Faster development cycles
  - Easier testing through dependency substitution

- **Best of Both Worlds**: While I love Spring's convenience, I also respect Go's philosophy of explicitness. That's why Go IoC:
  - Uses pure compile-time code generation
  - Makes dependencies traceable in generated code
  - Maintains compile-time type safety
  - Keeps the familiar developer experience without any runtime overhead

## How does it work?

Go IoC uses code generation to create a dependency injection system in three simple steps:

1. **Discovery**
   - Scans your code to find components and their dependencies
   - Identifies how components are connected through struct tags

2. **Analysis**
   - Figures out the correct order to create components
   - Validates that all dependencies can be satisfied

3. **Generation**
   - Generate type-safe code that initializes all components
   - Generate code to handle the component lifecycle (startup and cleanup)
   - Produces a single generated file that wires everything together

The result is a fast, type-safe dependency injection system with zero runtime reflection.

## IDE Support

### VS Code Extension

Get enhanced IDE support with the official Go IoC VS Code extension:

**[📦 Install from VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=keva-dev.go-ioc)**

**Features:**
- 🎨 **Visual Decorators**: IoC annotations highlighted with emojis (⚙️ Components, 🔗 Autowired, 🏷️ Qualifiers, 🔌 Interfaces)
- ⚠️ **Real-time Validation**: Live syntax checking for struct tags with actionable error messages
- 💡 **Smart Completions**: IntelliSense for IoC struct tags and annotations
- 🔍 **Hover Information**: Detailed documentation on IoC components and lifecycle methods
- 🔗 **Interface Navigation**: Ctrl+Click on interface names in `implements` tags to jump to definitions
- 🎯 **Commands**: Generate wire files, validate components, analyze dependencies, and show dependency graphs
- 🔧 **Code Snippets**: Quick snippets for creating components, adding dependencies, and lifecycle methods

Install via VS Code:
1. Open VS Code
2. Go to Extensions (Ctrl+Shift+X)
3. Search for "Go IoC"
4. Click Install

Or install via command line:
```bash
code --install-extension keva-dev.go-ioc
```

## Usage

To use Go IoC in your project, install the GO IoC code generator:

```bash
go install github.com/tuhuynh27/go-ioc/cmd/iocgen@latest
```

Please ensure that `$GOPATH/bin` is added to your $PATH.

### Defining Components

Components are defined using Go structs with specific annotations:

```go
// message/service.go
package message

type MessageService interface {
    SendMessage(msg string) string
}

type EmailService struct {
    Component struct{} // <- This is like a "@Component" annotation
    Implements struct{} `implements:"MessageService"` // <- Specify the interface that this struct implements
    Qualifier struct{} `value:"email"` // <- This is like a "@Qualifier" annotation
}

func NewEmailService() *EmailService {
    return &EmailService{}
}

func (s *EmailService) SendMessage(msg string) string {
    return fmt.Sprintf("Email: %s", msg)
}
```

### Defining Dependencies

Dependencies are defined using struct tags:

```go
// notification/service.go
package notification

type NotificationService struct {
    Component struct{}
    EmailSender message.MessageService `autowired:"true" qualifier:"email"` // <- This is like a "@Autowired" annotation
    SmsSender   message.MessageService `autowired:"true" qualifier:"sms"` // <- It also support "@Qualifier" annotation
}

func (s *NotificationService) SendNotifications(msg string) {
    s.EmailSender.SendMessage(msg)
    s.SmsSender.SendMessage(msg)
}
```

### Generating Initialization Code

Run the code generator in your project root:

```bash
iocgen
```

This will scan your project for components and generate the initialization code.

### Example Generated Code

Here's what the generated initialization code would look like:

```go
// wire/wire_gen.go
// Code generated by Go IoC. DO NOT EDIT.
//go:generate go run github.com/tuhuynh27/go-ioc/cmd/iocgen --dir=../
package wire

import (
    "your/project/message"
    "your/project/notification"
    "your/project/config"
)

type Container struct {
    EmailService        *message.EmailService
    SmsService          *message.SmsService
    NotificationService *notification.NotificationService
}

func Initialize() (*Container, func()) {
    container := &Container{}
    container.EmailService = message.NewEmailService()
    container.SmsService = &message.SmsService{}
    container.NotificationService = &notification.NotificationService{
        EmailSender: container.EmailService,
        SmsSender:   container.SmsService,
    }
    cleanup := func() {}
    return container, cleanup
}
```

### Using the Generated Code

Use the generated struct in your code:

```go
// main.go
package main

import (
    "your/project/wire"
)

func main() {
    container, cleanup := wire.Initialize()
    defer cleanup()
    // Get the service you need
    notificationService := container.NotificationService
    // Use it
    notificationService.SendNotifications("Hello World!")
}
```

## Component Discovery & Analysis

Go IoC provides powerful analysis tools to understand and optimize your component configurations:

### Component Listing

List all discovered components with detailed information:

```bash
iocgen --list
```

This shows:
- Component locations (file and line number)
- Dependencies and their qualifiers
- Interface implementations
- Lifecycle methods (PostConstruct/PreDestroy)
- Constructor functions

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

**🔍 Component Overview**
- Total components and dependencies
- Package distribution
- Component hierarchy depth

**🔄 Circular Dependency Detection**
- Identifies dependency cycles
- Shows complete dependency paths
- Suggests resolution strategies

**🗑️ Unused Component Detection**
- Finds components that aren't used as dependencies
- Helps identify potential cleanup opportunities
- Useful for removing dead code

**🔗 Orphaned Component Detection**
- Components with unsatisfied dependencies
- Missing interface implementations
- Configuration issues

**🔌 Interface Analysis**
- Multiple implementations per interface
- Interfaces with no implementations
- Qualifier conflicts and ambiguity

**📏 Dependency Depth Analysis**
- Calculates component initialization order
- Shows dependency hierarchy levels
- Identifies deeply nested dependencies

### Validation Without Generation

Validate your component configuration without generating files:

```bash
# Quick validation
iocgen --dry-run

# Detailed validation with warnings
iocgen --dry-run --verbose
```

### Dependency Graph Visualization

Visualize component relationships:

```bash
iocgen --graph
```

Shows a tree-like structure of all components and their dependencies, making it easy to understand the overall architecture.

## Comparison with Other DI Libraries

| Feature | Go IoC | Google Wire | Uber Dig | Facebook Inject |
|---------|-------|-------------|-----------|-----------------|
| Dependency Definition | Struct tags & marker structs | Function providers | Constructor functions | Struct tags |
| Runtime Overhead | None | None | Reflection-based | Reflection-based |
| Configuration Style | Spring-like annotations | Explicit provider functions | Constructor injection | Field tags |
| Interface Binding | Built-in | Manual provider setup | Manual provider setup | Limited support |
| Qualifier Support | Yes, via struct tags | No built-in support | Via name annotations | No |
| Learning Curve | Low (familiar to Spring devs) | Medium | Medium | Low |
| Code Generation | Yes | Yes | No | No |
| Compile-time Safety | Yes | Yes | Partial | No |
| Auto Component Scanning | Yes | No | No | No |
| Lifecycle Hooks | Yes | No | No | No |
| **Component Analysis** | **✅ Advanced** | **❌ None** | **❌ None** | **❌ None** |
| **Dependency Graph Visualization** | **✅ Built-in** | **❌ Manual** | **❌ Manual** | **❌ Manual** |
| **Circular Dependency Detection** | **✅ Automatic** | **⚠️ Build-time** | **⚠️ Runtime** | **❌ None** |
| **Unused Component Detection** | **✅ Yes** | **❌ No** | **❌ No** | **❌ No** |
| **Validation & Debugging** | **✅ Comprehensive** | **⚠️ Basic** | **⚠️ Basic** | **❌ None** |

## Test with Go IoC

Please check the [testing](docs/testing.md) guide for more information.

## Example

Please check the [example Git repository](https://github.com/tuhuynh27/go-ioc-gin-demo) (example with Go Gin web framework)

## Author's Note: Go Anti-Patterns and Migration Strategy

### ⚠️ Important Consideration

**This library intentionally violates Go idioms and best practices.** While Go IoC is technically functional and performant, it introduces several anti-patterns that go against Go's design philosophy:

**Anti-Patterns Used:**
- **Struct Embedding Abuse**: Empty struct fields as magic markers (`Component struct{}`, `Qualifier struct{}`)
- **Magic Behavior**: Heavy reliance on struct tags for core functionality instead of explicit code
- **Interface Pollution**: Upfront interface definitions rather than Go's "discover interfaces at point of use"
- **Global State Pattern**: Generated Container acting as a service locator
- **Convention over Configuration**: Java Spring naming conventions rather than Go idioms

### 🎯 Why Build This Despite Anti-Patterns?

**Migration Bridge for Java Teams**: This library serves as a **temporary migration tool** for teams transitioning from Java/Spring to Go. It provides:

1. **Familiar Patterns**: Reduces cognitive load for Spring developers learning Go
2. **Faster Initial Adoption**: Teams can be productive immediately without learning Go's dependency patterns
3. **Gradual Learning Curve**: Allows teams to focus on Go syntax and tooling before tackling Go's architectural patterns

### 🚀 Clean Migration Path

**The beauty of compile-time generation**: Since Go IoC leaves no runtime footprints, migration to idiomatic Go patterns is straightforward:

1. **Start with Go IoC**: Use familiar Spring-like patterns for initial development
2. **Learn Go Idioms**: Gradually understand Go's explicit dependency injection patterns
3. **Clean Removal**: Simply remove the marker structs and struct tags
4. **Refactor to Constructors**: Replace with explicit constructor functions

**Migration Example:**

```go
// Before (Go IoC)
type NotificationService struct {
    Component struct{}
    EmailSender message.MessageService `autowired:"true" qualifier:"email"`
    SmsSender   message.MessageService `autowired:"true" qualifier:"sms"`
}

// After (Idiomatic Go)
type NotificationService struct {
    emailSender message.MessageService
    smsSender   message.MessageService
}

func NewNotificationService(emailSender, smsSender message.MessageService) *NotificationService {
    return &NotificationService{
        emailSender: emailSender,
        smsSender:   smsSender,
    }
}
```

### 🎯 Recommendation

**For Production Go Applications**: Use explicit constructor patterns and Go's built-in dependency management.

**For Java Teams Learning Go**: Go IoC can serve as a stepping stone, but plan to migrate to idiomatic Go patterns as your team becomes comfortable with Go's philosophy of explicit, simple code.

**Remember**: The goal is not to make Go look like Java, but to help Java developers become productive Go developers faster.

## FAQ

### What's the performance impact?

None! Go IoC:
- Uses pure compile-time code generation, without runtime state or reflection, so not runtime cost
- Generates plain Go code that's as efficient as hand-written dependency injection
