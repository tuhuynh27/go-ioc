import type {ReactNode} from 'react';
import CodeBlock from '@theme/CodeBlock';
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';
import Terminal from '@site/src/components/Terminal';
import styles from './styles.module.css';

const serviceCode = `// Define your service interface
type UserService interface {
    CreateUser(name string) (*User, error)
    GetUser(id string) (*User, error)
}`;

const componentCode = `// Create IoC components with familiar syntax
type UserServiceImpl struct {
    Component struct{}                           // @Component
    Implements struct{} \`implements:"UserService"\` // @Service
    Qualifier struct{} \`value:"primary"\`           // @Qualifier
    
    // Dependency injection with @Autowired
    Database  DatabaseInterface \`autowired:"true"\`
    Logger    LoggerInterface   \`autowired:"true" qualifier:"console"\`
    Cache     CacheInterface    \`autowired:"true" qualifier:"redis"\`
}

func (s *UserServiceImpl) CreateUser(name string) (*User, error) {
    s.Logger.Info("Creating user: " + name)
    // Implementation here...
    return &User{Name: name}, nil
}`;

const generatedCode = `// Generated wire code (wire/wire_gen.go)
type Container struct {
    DatabaseService *database.DatabaseService
    LoggerService   *logger.ConsoleLogger
    CacheService    *cache.RedisCache
    UserService     *user.UserServiceImpl
}

func Initialize() (*Container, func()) {
    container := &Container{}
    
    // Dependencies resolved in correct order
    container.DatabaseService = database.NewDatabaseService()
    container.LoggerService = &logger.ConsoleLogger{}
    container.CacheService = cache.NewRedisCache()
    
    // Inject dependencies automatically
    container.UserService = &user.UserServiceImpl{
        Database: container.DatabaseService,
        Logger:   container.LoggerService,
        Cache:    container.CacheService,
    }
    
    return container, cleanup
}`;

const usageCode = `// Use in your application
func main() {
    // One line initialization
    container, cleanup := wire.Initialize()
    defer cleanup()
    
    // Get your services ready to use
    userService := container.UserService
    
    // All dependencies are injected and ready!
    user, err := userService.CreateUser("John Doe")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Created user: %s\\n", user.Name)
}`;

export default function CodeDemo(): ReactNode {
  return (
    <section className={styles.codeDemo}>
      <div className="container">
        <div className="row">
          <div className="col col--12">
            <div className="text--center margin-bottom--lg">
              <h2>See Go IoC in Action</h2>
              <p className="hero__subtitle">
                Familiar Spring-style syntax meets Go's compile-time safety
              </p>
            </div>
            
            <Tabs>
              <TabItem value="define" label="1. Define Services" default>
                <CodeBlock language="go" title="service/user.go">
                  {serviceCode}
                </CodeBlock>
              </TabItem>
              
              <TabItem value="component" label="2. Create Components">
                <CodeBlock language="go" title="service/user_impl.go">
                  {componentCode}
                </CodeBlock>
              </TabItem>
              
              <TabItem value="generate" label="3. Generate Wire Code">
                <div className={styles.generateStep}>
                  <Terminal title="go-ioc-project — -zsh — 80×24">
                    iocgen
                  </Terminal>
                  <div className={styles.arrow}>↓</div>
                  <CodeBlock language="go" title="wire/wire_gen.go (Generated)">
                    {generatedCode}
                  </CodeBlock>
                </div>
              </TabItem>
              
              <TabItem value="use" label="4. Use It">
                <CodeBlock language="go" title="main.go">
                  {usageCode}
                </CodeBlock>
              </TabItem>
            </Tabs>
            
            <div className={styles.highlights}>
              <div className="row">
                <div className="col col--4">
                  <div className={styles.highlight}>
                    <h4>🏷️ Familiar Annotations</h4>
                    <p>Spring-like <code>@Component</code>, <code>@Autowired</code>, and <code>@Qualifier</code> syntax</p>
                  </div>
                </div>
                <div className="col col--4">
                  <div className={styles.highlight}>
                    <h4>⚡ Zero Runtime Cost</h4>
                    <p>Pure compile-time code generation with no reflection or runtime overhead</p>
                  </div>
                </div>
                <div className="col col--4">
                  <div className={styles.highlight}>
                    <h4>🔒 Type Safe</h4>
                    <p>All dependencies validated at compile time with full Go type safety</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}