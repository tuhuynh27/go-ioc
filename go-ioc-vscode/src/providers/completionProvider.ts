import * as vscode from 'vscode';

export class CompletionProvider implements vscode.CompletionItemProvider {
    
    provideCompletionItems(
        document: vscode.TextDocument,
        position: vscode.Position,
        token: vscode.CancellationToken,
        context: vscode.CompletionContext
    ): vscode.ProviderResult<vscode.CompletionItem[] | vscode.CompletionList> {
        
        const lineText = document.lineAt(position).text;
        const textBeforeCursor = lineText.substring(0, position.character);
        
        // Check if we're inside a struct tag
        if (this.isInStructTag(textBeforeCursor)) {
            return this.getStructTagCompletions(textBeforeCursor);
        }
        
        // Check if we're defining a struct with IoC patterns
        if (this.isInStructDefinition(document, position)) {
            return this.getStructFieldCompletions();
        }
        
        return [];
    }

    private isInStructTag(text: string): boolean {
        // Check if cursor is inside backticks (struct tags)
        const backtickCount = (text.match(/`/g) || []).length;
        return backtickCount % 2 === 1;
    }

    private isInStructDefinition(document: vscode.TextDocument, position: vscode.Position): boolean {
        // Look backwards to see if we're in a struct definition
        for (let i = position.line; i >= 0; i--) {
            const line = document.lineAt(i).text.trim();
            if (line.includes('type ') && line.includes('struct {')) {
                return true;
            }
            if (line.includes('}') && i < position.line) {
                return false;
            }
        }
        return false;
    }

    private getStructTagCompletions(textBeforeCursor: string): vscode.CompletionItem[] {
        const completions: vscode.CompletionItem[] = [];

        // Autowired tag completion
        if (!textBeforeCursor.includes('autowired')) {
            const autowiredTrueItem = new vscode.CompletionItem('autowired:"true"', vscode.CompletionItemKind.Property);
            autowiredTrueItem.detail = 'Go IoC: Mark field as autowired dependency';
            autowiredTrueItem.documentation = new vscode.MarkdownString('Marks this field as an autowired dependency that will be injected by the IoC container.');
            autowiredTrueItem.insertText = 'autowired:"true"';
            completions.push(autowiredTrueItem);

            const autowiredFalseItem = new vscode.CompletionItem('autowired:"false"', vscode.CompletionItemKind.Property);
            autowiredFalseItem.detail = 'Go IoC: Disable autowiring for this field';
            autowiredFalseItem.documentation = new vscode.MarkdownString('Explicitly disables autowiring for this field.');
            autowiredFalseItem.insertText = 'autowired:"false"';
            completions.push(autowiredFalseItem);
        }

        // Qualifier tag completion
        if (!textBeforeCursor.includes('qualifier')) {
            const qualifierItem = new vscode.CompletionItem('qualifier', vscode.CompletionItemKind.Property);
            qualifierItem.detail = 'Go IoC: Specify dependency qualifier';
            qualifierItem.documentation = new vscode.MarkdownString('Specifies a qualifier to disambiguate between multiple implementations of the same interface.');
            qualifierItem.insertText = new vscode.SnippetString('qualifier:"${1:qualifier-name}"');
            completions.push(qualifierItem);
        }

        // Implements tag completion
        if (!textBeforeCursor.includes('implements')) {
            const implementsItem = new vscode.CompletionItem('implements', vscode.CompletionItemKind.Property);
            implementsItem.detail = 'Go IoC: Specify implemented interface';
            implementsItem.documentation = new vscode.MarkdownString('Specifies which interface this component implements.');
            implementsItem.insertText = new vscode.SnippetString('implements:"${1:InterfaceName}"');
            completions.push(implementsItem);
        }

        // Value tag completion (for Qualifier struct)
        if (!textBeforeCursor.includes('value')) {
            const valueItem = new vscode.CompletionItem('value', vscode.CompletionItemKind.Property);
            valueItem.detail = 'Go IoC: Specify qualifier value';
            valueItem.documentation = new vscode.MarkdownString('Specifies the value for a qualifier.');
            valueItem.insertText = new vscode.SnippetString('value:"${1:qualifier-value}"');
            completions.push(valueItem);
        }

        return completions;
    }

    private getStructFieldCompletions(): vscode.CompletionItem[] {
        const completions: vscode.CompletionItem[] = [];

        // Component marker field
        const componentItem = new vscode.CompletionItem('Component struct{}', vscode.CompletionItemKind.Field);
        componentItem.detail = 'Go IoC: Component marker';
        componentItem.documentation = new vscode.MarkdownString('Marks this struct as an IoC component that can be managed by the container.');
        componentItem.insertText = 'Component struct{} // IoC component marker';
        completions.push(componentItem);

        // Qualifier field
        const qualifierItem = new vscode.CompletionItem('Qualifier struct{}', vscode.CompletionItemKind.Field);
        qualifierItem.detail = 'Go IoC: Qualifier field';
        qualifierItem.documentation = new vscode.MarkdownString('Adds a qualifier to this component for disambiguation.');
        qualifierItem.insertText = new vscode.SnippetString('Qualifier struct{} `value:"${1:qualifier-name}"`');
        completions.push(qualifierItem);

        // Implements field
        const implementsItem = new vscode.CompletionItem('Implements struct{}', vscode.CompletionItemKind.Field);
        implementsItem.detail = 'Go IoC: Interface implementation marker';
        implementsItem.documentation = new vscode.MarkdownString('Marks which interface this component implements.');
        implementsItem.insertText = new vscode.SnippetString('Implements struct{} `implements:"${1:InterfaceName}"`');
        completions.push(implementsItem);

        return completions;
    }

    resolveCompletionItem(
        item: vscode.CompletionItem,
        token: vscode.CancellationToken
    ): vscode.ProviderResult<vscode.CompletionItem> {
        return item;
    }
}