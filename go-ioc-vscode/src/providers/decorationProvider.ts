import * as vscode from 'vscode';

export class DecorationProvider {
    private autowiredDecoration: vscode.TextEditorDecorationType;
    private qualifierDecoration: vscode.TextEditorDecorationType;
    private implementsDecoration: vscode.TextEditorDecorationType;
    private componentDecoration: vscode.TextEditorDecorationType;
    private errorDecoration: vscode.TextEditorDecorationType;

    constructor() {
        // Create decoration types for different IoC annotations
        this.autowiredDecoration = vscode.window.createTextEditorDecorationType({
            backgroundColor: 'rgba(0, 123, 255, 0.1)',
            borderRadius: '3px',
            after: {
                contentText: ' 🔗',
                color: '#007bff',
                fontWeight: 'bold'
            }
        });

        this.qualifierDecoration = vscode.window.createTextEditorDecorationType({
            backgroundColor: 'rgba(40, 167, 69, 0.1)',
            borderRadius: '3px',
            after: {
                contentText: ' 🏷️',
                color: '#28a745',
                fontWeight: 'bold'
            }
        });

        this.implementsDecoration = vscode.window.createTextEditorDecorationType({
            backgroundColor: 'rgba(255, 193, 7, 0.1)',
            borderRadius: '3px',
            after: {
                contentText: ' 🔌',
                color: '#ffc107',
                fontWeight: 'bold'
            }
        });

        this.componentDecoration = vscode.window.createTextEditorDecorationType({
            backgroundColor: 'rgba(108, 117, 125, 0.1)',
            borderRadius: '3px',
            after: {
                contentText: ' ⚙️',
                color: '#6c757d',
                fontWeight: 'bold'
            }
        });

        this.errorDecoration = vscode.window.createTextEditorDecorationType({
            backgroundColor: 'rgba(220, 53, 69, 0.1)',
            borderRadius: '3px',
            after: {
                contentText: ' ❌',
                color: '#dc3545',
                fontWeight: 'bold'
            }
        });

        // Register for active editor changes
        vscode.window.onDidChangeActiveTextEditor(this.updateDecorations, this);
        vscode.workspace.onDidChangeTextDocument(this.onDocumentChange, this);
        
        // Initial decoration update
        this.updateDecorations(vscode.window.activeTextEditor);
    }

    private onDocumentChange(event: vscode.TextDocumentChangeEvent) {
        const activeEditor = vscode.window.activeTextEditor;
        if (activeEditor && event.document === activeEditor.document) {
            // Debounce to avoid excessive updates
            setTimeout(() => this.updateDecorations(activeEditor), 100);
        }
    }

    private updateDecorations(editor?: vscode.TextEditor) {
        if (!editor || editor.document.languageId !== 'go') {
            return;
        }

        const text = editor.document.getText();
        const lines = text.split('\n');

        const autowiredRanges: vscode.Range[] = [];
        const qualifierRanges: vscode.Range[] = [];
        const implementsRanges: vscode.Range[] = [];
        const componentRanges: vscode.Range[] = [];
        const errorRanges: vscode.Range[] = [];

        for (let lineIndex = 0; lineIndex < lines.length; lineIndex++) {
            const line = lines[lineIndex];
            
            // Check for Component struct marker
            const componentMatch = line.match(/Component\s+struct\{\}/);
            if (componentMatch) {
                const startPos = line.indexOf(componentMatch[0]);
                const range = new vscode.Range(
                    lineIndex, startPos,
                    lineIndex, startPos + componentMatch[0].length
                );
                componentRanges.push(range);
            }

            // Check for autowired tags
            const autowiredMatches = [...line.matchAll(/autowired\s*:\s*"([^"]*)"?/g)];
            for (const match of autowiredMatches) {
                const startPos = match.index!;
                const range = new vscode.Range(
                    lineIndex, startPos,
                    lineIndex, startPos + match[0].length
                );
                
                // Validate autowired syntax
                if (this.validateAutowiredTag(match[0])) {
                    autowiredRanges.push(range);
                } else {
                    errorRanges.push(range);
                }
            }

            // Check for qualifier tags
            const qualifierMatches = [...line.matchAll(/(qualifier|value)\s*:\s*"([^"]*)"?/g)];
            for (const match of qualifierMatches) {
                const startPos = match.index!;
                const range = new vscode.Range(
                    lineIndex, startPos,
                    lineIndex, startPos + match[0].length
                );
                
                // Validate qualifier syntax
                if (this.validateQualifierTag(match[0])) {
                    qualifierRanges.push(range);
                } else {
                    errorRanges.push(range);
                }
            }

            // Check for implements tags
            const implementsMatches = [...line.matchAll(/implements\s*:\s*"([^"]*)"?/g)];
            for (const match of implementsMatches) {
                const startPos = match.index!;
                const range = new vscode.Range(
                    lineIndex, startPos,
                    lineIndex, startPos + match[0].length
                );
                
                // Validate implements syntax
                if (this.validateImplementsTag(match[0])) {
                    implementsRanges.push(range);
                } else {
                    errorRanges.push(range);
                }
            }

            // Check for Qualifier and Implements struct markers
            const structMarkerMatches = [...line.matchAll(/(Qualifier|Implements)\s+struct\{\}/g)];
            for (const match of structMarkerMatches) {
                const startPos = match.index!;
                const range = new vscode.Range(
                    lineIndex, startPos,
                    lineIndex, startPos + match[0].length
                );
                
                if (match[1] === 'Qualifier') {
                    qualifierRanges.push(range);
                } else {
                    implementsRanges.push(range);
                }
            }
        }

        // Apply decorations
        editor.setDecorations(this.autowiredDecoration, autowiredRanges);
        editor.setDecorations(this.qualifierDecoration, qualifierRanges);
        editor.setDecorations(this.implementsDecoration, implementsRanges);
        editor.setDecorations(this.componentDecoration, componentRanges);
        editor.setDecorations(this.errorDecoration, errorRanges);
    }

