import {
	lastPageForTotal,
	visiblePageNumbers,
	type PageToken,
} from '$lib/planning-meetings-pager';

export type ManagementUsersPageResponse = {
	items: unknown[];
	total: number;
	page: number;
	pageSize: number;
};

export type ManagementUserListQuery = {
	page: number;
	name: string;
	email: string;
	departmentUuids: string[];
};

export function buildManagementUsersQuery(q: ManagementUserListQuery): string {
	const params = new URLSearchParams();
	params.set('page', String(q.page));
	const name = q.name.trim();
	const email = q.email.trim();
	if (name) params.set('name', name);
	if (email) params.set('email', email);
	if (q.departmentUuids.length) {
		params.set('departmentUuids', q.departmentUuids.join(','));
	}
	return params.toString();
}

export { lastPageForTotal, visiblePageNumbers, type PageToken };