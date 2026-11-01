import { watch } from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { createKitBuilder } from './kit-build.mjs';

const root = path.join(path.dirname(fileURLToPath(import.meta.url)), '..');
const { runBuild, schedule } = createKitBuilder(root);

runBuild();

for (const rel of ['src', 'svelte.config.js', 'vite.config.ts']) {
	watch(path.join(root, rel), { recursive: rel === 'src' }, schedule);
}

console.log('[dev:kit] watching for changes...');