    private validateAutowiredTag(tag: string): boolean {
        // Valid formats: autowired:"true", autowired:"false"
        const match = tag.match(/autowired\s*:\s*"(true|false)"/);
        return match !== null;
    }

    private validateQualifierTag(tag: string): boolean {
        // Valid formats: qualifier:"value", value:"value"
        const match = tag.match(/(qualifier|value)\s*:\s*"([^"]+)"/);
        return match !== null && match[2].trim().length > 0;
    }

    private validateImplementsTag(tag: string): boolean {
        // Valid format: implements:"InterfaceName"
        const match = tag.match(/implements\s*:\s*"([A-Z][a-zA-Z0-9]*(?:Interface)?)"/);
        return match !== null && match[1].trim().length > 0;
    }

    public getValidationErrors(document: vscode.TextDocument): ValidationError[] {
        const errors: ValidationError[] = [];
        const text = document.getText();
        const lines = text.split('\n');

        for (let lineIndex = 0; lineIndex < lines.length; lineIndex++) {
            const line = lines[lineIndex];

            // Check autowired tags
            const autowiredMatches = [...line.matchAll(/autowired\s*:\s*"([^"]*)"?/g)];
            for (const match of autowiredMatches) {
                if (!this.validateAutowiredTag(match[0])) {
                    errors.push({
                        line: lineIndex,
                        column: match.index!,
                        length: match[0].length,
                        message: 'Invalid autowired tag. Must be autowired:"true" or autowired:"false"',
                        severity: vscode.DiagnosticSeverity.Error,
                        suggestion: 'Use autowired:"true" or autowired:"false"'
                    });
                }
            }

            // Check qualifier tags
            const qualifierMatches = [...line.matchAll(/(qualifier|value)\s*:\s*"([^"]*)"?/g)];
            for (const match of qualifierMatches) {
                if (!this.validateQualifierTag(match[0])) {
                    errors.push({
                        line: lineIndex,
                        column: match.index!,
                        length: match[0].length,
                        message: 'Invalid qualifier tag. Must have a non-empty value',
                        severity: vscode.DiagnosticSeverity.Error,
                        suggestion: 'Use qualifier:"your-qualifier-name" with a descriptive name'
                    });
                }
            }

            // Check implements tags
            const implementsMatches = [...line.matchAll(/implements\s*:\s*"([^"]*)"?/g)];
            for (const match of implementsMatches) {
                if (!this.validateImplementsTag(match[0])) {
                    errors.push({
                        line: lineIndex,
                        column: match.index!,
                        length: match[0].length,
                        message: 'Invalid implements tag. Must specify a valid interface name starting with capital letter',
                        severity: vscode.DiagnosticSeverity.Error,
                        suggestion: 'Use implements:"InterfaceName" where InterfaceName starts with a capital letter'
                    });
                }
            }

            // Check for missing Component marker in structs with IoC tags
            if (this.hasIocTags(line) && !this.hasComponentMarker(lines, lineIndex)) {
                const firstTag = line.search(/(autowired|qualifier|implements|value):/);
                if (firstTag !== -1) {
                    errors.push({
                        line: lineIndex,
                        column: firstTag,
                        length: 1,
                        message: 'IoC tags found but struct is missing Component marker',
                        severity: vscode.DiagnosticSeverity.Warning,
                        suggestion: 'Add "Component struct{}" field to mark this struct as an IoC component'
                    });
                }
            }
        }

        return errors;
    }

    private hasIocTags(line: string): boolean {
        return /(autowired|qualifier|implements|value)\s*:/.test(line);
    }

    private hasComponentMarker(lines: string[], currentLine: number): boolean {
        // Look backwards to find the struct definition and check for Component marker
        for (let i = currentLine; i >= 0; i--) {
            const line = lines[i];
            if (line.includes('Component struct{}')) {
                return true;
            }
            if (line.includes('type ') && line.includes('struct {')) {
                // Found struct start, check next few lines for Component marker
                for (let j = i + 1; j <= Math.min(i + 10, lines.length - 1); j++) {
                    if (lines[j].includes('Component struct{}')) {
                        return true;
                    }
                    if (lines[j].includes('}') && !lines[j].includes('struct{}')) {
                        break; // End of struct
                    }
                }
                return false;
            }
        }
        return false;
    }

    dispose() {
        this.autowiredDecoration.dispose();
        this.qualifierDecoration.dispose();
        this.implementsDecoration.dispose();
        this.componentDecoration.dispose();
        this.errorDecoration.dispose();
    }
}

export interface ValidationError {
    line: number;
    column: number;
    length: number;
    message: string;
    severity: vscode.DiagnosticSeverity;
    suggestion: string;
}