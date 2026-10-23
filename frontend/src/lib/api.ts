import { ensureSession, redirectToLogin } from './session';

export async function api<T>(
	method: string,
	path: string,
	body?: unknown,
	options?: { skipSessionCheck?: boolean }
): Promise<T> {
	if (!options?.skipSessionCheck) {
		const me = await ensureSession();
		if (!me) throw new Error('unauthorized');
	}
	const headers: Record<string, string> = {};
	if (body !== undefined && !(body instanceof FormData)) {
		headers['Content-Type'] = 'application/json';
	}
	const res = await fetch(path, {
		method,
		headers: Object.keys(headers).length ? headers : undefined,
		body: body instanceof FormData ? body : body !== undefined ? JSON.stringify(body) : undefined
	});
	if (res.status === 401 || res.status === 403) {
		redirectToLogin();
		throw new Error('unauthorized');
	}
	if (!res.ok) {
		const text = await res.text();
		throw new Error(text || res.statusText);
	}
	const ct = res.headers.get('content-type');
	if (ct?.includes('application/json')) {
		return res.json() as Promise<T>;
	}
	return {} as T;
}
