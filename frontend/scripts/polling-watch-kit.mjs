import { readdir, stat } from 'node:fs/promises';
import path from 'node:path';
import { fileURLToPath, pathToFileURL } from 'node:url';
import { createKitBuilder } from './kit-build.mjs';

const POLL_MS = 1000;
const FILES = ['svelte.config.js', 'vite.config.ts'];

export async function collectMtimes(root) {
	const mtimes = new Map();

	async function walk(dir) {
		let entries;
		try {
			entries = await readdir(dir, { withFileTypes: true });
		} catch {
			return;
		}
		for (const entry of entries) {
			const full = path.join(dir, entry.name);
			if (entry.isDirectory()) await walk(full);
			else if (entry.isFile()) {
				try {
					const info = await stat(full);
					mtimes.set(full, info.mtimeMs);
				} catch {
					// The file can disappear between readdir and stat.
				}
			}
		}
	}

	await walk(path.join(root, 'src'));
	for (const file of FILES) {
		const full = path.join(root, file);
		try {
			const info = await stat(full);
			mtimes.set(full, info.mtimeMs);
		} catch {
			// A missing config file is a change once it exists again.
		}
	}
	return mtimes;
}

export function mtimesChanged(prev, next) {
	if (prev.size !== next.size) return true;
	for (const [file, mtime] of next) {
		if (prev.get(file) !== mtime) return true;
	}
	return false;
}

async function main() {
	const root = path.join(path.dirname(fileURLToPath(import.meta.url)), '..');
	const { runBuild, schedule } = createKitBuilder(root);
	let current = await collectMtimes(root);
	let scanning = false;

	runBuild();
	setInterval(async () => {
		if (scanning) return;
		scanning = true;
		try {
			const next = await collectMtimes(root);
			if (mtimesChanged(current, next)) {
				current = next;
				schedule();
			}
		} finally {
			scanning = false;
		}
	}, POLL_MS);

	console.log(`[dev:kit] polling for changes every ${POLL_MS}ms...`);
}

if (process.argv[1] && import.meta.url === pathToFileURL(path.resolve(process.argv[1])).href) {
	main();
}
