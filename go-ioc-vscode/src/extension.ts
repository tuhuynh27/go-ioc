import * as vscode from 'vscode';
import { CompletionProvider } from './providers/completionProvider';
import { DiagnosticsProvider } from './providers/diagnosticsProvider';
import { HoverProvider } from './providers/hoverProvider';
import { DecorationProvider } from './providers/decorationProvider';
import { DefinitionProvider } from './providers/definitionProvider';

export function activate(context: vscode.ExtensionContext) {
    console.log('Go IoC extension is now active!');

    // Register completion provider for Go IoC struct tags
    const completionProvider = new CompletionProvider();
    context.subscriptions.push(
        vscode.languages.registerCompletionItemProvider(
            'go', 
            completionProvider, 
            '`', '"', ':'
        )
    );

    // Register diagnostics provider for validation
    const diagnosticsProvider = new DiagnosticsProvider();
    context.subscriptions.push(diagnosticsProvider);

    // Register hover provider for component information
    const hoverProvider = new HoverProvider();
    context.subscriptions.push(
        vscode.languages.registerHoverProvider('go', hoverProvider)
    );

    // Register definition provider for interface navigation
    const definitionProvider = new DefinitionProvider();
    context.subscriptions.push(
        vscode.languages.registerDefinitionProvider('go', definitionProvider)
    );

    // Register commands
    registerCommands(context);

    // Auto-generate on save if enabled
    context.subscriptions.push(
        vscode.workspace.onDidSaveTextDocument((document) => {
            if (document.languageId === 'go') {
                const config = vscode.workspace.getConfiguration('go-ioc');
                if (config.get<boolean>('autoGenerate', false)) {
                    executeIocgenCommand('');
                }
            }
        })
    );
}

function registerCommands(context: vscode.ExtensionContext) {
    // Generate wire files command
    context.subscriptions.push(
        vscode.commands.registerCommand('go-ioc.generate', () => {
            executeIocgenCommand('');
        })
    );

    // Validate components command
    context.subscriptions.push(
        vscode.commands.registerCommand('go-ioc.validate', () => {
            executeIocgenCommand('--dry-run');
        })
    );

    // Analyze dependencies command
    context.subscriptions.push(
        vscode.commands.registerCommand('go-ioc.analyze', () => {
            executeIocgenCommand('--analyze');
        })
    );

    // Show dependency graph command
    context.subscriptions.push(
        vscode.commands.registerCommand('go-ioc.graph', () => {
            executeIocgenCommand('--graph');
        })
    );

    // List components command
    context.subscriptions.push(
        vscode.commands.registerCommand('go-ioc.list', () => {
            executeIocgenCommand('--list');
        })
    );
}

function executeIocgenCommand(args: string) {
    const config = vscode.workspace.getConfiguration('go-ioc');
    const iocgenPath = config.get<string>('iocgenPath', 'iocgen');
    const verbose = config.get<boolean>('verboseOutput', false);
    
    let command = iocgenPath;
    if (args) {
        command += ` ${args}`;
    }
    if (verbose && !args.includes('--verbose')) {
        command += ' --verbose';
    }

    const terminal = vscode.window.createTerminal({
        name: 'Go IoC',
        cwd: vscode.workspace.rootPath
    });
    
    terminal.sendText(command);
    terminal.show();
}

export function deactivate() {
    console.log('Go IoC extension is now deactivated!');
}