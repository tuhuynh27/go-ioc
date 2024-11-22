# Go IoC

Go IoC is a lightweight, compile-time Inversion of Control (IoC) container for Go (inspired by Spring's APIs), designed to simplify dependency injection and object management in Go applications.

## Hasn't this been done already?

While there are other dependency injection solutions for Go, Go IoC takes a unique approach by:

- Providing compile-time dependency injection for better performance and safety
- Using code generation to avoid runtime reflection
- Providing a familiar Spring-like API that Java developers will recognize
- Using struct tags as "annotations" to keep the syntax clean and idiomatic to Go
- Supporting interface implementations and qualifiers in an elegant way
- Avoiding complex configuration files or setup - just use struct tags

## How does it work?

Go IoC uses code generation to:

1. Scan struct tags during build time to identify:
   - Components via the `Component` struct embedding
   - Interface implementations via `implements` tags
   - Qualifiers via `value` tags
   - Dependencies via `autowired` tags

2. Generate type-safe dependency injection code that:
   - Creates components in the correct order
   - Injects dependencies properly
   - Handles interface implementations and qualifiers
   - Provides compile-time safety

3. Create a wire_gen.go file that contains all the initialization logic

The container maintains internal maps to track components and their implementations, making dependency lookup fast and reliable.

## Features

- Compile-time Dependency Injection: No runtime reflection for better performance
- Type Safety: Errors are caught at compile time, not runtime
- Code Generation: Automatically generates initialization code
- Component Registration: Register components with annotations for easy management
- Interface Implementation: Support for multiple implementations of the same interface
- Qualifiers: Use qualifiers to distinguish between different implementations

## Usage

To use Go IoC in your project:

```bash
go get github.com/tuhuynh27/go-ioc@latest
go install github.com/tuhuynh27/go-ioc/cmd/generate@latest
```

### Defining Components

Components are defined using Go structs with specific annotations:

```go
type MessageService interface {
    SendMessage(msg string) string
}

type EmailService struct {
    ioc.Component
    Logger logger.Logger `autowired:""`
}

func (s *EmailService) SendMessage(msg string) string {
    s.Logger.Log("Sending email: " + msg)
    return fmt.Sprintf("Email: %s", msg)
}
```

### Defining Dependencies

Dependencies are defined using struct tags:

```go
type NotificationService struct {
    ioc.Component
    EmailSender message.MessageService `autowired:"" qualifier:"email"`
    SmsSender   message.MessageService `autowired:"" qualifier:"sms"`
}

func (s *NotificationService) SendNotifications(msg string) {
    s.EmailSender.SendMessage(msg)
    s.SmsSender.SendMessage(msg)
}
```

### Generating Container Code

Run the code generator to create the dependency injection container:

```bash
go generate ./wire
```

This will generate a `wire/wire_gen.go` file in your project directory.

### Using the Container

Use the generated container in your application:

```go
func main() {
    container, err := wire.InitializeContainer()
    if err != nil {
        panic(err)
    }

    // Get components by name
    notificationService := container.Get("NotificationService")

    // Get components by interface with qualifier
    emailSender := container.GetQualified("MessageService", "email")
    smsSender := container.GetQualified("MessageService", "sms")
}
```

## Demo

Please check the demo Git repository [here](https://github.com/tuhuynh27/go-ioc-gin-demo)
