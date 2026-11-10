export function statusLabel(status: number | string): string {
	const map: Record<string | number, string> = {
		1: '未着手',
		2: '下書き',
		3: '完了',
		4: '提出済み',
		pending: '未着手',
		draft: '下書き',
		done: '完了',
		submitted: '提出済み'
	};
	if (typeof status === 'string') return map[status] || status;
	return map[status] || String(status);
}

export function statusBadge(status: number | string): string {
	const s = status;
	if (s === 3 || s === 'done') return 'badge-success';
	if (s === 4 || s === 'submitted') return 'badge-info';
	if (s === 2 || s === 'draft') return 'badge-warning';
	return 'badge-ghost';
}
