import { spawn } from 'node:child_process';

export function createKitBuilder(cwd, debounceMs = 400) {
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
			cwd,
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
		timer = setTimeout(runBuild, debounceMs);
	}

	return { runBuild, schedule };
}
