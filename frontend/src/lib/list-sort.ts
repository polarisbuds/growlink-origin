export type SortDir = 'asc' | 'desc';

export function toggleSortKey<T extends string>(
	currentKey: T,
	currentDir: SortDir,
	key: T
): { sortKey: T; sortDir: SortDir } {
	if (currentKey === key) {
		return { sortKey: key, sortDir: currentDir === 'asc' ? 'desc' : 'asc' };
	}
	return { sortKey: key, sortDir: 'asc' };
}

export function sortIndicator<T extends string>(sortKey: T, currentKey: T, sortDir: SortDir): string {
	if (sortKey !== currentKey) return '';
	return sortDir === 'asc' ? ' ↑' : ' ↓';
}

export function compareJa(a: string, b: string): number {
	return a.localeCompare(b, 'ja');
}
