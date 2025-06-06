---
sidebar_position: 2
---

# Component Definition

Learn how to define and configure IoC components using Go IoC's Spring-like syntax.

## Component Markers

Go IoC uses struct markers to identify components. These are empty struct fields that serve as "annotations":

```go
type MyService struct {
    Component struct{} // Marks this struct as an IoC component
}
```

## Interface Implementation

Components can implement interfaces using the `Implements` marker:

```go
type UserService interface {
    GetUser(id string) (*User, error)
}

type UserServiceImpl struct {
    Component struct{}
    Implements struct{} `implements:"UserService"`
}

func (s *UserServiceImpl) GetUser(id string) (*User, error) {
    // Implementation here
    return &User{ID: id}, nil
}
```

## Qualifiers

When multiple components implement the same interface, use qualifiers to distinguish them:

```go
type EmailService struct {
    Component struct{}
    Implements struct{} `implements:"MessageService"`
    Qualifier struct{} `value:"email"` // Qualifier for this implementation
}

type SmsService struct {
    Component struct{}
    Implements struct{} `implements:"MessageService"`
    Qualifier struct{} `value:"sms"` // Different qualifier
}
```

## Dependency Injection

Use the `autowired` tag to inject dependencies:

```go
type NotificationService struct {
    Component struct{}
    
    // Inject by interface type
    Logger logger.Logger `autowired:"true"`
    
    // Inject with qualifier
    EmailSender message.MessageService `autowired:"true" qualifier:"email"`
    SmsSender   message.MessageService `autowired:"true" qualifier:"sms"`
}
```

## Constructor Functions

Go IoC can detect and use constructor functions automatically:

```go
type DatabaseService struct {
    Component struct{}
    connectionString string
}

// Constructor function - will be used for initialization
func NewDatabaseService() *DatabaseService {
    return &DatabaseService{
        connectionString: "default-connection",
    }
}
```

## Lifecycle Methods

Components can define lifecycle methods for initialization and cleanup:

```go
type CacheService struct {
    Component struct{}
    cache map[string]interface{}
}

// Called after dependency injection is complete
func (s *CacheService) PostConstruct() error {
    s.cache = make(map[string]interface{})
    fmt.Println("Cache initialized")
    return nil
}

// Called during application shutdown
func (s *CacheService) PreDestroy() error {
    s.cache = nil
    fmt.Println("Cache cleaned up")
    return nil
}
```

## Complete Example

Here's a comprehensive example showing all features:

```go
package service

type Logger interface {
    Log(message string)
}

type ConsoleLogger struct {
    Component struct{}
    Implements struct{} `implements:"Logger"`
    Qualifier struct{} `value:"console"`
}

func (l *ConsoleLogger) Log(message string) {
    fmt.Println("CONSOLE:", message)
}

type FileLogger struct {
    Component struct{}
    Implements struct{} `implements:"Logger"`
    Qualifier struct{} `value:"file"`
    filename string
}

func NewFileLogger() *FileLogger {
    return &FileLogger{filename: "app.log"}
}

func (l *FileLogger) Log(message string) {
    // Write to file implementation
}

func (l *FileLogger) PostConstruct() error {
    fmt.Println("File logger initialized with file:", l.filename)
    return nil
}

type UserService struct {
    Component struct{}
    
    // Dependencies
    ConsoleLogger Logger `autowired:"true" qualifier:"console"`
    FileLogger    Logger `autowired:"true" qualifier:"file"`
    Database      DatabaseInterface `autowired:"true"`
}

func (s *UserService) CreateUser(name string) {
    s.ConsoleLogger.Log("Creating user: " + name)
    s.FileLogger.Log("User created: " + name)
    // Database operations...
}
```

## Best Practices

### 1. Use Interfaces

Always define interfaces for your services to enable loose coupling:

```go
// Good
type UserService interface {
    CreateUser(name string) error
    GetUser(id string) (*User, error)
}

type UserServiceImpl struct {
    Component struct{}
    Implements struct{} `implements:"UserService"`
}
```

### 2. Meaningful Qualifiers

Use descriptive qualifiers that clearly indicate the purpose:

```go
// Good
Qualifier struct{} `value:"primary-database"`
Qualifier struct{} `value:"cache-redis"`
Qualifier struct{} `value:"smtp-email"`

// Avoid generic names
Qualifier struct{} `value:"impl1"`
Qualifier struct{} `value:"service"`
```

### 3. Constructor Functions

Use constructor functions for complex initialization:

```go
func NewDatabaseService(config *Config) *DatabaseService {
    return &DatabaseService{
        connectionString: config.DatabaseURL,
        maxConnections:   config.MaxConnections,
    }
}
```

### 4. Lifecycle Methods

Use lifecycle methods for resource management:

```go
func (s *DatabaseService) PostConstruct() error {
    return s.connect()
}

func (s *DatabaseService) PreDestroy() error {
    return s.disconnect()
}
```

## Validation

Go IoC provides comprehensive validation of your component configuration:

```bash
# Validate without generating files
iocgen --dry-run

# Show detailed validation information
iocgen --dry-run --verbose

# List all discovered components
iocgen --list
```

Common validation checks include:
- Missing interface implementations
- Unresolved dependencies
- Circular dependencies
- Qualifier conflicts
- Invalid struct tag syntax