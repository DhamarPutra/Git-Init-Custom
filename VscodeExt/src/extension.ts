import * as vscode from 'vscode';
import { exec } from 'child_process';
import * as path from 'path';

export function activate(context: vscode.ExtensionContext) {
    let disposable = vscode.commands.registerCommand('git-green-screen.bootstrap', async () => {
        // 1. Minta input nama folder/directory tujuan
        const destDir = await vscode.window.showInputBox({
            prompt: 'Masukkan nama folder project baru',
            placeHolder: 'my-awesome-app'
        });
        if (!destDir) return;

        // 2. Minta input template (bisa juga memakai showQuickPick dari daftar template)
        const template = await vscode.window.showInputBox({
            prompt: 'Masukkan nama/URL template (kosongkan untuk default)',
            placeHolder: 'react, laravel, dll.'
        });

        // Tentukan path ke binary git-new Anda (atau pastikan sudah ada di sistem PATH)
        const binaryPath = 'git-new'; // Atau path absolut ke git-new.exe

        // Bentuk command argumen
        let command = `${binaryPath} "${destDir}"`;
        if (template) {
            command += ` --template "${template}"`;
        }

        // 3. Jalankan command dengan Progress UI di VSCode
        vscode.window.withProgress({
            location: vscode.ProgressLocation.Notification,
            title: "Bootstrapping project...",
            cancellable: false
        }, async (progress) => {
            return new Promise<void>((resolve, reject) => {
                // Jalankan CLI
                exec(command, { cwd: vscode.workspace.workspaceFolders?.[0]?.uri.fsPath }, (error, stdout, stderr) => {
                    if (error) {
                        vscode.window.showErrorMessage(`Gagal bootstrap: ${stderr || error.message}`);
                        reject(error);
                    } else {
                        vscode.window.showInformationMessage('Project berhasil dibuat!');
                        resolve();
                    }
                });
            });
        });
    });

    context.subscriptions.push(disposable);
}
