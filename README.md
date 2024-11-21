# Go IoC

Go IoC is a lightweight Inversion of Control (IoC) container for Go (inspired by Spring's APIs), designed to simplify dependency injection and object management in Go applications.

## Hasn't this been done already?

While there are other dependency injection solutions for Go, Go IoC takes a unique approach by:

- Providing a familiar Spring-like API that Java developers will recognize
- Using struct tags as "annotations" to keep the syntax clean and idiomatic to Go
- Focusing on simplicity and ease of use while maintaining flexibility
- Supporting interface implementations and qualifiers in an elegant way
- Avoiding complex configuration files or setup - just use struct tags

## How does it work?

Go IoC uses Go's reflection capabilities to:

1. Scan struct tags during component registration to identify:
   - Components via the `Component` struct embedding
   - Interface implementations via `implements` tags
   - Qualifiers via `value` tags
   - Dependencies via `autowired` tags

2. Build a dependency graph by analyzing the relationships between components

3. Resolve and inject dependencies by:
   - Looking up components by name or interface
   - Matching qualifiers when multiple implementations exist
   - Using reflection to set field values

The container maintains internal maps to track components, their definitions, and interface implementations, making dependency resolution fast and reliable.

## Features

- Component Registration: Register components with annotations for easy management.
- Dependency Injection: Automatically inject dependencies into components.
- Interface Implementation: Support for multiple implementations of the same interface.
- Qualifiers: Use qualifiers to distinguish between different implementations of the same interface.

## Usage

To use Keva Go IoC in your project, you need first to download it:

```
go get github.com/tuhuynh27/go-ioc
```

Next, import the package in your Go code:

```
import "github.com/tuhuynh27/go-ioc/ioc"
```

### Defining Components

Components are defined using Go structs with specific annotations. Here's an example of a component definition:

```go
type MessageService interface {
	SendMessage(msg string) string
}

type EmailService struct {
	Component  ioc.Component
	Qualifier  struct{}      `value:"email"`
	Implements struct{}      `implements:"MessageService"`
}

func (s *EmailService) SendMessage(msg string) string {
    return fmt.Sprintf("Email: %s", msg)
}
```

- `Component`: Marks the struct as a component.
- `Qualifier`: Specifies a unique identifier for the component.
- `Implements`: Indicates the interface(s) the component implements.

### Defining dependencies

Dependencies are defined using Go structs with specific annotations. Here's an example of a dependency definition:

```go
type NotificationService struct {
	Component   ioc.Component
	EmailSender message.MessageService `autowired:"true" qualifier:"email"`
}

func (s *NotificationService) SendNotifications(msg string) {
	fmt.Println(s.EmailSender.SendMessage(msg))
}
```

- `Component`: Marks the struct as a component.
- `Autowired`: Marks a field for dependency injection.
- `Qualifier`: Specifies a unique identifier for the dependency.

### Registering Components

To register components, create a new container and use the `RegisterComponents` method:

```go
container := ioc.NewContainer()

err := container.RegisterComponents(
	&message.EmailService{},
	&notification.NotificationService{},
)

if err != nil {
	panic(err)
}
```

### Using Components

Once components are registered, you can retrieve and use them from the container:

```go
notificationService := 
    container.Get("NotificationService")
        .(*notification.NotificationService)

notificationService.SendNotifications("Hello, World!")
```
