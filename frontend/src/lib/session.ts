import type { AvatarMe } from './types';

let cached: AvatarMe | null = null;
let redirecting = false;

export function getAvatar(): AvatarMe | null {
	return cached;
}

export function setAvatar(me: AvatarMe | null) {
	cached = me;
}

export function redirectToLogin() {
	if (redirecting) return;
	redirecting = true;
	window.location.href = '/login';
}

export async function ensureSession(): Promise<AvatarMe | null> {
	const res = await fetch('/api/v1/avatar/me');
	if (res.status === 401 || res.status === 403) {
		redirectToLogin();
		return null;
	}
	if (!res.ok) {
		throw new Error('session check failed');
	}
	const me = (await res.json()) as AvatarMe;
	cached = me;
	return me;
}
