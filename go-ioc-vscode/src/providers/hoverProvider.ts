import * as vscode from 'vscode';

export class HoverProvider implements vscode.HoverProvider {
    
    provideHover(
        document: vscode.TextDocument,
        position: vscode.Position,
        token: vscode.CancellationToken
    ): vscode.ProviderResult<vscode.Hover> {
        
        const range = document.getWordRangeAtPosition(position);
        if (!range) {
            return;
        }

        const word = document.getText(range);
        const line = document.lineAt(position).text;
        
        // Check for IoC-specific patterns and provide hover information
        const hoverInfo = this.getHoverInfo(word, line, document, position);
        
        if (hoverInfo) {
            return new vscode.Hover(hoverInfo.content, hoverInfo.range || range);
        }

        return;
    }

    private getHoverInfo(word: string, line: string, document: vscode.TextDocument, position: vscode.Position): HoverInfo | null {
        
        // Component struct field hover
        if (word === 'Component' && line.includes('struct{}')) {
            return {
                content: this.createComponentMarkerHover(),
                range: document.getWordRangeAtPosition(position)
            };
        }

        // Qualifier struct field hover
        if (word === 'Qualifier' && line.includes('struct{}')) {
            return {
                content: this.createQualifierHover(line),
                range: document.getWordRangeAtPosition(position)
            };
        }

        // Implements struct field hover
        if (word === 'Implements' && line.includes('struct{}')) {
            return {
                content: this.createImplementsHover(line),
                range: document.getWordRangeAtPosition(position)
            };
        }

        // Autowired tag hover
        if (line.includes('autowired:') && this.isInStructTag(line, word)) {
            return {
                content: this.createAutowiredHover(line),
                range: this.getStructTagRange(document, position)
            };
        }

        // PostConstruct method hover
        if (word === 'PostConstruct' && line.includes('func')) {
            return {
                content: this.createPostConstructHover(),
                range: document.getWordRangeAtPosition(position)
            };
        }

        // PreDestroy method hover
        if (word === 'PreDestroy' && line.includes('func')) {
            return {
                content: this.createPreDestroyHover(),
                range: document.getWordRangeAtPosition(position)
            };
        }

        // Check if hovering over a dependency field
        const dependencyInfo = this.getDependencyFieldInfo(document, position);
        if (dependencyInfo) {
            return {
                content: this.createDependencyHover(dependencyInfo),
                range: document.getWordRangeAtPosition(position)
            };
        }

        return null;
    }

    private createComponentMarkerHover(): vscode.MarkdownString {
        const content = new vscode.MarkdownString();
        content.appendMarkdown('**Go IoC Component Marker**\n\n');
        content.appendMarkdown('This field marks the struct as an IoC component that will be managed by the dependency injection container.\n\n');
        content.appendMarkdown('**Features:**\n');
        content.appendMarkdown('- Automatic instantiation and dependency injection\n');
        content.appendMarkdown('- Lifecycle management (PostConstruct/PreDestroy)\n');
        content.appendMarkdown('- Interface-based dependency resolution\n\n');
        content.appendMarkdown('**Usage:** Add `Component struct{}` as the first field in your struct.');
        return content;
    }

    private createQualifierHover(line: string): vscode.MarkdownString {
        const content = new vscode.MarkdownString();
        content.appendMarkdown('**Go IoC Qualifier**\n\n');
        content.appendMarkdown('Qualifiers are used to disambiguate between multiple implementations of the same interface.\n\n');
        
        // Extract qualifier value if present
        const qualifierMatch = line.match(/value:"([^"]+)"/);
        if (qualifierMatch) {
            content.appendMarkdown(`**Current Qualifier:** \`${qualifierMatch[1]}\`\n\n`);
        }
        
