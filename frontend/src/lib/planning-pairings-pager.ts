import {
	lastPageForTotal,
	visiblePageNumbers,
	type MeetingsPageResponse,
	type PageToken,
} from '$lib/planning-meetings-pager';

export type PairingsPageResponse = MeetingsPageResponse;

export type PairingListQuery = {
	page?: number;
	mentorUuids: string[];
	menteeUuids: string[];
	departmentUuids: string[];
	teamUuids: string[];
};

export function buildPairingsQuery(q: PairingListQuery): string {
	const params = new URLSearchParams();
	if (q.page != null) params.set('page', String(q.page));
	if (q.mentorUuids.length) params.set('mentorUuids', q.mentorUuids.join(','));
	if (q.menteeUuids.length) params.set('menteeUuids', q.menteeUuids.join(','));
	if (q.departmentUuids.length) params.set('departmentUuids', q.departmentUuids.join(','));
	if (q.teamUuids.length) params.set('teamUuids', q.teamUuids.join(','));
	return params.toString();
}

export { lastPageForTotal, visiblePageNumbers, type PageToken };