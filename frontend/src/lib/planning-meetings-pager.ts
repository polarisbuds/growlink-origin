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
	deleted: boolean;
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
	if (q.deleted) params.set('deleted', '1');
	return params.toString();
}

export function lastPageForTotal(total: number, pageSize: number): number {
	if (total <= 0) return 1;
	return Math.ceil(total / pageSize);
}

export type PageToken = number | 'ellipsis';

export function visiblePageNumbers(current: number, last: number): PageToken[] {
	if (last <= 1) return [1];
	if (last <= 7) {
		return Array.from({ length: last }, (_, i) => i + 1);
	}

	const needLeftEllipsis = current > 3;
	const needRightEllipsis = current < last - 2;

	if (!needLeftEllipsis) {
		return [1, 2, 3, 4, 5, 'ellipsis', last];
	}
	if (!needRightEllipsis) {
		return [1, 'ellipsis', last - 4, last - 3, last - 2, last - 1, last];
	}
	return [1, 'ellipsis', current - 1, current, current + 1, 'ellipsis', last];
}