export function filterItems<T>(items: T[], query: string, labelOf: (item: T) => string): T[] {
	const q = (query || '').trim().toLowerCase();
	if (!q) return items;
	return items.filter((item) => labelOf(item).toLowerCase().includes(q));
}
