import * as vscode from 'vscode';
import { spawn } from 'child_process';
import { DecorationProvider, ValidationError } from './decorationProvider';

export class DiagnosticsProvider {
    private diagnostics: vscode.DiagnosticCollection;
    private statusBarItem: vscode.StatusBarItem;
    private decorationProvider: DecorationProvider;

    constructor() {
        this.diagnostics = vscode.languages.createDiagnosticCollection('go-ioc');
        this.statusBarItem = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Left, 100);
        this.statusBarItem.command = 'go-ioc.validate';
        this.decorationProvider = new DecorationProvider();
        this.updateStatusBar('Ready');
        this.statusBarItem.show();

        // Register event listeners
        vscode.workspace.onDidSaveTextDocument(this.validateDocument, this);
        vscode.workspace.onDidOpenTextDocument(this.validateDocument, this);
        vscode.workspace.onDidChangeTextDocument((e) => {
            if (e.document.languageId === 'go') {
                this.debounceValidation(e.document);
            }
        });
    }

    private debounceTimer: NodeJS.Timeout | undefined;

    private debounceValidation(document: vscode.TextDocument) {
        if (this.debounceTimer) {
            clearTimeout(this.debounceTimer);
        }
        this.debounceTimer = setTimeout(() => {
            this.validateDocument(document);
        }, 1000); // Wait 1 second after last change
    }

    private async validateDocument(document: vscode.TextDocument) {
        if (document.languageId !== 'go') {
            return;
        }

        // Check if document contains IoC components
        const text = document.getText();
        if (!this.containsIocComponents(text)) {
            this.diagnostics.delete(document.uri);
            return;
        }

        this.updateStatusBar('Validating...');
        
        try {
            // Get syntax validation errors from decoration provider
            const syntaxErrors = this.decorationProvider.getValidationErrors(document);
            
            // Get semantic validation from iocgen
            const semanticIssues = await this.runIocgenValidation();
            
            // Combine both types of issues
            const allIssues = [...this.convertValidationErrors(syntaxErrors), ...semanticIssues];
            
            this.updateDiagnostics(document, allIssues);
            this.updateStatusBar(allIssues.length > 0 ? `${allIssues.length} issues` : 'Valid');
        } catch (error) {
            this.updateStatusBar('Validation failed');
            console.error('IoC validation error:', error);
        }
    }

    private containsIocComponents(text: string): boolean {
        return text.includes('Component struct{}') ||
               text.includes('autowired:') ||
               text.includes('implements:') ||
               text.includes('Qualifier struct{}');
    }

    private async runIocgenValidation(): Promise<IocIssue[]> {
        return new Promise((resolve, reject) => {
            const config = vscode.workspace.getConfiguration('go-ioc');
            const iocgenPath = config.get<string>('iocgenPath', 'iocgen');
            
            const child = spawn(iocgenPath, ['--dry-run', '--verbose'], {
                cwd: vscode.workspace.rootPath,
                stdio: 'pipe'
            });

            let stdout = '';
            let stderr = '';

            child.stdout.on('data', (data) => {
                stdout += data.toString();
            });

            child.stderr.on('data', (data) => {
                stderr += data.toString();
            });

            child.on('close', (code) => {
                if (code === 0) {
                    resolve(this.parseValidationOutput(stdout));
                } else {
                    resolve(this.parseErrorOutput(stderr));
                }
            });

            child.on('error', (error) => {
                reject(error);
            });
        });
    }

    private parseValidationOutput(output: string): IocIssue[] {
        const issues: IocIssue[] = [];
        const lines = output.split('\n');

        for (const line of lines) {
            if (line.includes('Warning:') || line.includes('Error:')) {
                const issue = this.parseIssueLine(line);
                if (issue) {
                    issues.push(issue);
                }
            }
        }

        return issues;
    }

    private parseErrorOutput(output: string): IocIssue[] {
        const issues: IocIssue[] = [];
        
        if (output.includes('Could not resolve dependency')) {
            const issue: IocIssue = {
                message: 'Could not resolve dependency',
                severity: vscode.DiagnosticSeverity.Error,
                line: 0,
                column: 0,
                source: 'go-ioc'
            };
            issues.push(issue);
        }

        return issues;
    }

    private parseIssueLine(line: string): IocIssue | null {
        // Example: "Warning: Component at config.go:15 has unused dependency"
        const match = line.match(/(Warning|Error):\s*(.+?)\s*at\s*([^:]+):(\d+)/);
        
        if (match) {
            return {
                message: match[2],
                severity: match[1] === 'Error' ? vscode.DiagnosticSeverity.Error : vscode.DiagnosticSeverity.Warning,
                line: parseInt(match[4]) - 1, // VS Code uses 0-based line numbers
                column: 0,
                source: 'go-ioc',
                file: match[3]
            };
        }

        return null;
    }

    private updateDiagnostics(document: vscode.TextDocument, issues: IocIssue[]) {
        const diagnostics: vscode.Diagnostic[] = [];

        for (const issue of issues) {
            // Only add diagnostics for the current document
            if (!issue.file || document.fileName.endsWith(issue.file)) {
                const range = new vscode.Range(
                    new vscode.Position(issue.line, issue.column),
                    new vscode.Position(issue.line, issue.column + 1)
                );

                const diagnostic = new vscode.Diagnostic(
                    range,
                    issue.message,
                    issue.severity
                );
                diagnostic.source = issue.source;
                diagnostics.push(diagnostic);
            }
        }

        this.diagnostics.set(document.uri, diagnostics);
    }

    private convertValidationErrors(errors: ValidationError[]): IocIssue[] {
        return errors.map(error => ({
            message: `${error.message}. Suggestion: ${error.suggestion}`,
            severity: error.severity,
            line: error.line,
            column: error.column,
            source: 'go-ioc-syntax'
        }));
    }

    private updateStatusBar(text: string) {
        this.statusBarItem.text = `$(beaker) Go IoC: ${text}`;
        this.statusBarItem.tooltip = 'Click to validate IoC components';
    }

    dispose() {
        this.diagnostics.dispose();
        this.statusBarItem.dispose();
        this.decorationProvider.dispose();
        if (this.debounceTimer) {
            clearTimeout(this.debounceTimer);
        }
    }
}

interface IocIssue {
    message: string;
    severity: vscode.DiagnosticSeverity;
    line: number;
    column: number;
    source: string;
    file?: string;
}