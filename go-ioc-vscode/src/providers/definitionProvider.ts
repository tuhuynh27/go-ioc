import * as vscode from 'vscode';
import * as fs from 'fs';
import * as path from 'path';

export class DefinitionProvider implements vscode.DefinitionProvider {

    async provideDefinition(
        document: vscode.TextDocument,
        position: vscode.Position,
        token: vscode.CancellationToken
    ): Promise<vscode.Definition | vscode.DefinitionLink[] | undefined> {

        const line = document.lineAt(position).text;
        
        // Check if we're in an implements tag
        const implementsMatch = this.getImplementsTagAtPosition(line, position.character);
        if (!implementsMatch) {
            return undefined;
        }

        const interfaceName = implementsMatch.interfaceName;
        const interfaceRange = new vscode.Range(
            position.line,
            implementsMatch.startIndex,
            position.line,
            implementsMatch.endIndex
        );

        // Check if the cursor is within the interface name
        if (position.character < implementsMatch.startIndex || position.character > implementsMatch.endIndex) {
            return undefined;
        }

        // Search for the interface definition
        const interfaceDefinition = await this.findInterfaceDefinition(interfaceName, document.uri);
        if (!interfaceDefinition) {
            return undefined;
        }

        return [{
            targetUri: interfaceDefinition.uri,
            targetRange: interfaceDefinition.range,
            targetSelectionRange: interfaceDefinition.selectionRange,
            originSelectionRange: interfaceRange
        }];
    }

