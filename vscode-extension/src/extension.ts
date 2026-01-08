import * as fs from 'node:fs';
import * as os from 'node:os';
import * as path from 'node:path';
import {
    type ExtensionContext,
    type OutputChannel,
    StatusBarAlignment,
    type StatusBarItem,
    commands,
    window,
    workspace,
} from 'vscode';
import {
    LanguageClient,
    type LanguageClientOptions,
    type ServerOptions,
    TransportKind,
} from 'vscode-languageclient/node';

let client: LanguageClient | undefined;
let statusBarItem: StatusBarItem;
let outputChannel: OutputChannel;

export async function activate(context: ExtensionContext): Promise<void> {
    outputChannel = window.createOutputChannel('Godot Language Server');
    context.subscriptions.push(outputChannel);

    statusBarItem = window.createStatusBarItem(StatusBarAlignment.Right, 100);
    statusBarItem.text = '$(loading~spin) GDLS';
    statusBarItem.tooltip = 'Godot Language Server';
    statusBarItem.command = 'gdls.showOutputChannel';
    context.subscriptions.push(statusBarItem);

    // Register commands
    context.subscriptions.push(
        commands.registerCommand('gdls.restartServer', async () => {
            await restartServer(context);
        }),
    );

    context.subscriptions.push(
        commands.registerCommand('gdls.showOutputChannel', () => {
            outputChannel.show();
        }),
    );

    // Start the language server
    const config = workspace.getConfiguration('gdls');
    if (config.get<boolean>('server.enabled', true)) {
        await startServer(context);
    }

    // Watch for configuration changes
    context.subscriptions.push(
        workspace.onDidChangeConfiguration(async (e) => {
            if (e.affectsConfiguration('gdls.server.enabled')) {
                const enabled = workspace
                    .getConfiguration('gdls')
                    .get<boolean>('server.enabled', true);
                if (enabled && !client) {
                    await startServer(context);
                } else if (!enabled && client) {
                    await stopServer();
                }
            } else if (e.affectsConfiguration('gdls.server.path')) {
                await restartServer(context);
            }
        }),
    );
}

export async function deactivate(): Promise<void> {
    await stopServer();
}

async function startServer(context: ExtensionContext): Promise<void> {
    const serverPath = getServerPath(context);

    if (!serverPath) {
        const message =
            'Could not find gdls executable. Please install it or configure the path in settings.';
        outputChannel.appendLine(`Error: ${message}`);
        window.showErrorMessage(message);
        updateStatusBar('error', 'GDLS server not found');
        return;
    }

    outputChannel.appendLine(`Starting Godot language server: ${serverPath}`);
    statusBarItem.show();
    updateStatusBar('loading', 'Starting GDLS server...');

    const serverOptions: ServerOptions = {
        run: {
            command: serverPath,
            transport: TransportKind.stdio,
        },
        debug: {
            command: serverPath,
            transport: TransportKind.stdio,
        },
    };

    const clientOptions: LanguageClientOptions = {
        documentSelector: [
            { scheme: 'file', language: 'tscn' },
            { scheme: 'file', language: 'gdshader' },
        ],
        synchronize: {
            fileEvents: workspace.createFileSystemWatcher(
                '**/*.{tscn,escn,gdshader,gdshaderinc}',
            ),
        },
        outputChannel,
        traceOutputChannel: outputChannel,
    };

    client = new LanguageClient(
        'gdls',
        'Godot Language Server',
        serverOptions,
        clientOptions,
    );

    try {
        await client.start();
        outputChannel.appendLine('Godot language server started successfully');
        updateStatusBar('ready', 'GDLS server running');
    } catch (error) {
        const errorMessage =
            error instanceof Error ? error.message : String(error);
        outputChannel.appendLine(`Failed to start server: ${errorMessage}`);
        window.showErrorMessage(
            `Failed to start Godot language server: ${errorMessage}`,
        );
        updateStatusBar('error', 'GDLS server failed to start');
        client = undefined;
    }
}

async function stopServer(): Promise<void> {
    if (client) {
        outputChannel.appendLine('Stopping Godot language server...');
        await client.stop();
        client = undefined;
        statusBarItem.hide();
        outputChannel.appendLine('Godot language server stopped');
    }
}

async function restartServer(context: ExtensionContext): Promise<void> {
    outputChannel.appendLine('Restarting Godot language server...');
    await stopServer();
    await startServer(context);
}

function getServerPath(context: ExtensionContext): string | undefined {
    const config = workspace.getConfiguration('gdls');

    // 1. Check user-configured path
    const configuredPath = config.get<string>('server.path');
    if (configuredPath && fs.existsSync(configuredPath)) {
        return configuredPath;
    }

    // 2. Check bundled binary
    const bundledPath = getBundledServerPath(context);
    if (bundledPath && fs.existsSync(bundledPath)) {
        return bundledPath;
    }

    // 3. Check PATH
    const pathEnvPath = findInPath('gdls');
    if (pathEnvPath) {
        return pathEnvPath;
    }

    return undefined;
}

function getBundledServerPath(context: ExtensionContext): string | undefined {
    const platform = os.platform();
    const arch = os.arch();

    let binaryName: string;

    if (platform === 'win32') {
        binaryName = 'gdls-windows-amd64.exe';
    } else if (platform === 'darwin') {
        binaryName =
            arch === 'arm64' ? 'gdls-darwin-arm64' : 'gdls-darwin-amd64';
    } else if (platform === 'linux') {
        binaryName = arch === 'arm64' ? 'gdls-linux-arm64' : 'gdls-linux-amd64';
    } else {
        return undefined;
    }

    return path.join(context.extensionPath, 'bin', binaryName);
}

function findInPath(executable: string): string | undefined {
    const envPath = process.env.PATH || '';
    const pathSeparator = os.platform() === 'win32' ? ';' : ':';
    const paths = envPath.split(pathSeparator);

    const extensions =
        os.platform() === 'win32' ? ['', '.exe', '.cmd', '.bat'] : [''];

    for (const dir of paths) {
        for (const ext of extensions) {
            const fullPath = path.join(dir, executable + ext);
            if (fs.existsSync(fullPath)) {
                return fullPath;
            }
        }
    }

    return undefined;
}

function updateStatusBar(
    status: 'loading' | 'ready' | 'error',
    tooltip: string,
): void {
    switch (status) {
        case 'loading':
            statusBarItem.text = '$(loading~spin) GDLS';
            statusBarItem.backgroundColor = undefined;
            break;
        case 'ready':
            statusBarItem.text = '$(check) GDLS';
            statusBarItem.backgroundColor = undefined;
            break;
        case 'error':
            statusBarItem.text = '$(error) GDLS';
            statusBarItem.backgroundColor = undefined;
            break;
    }
    statusBarItem.tooltip = tooltip;
}
