# Go IoC VS Code Extension

VS Code extension providing IDE support for the Go IoC dependency injection framework.

## Features

### 🔧 Code Snippets
- **Component Creation**: `ioc-component` - Create IoC component structure
- **Interface Implementation**: `ioc-interface` - Component with interface implementation
- **Qualified Components**: `ioc-qualified` - Component with qualifier
- **Dependency Injection**: `ioc-autowired` - Add autowired dependencies
- **Lifecycle Methods**: `ioc-postconstruct`, `ioc-predestroy` - Lifecycle callbacks
- **Full Component**: `ioc-full` - Complete component with all features

### 💡 IntelliSense
- **Auto-completion** for IoC struct tags and annotations
- **Smart suggestions** for `autowired`, `qualifier`, `implements` tags
- **Context-aware** completion for component fields

### 🔍 Hover Information
- **Detailed documentation** on hover for IoC components and annotations
- **Dependency analysis** showing autowiring status and qualifiers
- **Lifecycle method** explanations and usage examples

### ⚠️ Real-time Validation
- **Live validation** of IoC components using `iocgen --dry-run`
- **Syntax validation** for struct tags (autowired, qualifier, implements)
- **Error highlighting** for dependency resolution issues and syntax errors
- **Status bar** integration showing validation status
- **Problem panel** integration for quick issue navigation

### 🎨 Visual Decorators
- **Annotation decorators** with emojis for better visual recognition:
  - `⚙️` Component markers (`Component struct{}`)
  - `🔗` Autowired dependencies (`autowired:"true"`)
  - `🏷️` Qualifiers (`qualifier:"name"` or `value:"name"`)
  - `🔌` Interface implementations (`implements:"InterfaceName"`)
  - `❌` Syntax errors and invalid annotations

### 🔗 Interface Navigation
- **Go to Definition** for interface names in `implements` tags
- **Ctrl+Click** (or **Cmd+Click** on Mac) on interface names to navigate to their definitions
- **Workspace-wide search** for interface definitions across all Go files
- **Multi-file support** with intelligent package resolution

### 🎯 Commands
- **Go IoC: Generate Wire Files** - Run `iocgen` to generate dependency injection code
- **Go IoC: Validate Components** - Validate IoC configuration without generating files
- **Go IoC: Analyze Dependencies** - Run comprehensive dependency analysis
- **Go IoC: Show Dependency Graph** - Display visual dependency relationships
- **Go IoC: List Components** - Show all discovered IoC components

## Requirements

- **Go IoC Framework**: Install the `iocgen` CLI tool
  ```bash
  go install github.com/tuhuynh27/go-ioc/cmd/iocgen@latest
  ```
- **Go Language**: Go 1.19 or later
- **VS Code**: Version 1.60.0 or later

## Extension Settings

This extension contributes the following settings:

- `go-ioc.iocgenPath`: Path to the iocgen binary (default: "iocgen")
- `go-ioc.autoGenerate`: Automatically generate wire files on save (default: false)
- `go-ioc.verboseOutput`: Enable verbose output for iocgen commands (default: false)

## Usage

### Creating IoC Components

1. **Type `ioc-component`** in a Go file and press Tab to create a basic component:
   ```go
   type ServiceName struct {
       Component struct{} // IoC component marker
       // Add fields here
   }
   ```

2. **Type `ioc-interface`** for a component implementing an interface:
   ```go
   type ServiceName struct {
       Component struct{}
       Implements struct{} `implements:"InterfaceName"`
       Qualifier struct{} `value:"qualifier"`
       // Add fields here
   }
   ```

3. **Type `ioc-autowired`** to add dependencies:
   ```go
   FieldName InterfaceType `autowired:"true" qualifier:"value"`
   ```

### Interface Navigation

- **Ctrl+Click** (or **Cmd+Click** on Mac) on interface names in `implements` tags
- Example: Click on `"UserService"` in `implements:"UserService"` to jump to the interface definition
- Works across files in your workspace

### Using Commands

- **Command Palette** (`Cmd+Shift+P` / `Ctrl+Shift+P`): Type "Go IoC" to see available commands
- **Context Menu**: Right-click in Go files for IoC-specific actions
- **Status Bar**: Click the IoC status indicator to validate components

### Validation and Debugging

The extension automatically validates your IoC configuration:

- **Green Status**: All components valid ✅
- **Warning Status**: Issues detected ⚠️  
- **Error Status**: Validation failed ❌

Hover over highlighted issues for detailed error information and suggested fixes.

## Example Workflow

1. **Create a service interface**:
   ```go
   type UserService interface {
       GetUser(id string) (*User, error)
   }
   ```

2. **Implement the service** using `ioc-interface` snippet:
   ```go
   type UserServiceImpl struct {
       Component struct{}
       Implements struct{} `implements:"UserService"`
       
       DB DatabaseInterface `autowired:"true"`
   }
   ```

3. **Add lifecycle methods** using snippets:
   ```go
   func (s *UserServiceImpl) PostConstruct() error {
       // Initialization logic
       return nil
   }
   ```

4. **Generate wire files** using Command Palette → "Go IoC: Generate Wire Files"

## Contributing

Issues and feature requests welcome! Visit the [GitHub repository](https://github.com/tuhuynh27/go-ioc) to contribute.

## License

MIT License - see the [LICENSE](LICENSE) file for details.