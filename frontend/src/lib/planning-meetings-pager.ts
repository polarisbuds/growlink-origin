export type MeetingsPageResponse = {
	items: unknown[];
	total: number;
	page: number;
	pageSize: number;
};

export type MeetingListQuery = {
	page: number;
	dir: 'asc' | 'desc';
	mentorUuids: string[];
	menteeUuids: string[];
	departmentUuids: string[];
	status: string;
	from: string;
	to: string;
};

export function buildMeetingsAllQuery(q: MeetingListQuery): string {
	const params = new URLSearchParams();
	params.set('page', String(q.page));
	params.set('dir', q.dir);
	params.set('pageSize', '20');
	if (q.mentorUuids.length) params.set('mentorUuids', q.mentorUuids.join(','));
	if (q.menteeUuids.length) params.set('menteeUuids', q.menteeUuids.join(','));
	if (q.departmentUuids.length) params.set('departmentUuids', q.departmentUuids.join(','));
	if (q.status) params.set('status', q.status);
	if (q.from) params.set('from', q.from);
	if (q.to) params.set('to', q.to);
	return params.toString();
}

export function lastPageForTotal(total: number, pageSize: number): number {
	if (total <= 0) return 1;
	return Math.ceil(total / pageSize);
}

export type PageToken = number | 'ellipsis';

export function visiblePageNumbers(current: number, last: number): PageToken[] {
	if (last <= 1) return [1];
	const pages = new Set<number>();
	pages.add(1);
	pages.add(last);
	for (let p = current - 1; p <= current + 1; p++) {
		if (p >= 1 && p <= last) pages.add(p);
	}
	const sorted = [...pages].sort((a, b) => a - b);
	const out: PageToken[] = [];
	let prev = 0;
	for (const p of sorted) {
		if (prev && p - prev > 1) out.push('ellipsis');
		out.push(p);
		prev = p;
	}
	return out;
}