---
sidebar_position: 1
---

# Getting Started

Let's get **Go IoC up and running in less than 5 minutes**.

Go IoC brings Spring-style autowiring to Go, offering a compile-time Inversion of Control (IoC) solution. Unlike Spring's runtime reflection-based `@Autowired`, Go IoC uses code generation to make dependencies explicit and type-safe at compile time.

## Installation

Install the Go IoC code generator:

```bash
go install github.com/tuhuynh27/go-ioc/cmd/iocgen@latest
```

Make sure `$GOPATH/bin` is added to your `$PATH`.

## Quick Example

### 1. Define a Service Interface

```go
// message/service.go
package message

type MessageService interface {
    SendMessage(msg string) string
}
```

### 2. Create a Component

```go
// message/email_service.go
package message

import "fmt"

type EmailService struct {
    Component struct{} // <- Component marker
    Implements struct{} `implements:"MessageService"` // <- Interface implementation
    Qualifier struct{} `value:"email"` // <- Qualifier for disambiguation
}

func (s *EmailService) SendMessage(msg string) string {
    return fmt.Sprintf("Email: %s", msg)
}
```

### 3. Use Dependency Injection

```go
// notification/service.go
package notification

type NotificationService struct {
    Component struct{} // <- Component marker
    EmailSender message.MessageService `autowired:"true" qualifier:"email"` // <- Autowired dependency
}

func (s *NotificationService) SendNotifications(msg string) {
    s.EmailSender.SendMessage(msg)
}
```

### 4. Generate Wire Code

Run the code generator in your project root:

```bash
iocgen
```

### 5. Use the Generated Container

```go
// main.go
package main

import "your/project/wire"

func main() {
    container, cleanup := wire.Initialize()
    defer cleanup()
    
    notificationService := container.NotificationService
    notificationService.SendNotifications("Hello World!")
}
```

## Key Features

- **🍃 Spring-like Syntax**: Familiar `@Autowired` and `@Component` annotations via struct tags
- **⚡ Zero Runtime Overhead**: Pure compile-time code generation, no reflection
- **🔒 Type Safety**: All dependencies resolved and validated at compile time
- **🔍 Advanced Analysis**: Built-in dependency graph visualization and circular dependency detection
- **🎯 Interface Support**: Clean interface binding with qualifier support
- **♻️ Lifecycle Hooks**: PostConstruct and PreDestroy method support

## Next Steps

- Learn more about [Component Definition](./components)
- Explore [Advanced Features](./advanced)
- Try the [VS Code Extension](./vscode-extension)
- Check out the [Testing Guide](./testing)
