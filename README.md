# Go IoC: Bring "@Autowired" to Go!

Go IoC brings Spring-style autowiring to Go, offering a compile-time Inversion of Control (IoC) solution. Unlike Spring's runtime reflection-based `@Autowired`, Go IoC uses code generation to make dependencies explicit and type-safe at compile time, while keeping the familiar and intuitive struct tag syntax for dependency injection.

## Hasn't this been done already?

While other DI solutions exist in Go ([Google's Wire](https://github.com/google/wire/tree/main), [Uber's Dig](https://github.com/uber-go/dig) or [Facebook's Inject](https://github.com/facebookarchive/inject)), Go IoC takes a unique approach by:

- Providing a familiar Spring-like API that Java developers will recognize
- Using code generation for compile-time safety (unlike Spring's runtime DI)
- Using struct tags and marker structs as clean "annotations" 
- Supporting interface implementations and qualifiers elegantly
- Enabling automatic component scanning via struct tags

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

Go IoC uses code generation to:

1. Scan struct tags at build time to identify components, interfaces, qualifiers and dependencies
2. Generate type-safe dependency injection code that creates and wires components correctly
3. Create initialization functions that handle all the wiring

## Usage

To use Go IoC in your project, install the code generator:

```bash
go install github.com/tuhuynh27/go-ioc/cmd/ioc-generate@latest
```

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
ioc-generate
```

This will scan your project for components and generate the initialization code.

### Example Generated Code

Here's what the generated initialization code would look like:

```go
// wire/wire_gen.go
// Code generated by Go IoC. DO NOT EDIT.
//go:generate go run github.com/tuhuynh27/go-ioc/cmd/ioc-generate -dir=../
package wire

import (
    "your/project/message"
    "your/project/notification"
    "your/project/config"
)

// Application holds all the wired components
type Application struct {
    emailService      *message.EmailService
    smsService       *message.SmsService
    notificationService *notification.NotificationService
}

// Initialize creates and wires all components
func Initialize() *Application {
    app := &Application{}
    
    app.emailService = &message.EmailService{
    }
    
    app.smsService = &message.SmsService{
    }
    
    app.notificationService = &notification.NotificationService{
        EmailSender: app.emailService,
        SmsSender:   app.smsService,
    }
    
    return app
}

// GetEmailService returns the EmailService instance
func (app *Application) GetEmailService() *message.EmailService {
    return app.emailService
}

// GetSmsService returns the SmsService instance
func (app *Application) GetSmsService() *message.SmsService {
    return app.smsService
}

// GetNotificationService returns the NotificationService instance
func (app *Application) GetNotificationService() *notification.NotificationService {
    return app.notificationService
}
```

### Using the Generated Code

Use the generated Application struct in your code:

```go
// main.go
package main

import (
    "your/project/wire"
)

func main() {
    app := wire.Initialize()
    
    // Get the service you need
    notificationService := app.GetNotificationService()
    
    // Use it
    notificationService.SendNotifications("Hello World!")
}
```

## Comparison with Other DI Libraries

| Feature | Go IoC | Google Wire | Uber Dig | Facebook Inject |
|---------|--------|-------------|-----------|-----------------|
| Dependency Definition | Struct tags & marker structs | Function providers | Constructor functions | Struct tags |
| Runtime Overhead | None | None | Reflection-based | Reflection-based |
| Configuration Style | Spring-like annotations | Explicit provider functions | Constructor injection | Field tags |
| Interface Binding | Built-in | Manual provider setup | Manual provider setup | Limited support |
| Qualifier Support | Yes, via struct tags | No built-in support | Via name annotations | No |
| Learning Curve | Low (familiar to Spring devs) | Medium | Medium | Low |
| Code Generation | Yes | Yes | No | No |
| Compile-time Safety | Yes | Yes | Partial | No |
| Auto Component Scanning | Yes (via struct tags) | No | No | No |

## Test with Go IoC

Please check the [testing](docs/testing.md) guide for more information.

## Demo

Please check the demo Git repository (example with Go Gin web framework) [here](https://github.com/tuhuynh27/go-ioc-gin-demo)

## FAQ

### What's the performance impact?

None! Go IoC:
- Uses pure compile-time code generation
- Has zero runtime overhead
- Generates plain Go code that's as efficient as hand-written dependency injection
- No reflection, no runtime container