        content.appendMarkdown('**Example:**\n');
        content.appendCodeblock('go', 'Qualifier struct{} `value:"primary"`');
        content.appendMarkdown('\n**Dependencies can reference this qualifier:**\n');
        content.appendCodeblock('go', 'MyService SomeInterface `autowired:"true" qualifier:"primary"`');
        return content;
    }

    private createImplementsHover(line: string): vscode.MarkdownString {
        const content = new vscode.MarkdownString();
        content.appendMarkdown('**Go IoC Interface Implementation**\n\n');
        content.appendMarkdown('This field specifies which interface this component implements, enabling interface-based dependency injection.\n\n');
        
        // Extract interface name if present
        const implementsMatch = line.match(/implements:"([^"]+)"/);
        if (implementsMatch) {
            content.appendMarkdown(`**Implements Interface:** \`${implementsMatch[1]}\`\n\n`);
        }
        
        content.appendMarkdown('**Example:**\n');
        content.appendCodeblock('go', 'Implements struct{} `implements:"UserService"`');
        content.appendMarkdown('\n**Other components can depend on this interface:**\n');
        content.appendCodeblock('go', 'UserSvc UserService `autowired:"true"`');
        return content;
    }

    private createAutowiredHover(line: string): vscode.MarkdownString {
        const content = new vscode.MarkdownString();
        content.appendMarkdown('**Go IoC Autowired Dependency**\n\n');
        content.appendMarkdown('This field will be automatically injected by the IoC container.\n\n');
        
        // Extract qualifier if present
        const qualifierMatch = line.match(/qualifier:"([^"]+)"/);
        if (qualifierMatch) {
            content.appendMarkdown(`**Qualifier:** \`${qualifierMatch[1]}\`\n\n`);
            content.appendMarkdown('The container will inject an implementation with this specific qualifier.\n\n');
        } else {
            content.appendMarkdown('The container will inject any available implementation of this interface.\n\n');
        }
        
        content.appendMarkdown('**Requirements:**\n');
        content.appendMarkdown('- The dependency type must be an interface\n');
        content.appendMarkdown('- At least one component must implement this interface\n');
        content.appendMarkdown('- If multiple implementations exist, use qualifiers to specify which one\n');
        return content;
    }

    private createPostConstructHover(): vscode.MarkdownString {
        const content = new vscode.MarkdownString();
        content.appendMarkdown('**Go IoC PostConstruct Lifecycle Method**\n\n');
        content.appendMarkdown('This method is called automatically after the component is instantiated and all dependencies are injected.\n\n');
        content.appendMarkdown('**Use Cases:**\n');
        content.appendMarkdown('- Initialize connections (database, cache, etc.)\n');
        content.appendMarkdown('- Validate configuration\n');
        content.appendMarkdown('- Set up internal state\n');
        content.appendMarkdown('- Start background processes\n\n');
        content.appendMarkdown('**Signature:**\n');
        content.appendCodeblock('go', 'func (component *ComponentName) PostConstruct() error');
        content.appendMarkdown('\n**Note:** If this method returns an error, component initialization will fail.');
        return content;
    }

    private createPreDestroyHover(): vscode.MarkdownString {
        const content = new vscode.MarkdownString();
        content.appendMarkdown('**Go IoC PreDestroy Lifecycle Method**\n\n');
        content.appendMarkdown('This method is called when the container is being shut down, allowing for cleanup.\n\n');
        content.appendMarkdown('**Use Cases:**\n');
        content.appendMarkdown('- Close database connections\n');
        content.appendMarkdown('- Stop background processes\n');
        content.appendMarkdown('- Release resources\n');
        content.appendMarkdown('- Save state\n\n');
        content.appendMarkdown('**Signature:**\n');
        content.appendCodeblock('go', 'func (component *ComponentName) PreDestroy() error');
        content.appendMarkdown('\n**Note:** This method is called during container shutdown.');
        return content;
    }

    private createDependencyHover(info: DependencyInfo): vscode.MarkdownString {
        const content = new vscode.MarkdownString();
        content.appendMarkdown(`**Go IoC Dependency: \`${info.fieldName}\`**\n\n`);
        content.appendMarkdown(`**Type:** \`${info.fieldType}\`\n`);
        
        if (info.isAutowired) {
            content.appendMarkdown('**Status:** Autowired ✅\n');
            if (info.qualifier) {
                content.appendMarkdown(`**Qualifier:** \`${info.qualifier}\`\n`);
            }
        } else {
            content.appendMarkdown('**Status:** Not autowired ⚠️\n');
            content.appendMarkdown('*Add `autowired:"true"` tag to enable dependency injection*\n');
        }
        
        content.appendMarkdown('\n**Dependency Resolution:**\n');
        content.appendMarkdown('The container will look for components that implement this interface');
        if (info.qualifier) {
            content.appendMarkdown(` with qualifier "${info.qualifier}"`);
        }
        content.appendMarkdown('.\n');
        
        return content;
    }

    private isInStructTag(line: string, word: string): boolean {
        const wordIndex = line.indexOf(word);
        if (wordIndex === -1) return false;
        
        const beforeWord = line.substring(0, wordIndex);
        const backtickCount = (beforeWord.match(/`/g) || []).length;
        return backtickCount % 2 === 1;
    }

    private getStructTagRange(document: vscode.TextDocument, position: vscode.Position): vscode.Range {
        const line = document.lineAt(position).text;
        const startIndex = line.lastIndexOf('`', position.character);
        const endIndex = line.indexOf('`', position.character);
        
        if (startIndex !== -1 && endIndex !== -1) {
            return new vscode.Range(
                new vscode.Position(position.line, startIndex),
                new vscode.Position(position.line, endIndex + 1)
            );
        }
        
        return document.getWordRangeAtPosition(position) || new vscode.Range(position, position);
    }

    private getDependencyFieldInfo(document: vscode.TextDocument, position: vscode.Position): DependencyInfo | null {
        const line = document.lineAt(position).text;
        
        // Check if this line looks like a field declaration with potential autowiring
        const fieldMatch = line.match(/(\w+)\s+(\w+).*`([^`]*)`/);
        if (!fieldMatch) return null;
        
        const fieldName = fieldMatch[1];
        const fieldType = fieldMatch[2];
        const tags = fieldMatch[3];
        
        const isAutowired = tags.includes('autowired:"true"');
        const qualifierMatch = tags.match(/qualifier:"([^"]+)"/);
        const qualifier = qualifierMatch ? qualifierMatch[1] : undefined;
        
        return {
            fieldName,
            fieldType,
            isAutowired,
            qualifier
        };
    }
}

interface HoverInfo {
    content: vscode.MarkdownString;
    range?: vscode.Range;
}

interface DependencyInfo {
    fieldName: string;
    fieldType: string;
    isAutowired: boolean;
    qualifier?: string;
}