    private getImplementsTagAtPosition(line: string, position: number): InterfaceMatch | null {
        // Find all implements tags in the line
        const implementsRegex = /implements\s*:\s*"([^"]+)"/g;
        let match;

        while ((match = implementsRegex.exec(line)) !== null) {
            const fullMatch = match[0];
            const interfaceName = match[1].trim();
            const tagStartIndex = match.index!;
            const tagEndIndex = tagStartIndex + fullMatch.length;

            // Skip empty interface names
            if (!interfaceName) {
                continue;
            }

            // Find the position of the interface name within the tag
            const interfaceStartInTag = fullMatch.indexOf(`"${interfaceName}"`);
            const interfaceNameStartIndex = tagStartIndex + interfaceStartInTag + 1; // +1 for opening quote
            const interfaceNameEndIndex = interfaceNameStartIndex + interfaceName.length;

            // Check if position is within the interface name
            if (position >= interfaceNameStartIndex && position <= interfaceNameEndIndex) {
                return {
                    interfaceName,
                    startIndex: interfaceNameStartIndex,
                    endIndex: interfaceNameEndIndex,
                    tagStartIndex,
                    tagEndIndex
                };
            }
        }

        return null;
    }

    private async findInterfaceDefinition(interfaceName: string, currentUri: vscode.Uri): Promise<InterfaceDefinition | null> {
        // First, try to find in the current file
        const currentFileDefinition = await this.searchInterfaceInFile(interfaceName, currentUri);
        if (currentFileDefinition) {
            return currentFileDefinition;
        }

        // Search in workspace files
        const workspaceDefinition = await this.searchInterfaceInWorkspace(interfaceName, currentUri);
        return workspaceDefinition;
    }

    private async searchInterfaceInFile(interfaceName: string, fileUri: vscode.Uri): Promise<InterfaceDefinition | null> {
        try {
            const document = await vscode.workspace.openTextDocument(fileUri);
            return this.findInterfaceInDocument(interfaceName, document);
        } catch (error) {
            console.error(`Error reading file ${fileUri.fsPath}:`, error);
            return null;
        }
    }

    private async searchInterfaceInWorkspace(interfaceName: string, excludeUri: vscode.Uri): Promise<InterfaceDefinition | null> {
        const workspaceFolders = vscode.workspace.workspaceFolders;
        if (!workspaceFolders) {
            return null;
        }

        // Search through all Go files in the workspace
        for (const folder of workspaceFolders) {
            const goFiles = await vscode.workspace.findFiles(
                new vscode.RelativePattern(folder, '**/*.go'),
                '**/node_modules/**', // Exclude node_modules
                1000 // Limit to 1000 files for performance
            );

            for (const fileUri of goFiles) {
                // Skip the current file as we already searched it
                if (fileUri.toString() === excludeUri.toString()) {
                    continue;
                }

                const definition = await this.searchInterfaceInFile(interfaceName, fileUri);
                if (definition) {
                    return definition;
                }
            }
        }

        return null;
    }

    private findInterfaceInDocument(interfaceName: string, document: vscode.TextDocument): InterfaceDefinition | null {
        const text = document.getText();
        const lines = text.split('\n');

        for (let lineIndex = 0; lineIndex < lines.length; lineIndex++) {
            const line = lines[lineIndex];
            
            // Look for interface declarations
            // Pattern: type InterfaceName interface {
            const interfaceRegex = new RegExp(`\\btype\\s+(${interfaceName})\\s+interface\\s*\\{`, 'g');
            const match = interfaceRegex.exec(line);

            if (match) {
                const interfaceNameInLine = match[1];
                const interfaceStartIndex = match.index! + match[0].indexOf(interfaceNameInLine);
                const interfaceEndIndex = interfaceStartIndex + interfaceNameInLine.length;

                // Find the end of the interface (closing brace)
                let endLine = lineIndex;
                let braceCount = 1; // We start after the opening brace
                let endChar = line.length;

                // Look for the closing brace
                for (let i = lineIndex; i < lines.length && braceCount > 0; i++) {
                    const currentLine = lines[i];
                    let startPos = i === lineIndex ? line.indexOf('{') + 1 : 0;

                    for (let j = startPos; j < currentLine.length; j++) {
                        if (currentLine[j] === '{') {
                            braceCount++;
                        } else if (currentLine[j] === '}') {
                            braceCount--;
                            if (braceCount === 0) {
                                endLine = i;
                                endChar = j + 1;
                                break;
                            }
                        }
                    }
                }

                const range = new vscode.Range(
                    new vscode.Position(lineIndex, 0),
                    new vscode.Position(endLine, endChar)
                );

                const selectionRange = new vscode.Range(
                    new vscode.Position(lineIndex, interfaceStartIndex),
                    new vscode.Position(lineIndex, interfaceEndIndex)
                );

                return {
                    uri: document.uri,
                    range,
                    selectionRange
                };
            }
        }

        return null;
    }

    /**
     * Enhanced interface search that also handles interfaces from different packages
     */
    private async searchInterfaceWithPackageSupport(interfaceName: string, currentDocument: vscode.TextDocument): Promise<InterfaceDefinition | null> {
        // First, check if the interface has a package prefix (e.g., "http.Handler")
        const packageMatch = interfaceName.match(/^(\w+)\.(\w+)$/);
        
        if (packageMatch) {
            const packageName = packageMatch[1];
            const actualInterfaceName = packageMatch[2];
            
            // Look for imports in the current document to find the package path
            const imports = this.extractImports(currentDocument);
            const packagePath = imports[packageName];
            
            if (packagePath) {
                // Search in the specific package
                return await this.searchInterfaceInPackage(actualInterfaceName, packagePath);
            }
        }

        // Fallback to regular search
        return await this.findInterfaceDefinition(interfaceName, currentDocument.uri);
    }

    private extractImports(document: vscode.TextDocument): { [alias: string]: string } {
        const imports: { [alias: string]: string } = {};
        const text = document.getText();
        
        // Match both single imports and import blocks
        const importRegex = /import\s+(?:\(\s*([\s\S]*?)\s*\)|(\S+\s+)?"([^"]+)")/g;
        let match;

        while ((match = importRegex.exec(text)) !== null) {
            if (match[1]) {
                // Import block
                const blockContent = match[1];
                const lineRegex = /(?:(\w+)\s+)?"([^"]+)"/g;
                let lineMatch;
                
                while ((lineMatch = lineRegex.exec(blockContent)) !== null) {
                    const alias = lineMatch[1] || path.basename(lineMatch[2]);
                    const packagePath = lineMatch[2];
                    imports[alias] = packagePath;
                }
            } else {
                // Single import
                const alias = match[2] ? match[2].trim() : path.basename(match[3]);
                const packagePath = match[3];
                imports[alias] = packagePath;
            }
        }

        return imports;
    }

    private async searchInterfaceInPackage(interfaceName: string, packagePath: string): Promise<InterfaceDefinition | null> {
        // This is a simplified implementation
        // In a full implementation, you'd need to resolve the actual file paths from package imports
        // For now, we'll search in the workspace for files that might contain this package
        
        const workspaceFolders = vscode.workspace.workspaceFolders;
        if (!workspaceFolders) {
            return null;
        }

        // Look for files in directories that match the package name
        const packageName = path.basename(packagePath);
        
        for (const folder of workspaceFolders) {
            const pattern = `**/${packageName}/**/*.go`;
            const goFiles = await vscode.workspace.findFiles(
                new vscode.RelativePattern(folder, pattern),
                '**/node_modules/**',
                100
            );

            for (const fileUri of goFiles) {
                const definition = await this.searchInterfaceInFile(interfaceName, fileUri);
                if (definition) {
                    return definition;
                }
            }
        }

        return null;
    }
}

interface InterfaceMatch {
    interfaceName: string;
    startIndex: number;
    endIndex: number;
    tagStartIndex: number;
    tagEndIndex: number;
}

interface InterfaceDefinition {
    uri: vscode.Uri;
    range: vscode.Range;
    selectionRange: vscode.Range;
}