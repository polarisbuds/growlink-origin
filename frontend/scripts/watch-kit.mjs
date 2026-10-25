import { spawn } from 'node:child_process';
import { watch } from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const root = path.join(path.dirname(fileURLToPath(import.meta.url)), '..');
let timer;
let running = false;
let queued = false;

function runBuild() {
	if (running) {
		queued = true;
		return;
	}
	running = true;
	const child = spawn('bun', ['run', 'build:kit'], {
		cwd: root,
		stdio: 'inherit',
		shell: true
	});
	child.on('exit', (code) => {
		running = false;
		if (code !== 0) {
			console.error(`[dev:kit] build:kit failed with code ${code}`);
		}
		if (queued) {
			queued = false;
			runBuild();
		}
	});
}

function schedule() {
	clearTimeout(timer);
	timer = setTimeout(runBuild, 400);
}

runBuild();

for (const rel of ['src', 'svelte.config.js', 'vite.config.ts']) {
	watch(path.join(root, rel), { recursive: rel === 'src' }, schedule);
}

console.log('[dev:kit] watching for changes...');
