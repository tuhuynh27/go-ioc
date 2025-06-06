---
sidebar_position: 5
---

# VS Code Extension

Get enhanced IDE support with the official Go IoC VS Code extension for a better development experience.

## Installation

### From VS Code Marketplace

**[📦 Install from VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=keva-dev.go-ioc)**

1. Open VS Code
2. Go to Extensions (Ctrl+Shift+X)
3. Search for "Go IoC"
4. Click Install

### Command Line Installation

```bash
code --install-extension keva-dev.go-ioc
```

## Features

### 🎨 Visual Decorators

IoC annotations are highlighted with emojis for better visual recognition:

- `⚙️` Component markers (`Component struct{}`)
- `🔗` Autowired dependencies (`autowired:"true"`)
- `🏷️` Qualifiers (`qualifier:"name"` or `value:"name"`)
- `🔌` Interface implementations (`implements:"InterfaceName"`)
- `❌` Syntax errors and invalid annotations

### 💡 IntelliSense & Auto-completion

Smart code completion for IoC struct tags and annotations:

- **Auto-completion** for `autowired`, `qualifier`, `implements` tags
- **Context-aware** suggestions for component fields
- **Syntax validation** while typing

### 🔍 Hover Information

Get detailed documentation on hover:

- **Component information** showing dependencies and qualifiers
- **Dependency analysis** with autowiring status
- **Lifecycle method** explanations and usage examples

### ⚠️ Real-time Validation

Live validation of IoC components:

- **Syntax validation** for struct tags
- **Dependency resolution** checking
- **Error highlighting** with actionable messages
- **Status bar** integration showing validation status
- **Problem panel** integration for quick issue navigation

### 🔗 Interface Navigation

Navigate to interface definitions with ease:

- **Go to Definition** for interface names in `implements` tags
- **Ctrl+Click** (or **Cmd+Click** on Mac) navigation
- **Workspace-wide search** for interface definitions
- **Multi-file support** with intelligent package resolution

### 🎯 Commands

Access Go IoC commands from the Command Palette (`Cmd+Shift+P` / `Ctrl+Shift+P`):

- **Go IoC: Generate Wire Files** - Run `iocgen` to generate dependency injection code
- **Go IoC: Validate Components** - Validate IoC configuration without generating files
- **Go IoC: Analyze Dependencies** - Run comprehensive dependency analysis
- **Go IoC: Show Dependency Graph** - Display visual dependency relationships
- **Go IoC: List Components** - Show all discovered IoC components

### 🔧 Code Snippets

Quick snippets for common IoC patterns:

#### Component Creation
- **`ioc-component`** - Create basic IoC component structure
- **`ioc-interface`** - Component with interface implementation
- **`ioc-qualified`** - Component with qualifier
- **`ioc-full`** - Complete component with all features

#### Dependency Injection
- **`ioc-autowired`** - Add autowired dependencies

#### Lifecycle Methods
- **`ioc-postconstruct`** - PostConstruct lifecycle method
- **`ioc-predestroy`** - PreDestroy lifecycle method

## Configuration

### Extension Settings

Configure the extension in VS Code settings:

```json
{
  "go-ioc.iocgenPath": "iocgen",
  "go-ioc.autoGenerate": false,
  "go-ioc.verboseOutput": false
}
```

#### Available Settings

- `go-ioc.iocgenPath`: Path to the iocgen binary (default: "iocgen")
- `go-ioc.autoGenerate`: Automatically generate wire files on save (default: false)
- `go-ioc.verboseOutput`: Enable verbose output for iocgen commands (default: false)

## Usage Examples

### Creating IoC Components

1. **Type `ioc-component`** and press Tab to create a basic component:

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

### 1. Create a Service Interface

```go
type UserService interface {
    GetUser(id string) (*User, error)
}
```

### 2. Implement the Service

Use the `ioc-interface` snippet:

```go
type UserServiceImpl struct {
    Component struct{}
    Implements struct{} `implements:"UserService"`
    
    DB DatabaseInterface `autowired:"true"`
}
```

### 3. Add Lifecycle Methods

Use the `ioc-postconstruct` snippet:

```go
func (s *UserServiceImpl) PostConstruct() error {
    // Initialization logic
    return nil
}
```

### 4. Generate Wire Files

Use Command Palette → "Go IoC: Generate Wire Files"

## Troubleshooting

### Extension Not Working

1. **Check Go IoC Installation**:
   ```bash
   go install github.com/tuhuynh27/go-ioc/cmd/iocgen@latest
   ```

2. **Verify PATH**: Ensure `$GOPATH/bin` is in your PATH

3. **Check Extension Settings**: Verify `go-ioc.iocgenPath` points to the correct binary

### Validation Issues

1. **Check Syntax**: Ensure struct tags are properly formatted
2. **Run Manual Validation**: Use `iocgen --dry-run` in terminal
3. **Check Dependencies**: Ensure all interface implementations exist

### Performance Issues

1. **Large Projects**: The extension scans all Go files; consider excluding large vendor directories
2. **Auto-generation**: Disable `go-ioc.autoGenerate` if it impacts performance

## Requirements

- **Go IoC Framework**: Install the `iocgen` CLI tool
- **Go Language**: Go 1.19 or later
- **VS Code**: Version 1.60.0 or later

## Contributing

Issues and feature requests welcome! Visit the [GitHub repository](https://github.com/tuhuynh27/go-ioc) to contribute.

The VS Code extension source code is available in the [`go-ioc-vscode`](https://github.com/tuhuynh27/go-ioc/tree/main/go-ioc-vscode) directory.