# Testing with Go IoC

This guide demonstrates how to effectively test applications using Go IoC, covering both unit testing of individual components and integration testing.

## Unit Testing Components

### Example Component Structure

Let's use the message service example from our main documentation:

```go:message/service.go
type MessageService interface {
    SendMessage(msg string) string
}

type EmailService struct {
    Component struct{} `implements:"MessageService"`
    Qualifier struct{} `value:"email"`
}

func (s *EmailService) SendMessage(msg string) string {
    return fmt.Sprintf("Email: %s", msg)
}
```

### Creating Mock Components

Create mock implementations using the same struct tags:

```go:message/mock_service.go
type MockMessageService struct {
    Component struct{} `implements:"MessageService"`
    Qualifier struct{} `value:"mock"`
    
    // Add fields to track calls
    Messages []string
}

func (s *MockMessageService) SendMessage(msg string) string {
    s.Messages = append(s.Messages, msg)
    return fmt.Sprintf("Mock: %s", msg)
}
```

### Unit Test Example

```go:message/service_test.go
func TestNotificationService(t *testing.T) {
    // Create a test-specific application configuration
    mockService := &MockMessageService{
        Messages: make([]string, 0),
    }
    
    // Create the service under test with the mock dependency
    notificationService := &notification.NotificationService{
        EmailSender: mockService,
    }
    
    // Run the test
    testMsg := "Hello Test"
    notificationService.SendNotifications(testMsg)
    
    // Verify the mock was called correctly
    if len(mockService.Messages) != 1 {
        t.Errorf("Expected 1 message, got %d", len(mockService.Messages))
    }
    if mockService.Messages[0] != testMsg {
        t.Errorf("Expected message '%s', got '%s'", testMsg, mockService.Messages[0])
    }
}
```

## Integration Testing

For integration tests, you can create a test-specific wire configuration:

### Test Configuration

```go:wire/wire_test.go
// TestApplication extends the main Application with test-specific components
type TestApplication struct {
    *Application
    mockEmailService    *message.MockMessageService
    mockSmsService      *message.MockMessageService
}

// InitializeTestApplication creates a test application with mock services
func InitializeTestApplication() *TestApplication {
    app := &TestApplication{}
    
    // Initialize mock services
    app.mockEmailService = &message.MockMessageService{
        Messages: make([]string, 0),
    }
    
    app.mockSmsService = &message.MockMessageService{
        Messages: make([]string, 0),
    }
    
    // Initialize the notification service with mocks
    app.notificationService = &notification.NotificationService{
        EmailSender: app.mockEmailService,
        SmsSender:   app.mockSmsService,
    }
    
    return app
}
```

### Integration Test Example

```go:tests/integration_test.go
func TestNotificationIntegration(t *testing.T) {
    // Initialize the test application
    app := wire.InitializeTestApplication()
    
    // Get the service
    notificationService := app.GetNotificationService()
    
    // Run the integration test
    testMsg := "Integration Test"
    notificationService.SendNotifications(testMsg)
    
    // Verify both mock services were called
    mockEmail := app.mockEmailService
    mockSms := app.mockSmsService
    
    if len(mockEmail.Messages) != 1 {
        t.Errorf("Expected 1 email, got %d", len(mockEmail.Messages))
    }
    
    if len(mockSms.Messages) != 1 {
        t.Errorf("Expected 1 SMS, got %d", len(mockSms.Messages))
    }
}
```

## Testing Best Practices

1. **Mock Interface Creation**
   - Create mock implementations using the same struct tags as real components
   - Add fields to track calls, arguments, and control return values
   - Consider using a mocking library like `testify/mock` for more complex scenarios

2. **Test Configuration**
   - Create separate initialization functions for tests
   - Use qualifier tags to distinguish between production and test components
   - Keep test wire configurations in separate files

3. **Integration Testing**
   - Create a test-specific Application struct that extends the main one
   - Initialize with mock dependencies as needed
   - Use the same wire pattern as production code

4. **Test Organization**
   - Keep unit tests close to the components they test
   - Place integration tests in a separate `tests` package
   - Use table-driven tests for comprehensive coverage

## Example with External Dependencies

When testing components with external dependencies (like databases or APIs), you can create mock implementations:

```go:repository/mock_repository.go
type MockUserRepository struct {
    Component struct{} `implements:"UserRepository"`
    Qualifier struct{} `value:"mock"`
    
    Users map[string]User
}

func (r *MockUserRepository) GetUser(id string) (User, error) {
    if user, exists := r.Users[id]; exists {
        return user, nil
    }
    return User{}, errors.New("user not found")
}
```

Then use it in your tests:

```go:service/user_service_test.go
func TestUserService(t *testing.T) {
    // Create mock repository with test data
    mockRepo := &repository.MockUserRepository{
        Users: map[string]User{
            "123": {ID: "123", Name: "Test User"},
        },
    }
    
    // Initialize service with mock
    userService := &UserService{
        Repository: mockRepo,
    }
    
    // Run tests
    user, err := userService.GetUser("123")
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    if user.Name != "Test User" {
        t.Errorf("Expected 'Test User', got '%s'", user.Name)
    }
}
```

## Using Test Containers

For integration tests that require real databases or external services, consider using [testcontainers-go](https://github.com/testcontainers/testcontainers-go):

```go:tests/database_test.go
func TestWithDatabase(t *testing.T) {
    // Start a PostgreSQL container
    postgres, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:14"),
        postgres.WithDatabase("testdb"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
    )
    if err != nil {
        t.Fatal(err)
    }
    defer postgres.Terminate(ctx)
    
    // Get connection string
    connStr, err := postgres.ConnectionString(ctx)
    if err != nil {
        t.Fatal(err)
    }
    
    // Initialize application with real database
    app := wire.InitializeTestApplicationWithDB(connStr)
    
    // Run your integration tests
    // ...
}
```
