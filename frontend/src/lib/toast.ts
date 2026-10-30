export function showToast(message: string, type: 'success' | 'error' | 'info' = 'success') {
	const el = document.createElement('div');
	el.className = `alert alert-${type === 'error' ? 'error' : type === 'info' ? 'info' : 'success'} shadow-lg fixed bottom-4 right-4 z-50 w-auto max-w-sm`;
	el.setAttribute('role', 'alert');
	el.textContent = message;
	document.body.appendChild(el);
	setTimeout(() => el.remove(), 3200);
}
