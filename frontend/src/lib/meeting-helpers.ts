export function formatDate(iso?: string): string {
	if (!iso) return '—';
	const d = new Date(iso);
	if (Number.isNaN(d.getTime())) return '—';
	return d.toLocaleDateString('ja-JP');
}

export type MeetingPairingRef = {
	mentor?: { uuid?: string; name?: string };
	mentee?: { uuid?: string; name?: string };
	mentees?: { name?: string }[];
};

export type MeetingListItem = {
	uuid: string;
	implementedAt?: string;
	status?: number;
	template?: { title?: string; description?: string };
	pairing?: MeetingPairingRef;
};

export type MeetingTopicRef = {
	uuid?: string;
	displayNo?: number;
	question?: string;
	description?: string;
};

export type MeetingAnswerRef = {
	uuid?: string;
	content?: string;
	topic?: MeetingTopicRef;
	user?: { uuid?: string; name?: string; email?: string };
};

export type MeetingUserRef = {
	uuid?: string;
	name?: string;
	email?: string;
};

export type MeetingCommentRef = {
	uuid?: string;
	content?: string;
	isLimited?: boolean;
	createdAt?: string;
	user?: { name?: string };
};

export type DisplayTopicRow = {
	displayNo: number;
	topic?: MeetingTopicRef;
	answer?: MeetingAnswerRef;
};

export type DualTopicRow = {
	displayNo: number;
	topic?: MeetingTopicRef;
	menteeAnswer?: MeetingAnswerRef;
	mentorAnswer?: MeetingAnswerRef;
};

export const MEETING_STATUS_DONE = 3;

export function isMeetingDone(status?: number): boolean {
	return status === MEETING_STATUS_DONE;
}

export function menteeTargetLabel(pairing?: MeetingPairingRef): string {
	if (!pairing) return '—';
	const mentees = (pairing.mentees || []).map((m) => m.name).filter(Boolean);
	if (mentees.length) return mentees.join('・');
	return pairing.mentee?.name || '—';
}

export function canUserAnswerMeeting(
	item: { pairing?: MeetingPairingRef; status?: number },
	userUuid: string | undefined
): boolean {
	if (!userUuid || !item.pairing) return false;
	const done = isMeetingDone(item.status);
	if (isUserMentorOfMeeting(item, userUuid)) {
		return done;
	}
	if (isUserMenteeOfMeeting(item, userUuid)) {
		return !done;
	}
	return false;
}

/** 一覧の「閲覧」。記入と同時には true にならない */
export function canUserViewMeeting(
	item: { pairing?: MeetingPairingRef; status?: number },
	userUuid: string | undefined
): boolean {
	if (!item.pairing) return false;
	const done = isMeetingDone(item.status);
	if (userUuid && isUserMenteeOfMeeting(item, userUuid)) {
		return done;
	}
	if (userUuid && isUserMentorOfMeeting(item, userUuid)) {
		return !done;
	}
	return true;
}

/** メンター・共有者はメンティー提出前に面談の中身を画面に出さない */
export function shouldHideMeetingContentBeforeSubmit(
	viewerRole: 'mentee' | 'mentor' | 'viewer' | undefined,
	status?: number
): boolean {
	if (viewerRole === 'mentee') return false;
	return !isMeetingDone(status);
}

export function pairLabel(item: { pairing?: MeetingPairingRef }): string {
	const p = item.pairing;
	if (!p) return '';
	const mentees = (p.mentees || []).map((m) => m.name).filter(Boolean);
	const menteeLabel = mentees.length ? mentees.join('・') : p.mentee?.name || '';
	return `${p.mentor?.name || ''} / ${menteeLabel}`;
}

export function splitByToday<T extends { implementedAt?: string }>(meetings: T[]) {
	const today = new Date();
	today.setHours(23, 59, 59, 999);
	const untilToday: T[] = [];
	const fromTomorrow: T[] = [];
	for (const m of meetings) {
		const d = new Date(m.implementedAt || '');
		if (d <= today) untilToday.push(m);
		else fromTomorrow.push(m);
	}
	return { untilToday, fromTomorrow };
}

export type MeetingPeriodSplit<T> = {
	thisMonth: T[];
	past: T[];
	planned: T[];
};

function monthBounds(now = new Date()) {
	const y = now.getFullYear();
	const m = now.getMonth();
	const thisMonthStart = new Date(y, m, 1, 0, 0, 0, 0);
	const nextMonthStart = new Date(y, m + 1, 1, 0, 0, 0, 0);
	return { thisMonthStart, nextMonthStart };
}

export function splitByMeetingPeriod<T extends { implementedAt?: string }>(
	meetings: T[],
	now = new Date()
): MeetingPeriodSplit<T> {
	const { thisMonthStart, nextMonthStart } = monthBounds(now);
	const thisMonth: T[] = [];
	const past: T[] = [];
	const planned: T[] = [];
	for (const meeting of meetings) {
		const raw = meeting.implementedAt;
		if (!raw) {
			thisMonth.push(meeting);
			continue;
		}
		const d = new Date(raw);
		if (Number.isNaN(d.getTime())) {
			thisMonth.push(meeting);
			continue;
		}
		if (d < thisMonthStart) {
			past.push(meeting);
		} else if (d >= nextMonthStart) {
			planned.push(meeting);
		} else {
			thisMonth.push(meeting);
		}
	}
	return { thisMonth, past, planned };
}

export function isUserMentorOfMeeting(
	item: { pairing?: MeetingPairingRef },
	userUuid: string | undefined
): boolean {
	if (!userUuid || !item.pairing) return false;
	return item.pairing.mentor?.uuid === userUuid;
}

export function isUserMenteeOfMeeting(
	item: { pairing?: MeetingPairingRef },
	userUuid: string | undefined
): boolean {
	if (!userUuid || !item.pairing) return false;
	return item.pairing.mentee?.uuid === userUuid;
}

export const BADGE_YOU_ARE_MENTEE = 'あなたがメンティーの面談';
export const BADGE_YOU_ARE_MENTOR = 'あなたがメンターの面談';

export function meetingRoleBadgeLabel(
	item: { pairing?: MeetingPairingRef },
	userUuid: string | undefined,
	mode: 'mentee-only' | 'mentor-always' | 'none'
): string | null {
	if (mode === 'none' || !userUuid) return null;
	if (mode === 'mentor-always' || isUserMentorOfMeeting(item, userUuid)) {
		if (isUserMentorOfMeeting(item, userUuid)) return BADGE_YOU_ARE_MENTOR;
	}
	if (mode === 'mentee-only' && isUserMenteeOfMeeting(item, userUuid)) {
		return BADGE_YOU_ARE_MENTEE;
	}
	return null;
}

export function canShowMeetingSetting(implementedAt?: string): boolean {
	if (!implementedAt) return true;
	const d = new Date(implementedAt);
	if (Number.isNaN(d.getTime())) return true;
	const today = new Date();
	today.setHours(23, 59, 59, 999);
	return d <= today;
}

export function splitInvolvedMeetingsForWorkspace(
	meetings: MeetingListItem[],
	userUuid: string | undefined
): { own: MeetingListItem[]; asMentor: MeetingListItem[]; showMentorTab: boolean } {
	const asMentor: MeetingListItem[] = [];
	const own: MeetingListItem[] = [];
	for (const m of meetings) {
		if (isUserMentorOfMeeting(m, userUuid)) {
			asMentor.push(m);
		} else {
			own.push(m);
		}
	}
	return { own, asMentor, showMentorTab: asMentor.length > 0 };
}

function topicKey(topic?: MeetingTopicRef): string | undefined {
	if (!topic) return undefined;
	return topic.uuid || (topic.displayNo != null ? `display-${topic.displayNo}` : undefined);
}

export function mentorTabPairLabel(item: { pairing?: MeetingPairingRef }): string {
	const mentee = menteeTargetLabel(item.pairing);
	const mentor = item.pairing?.mentor?.name || '';
	if (mentee && mentee !== '—') return `メンティー: ${mentee}`;
	return pairLabel(item) || mentor;
}

export type MeetingPayload = {
	uuid: string;
	meeting: MeetingListItem;
	answers: MeetingAnswerRef[];
	topics: MeetingTopicRef[];
	meetingComments?: MeetingCommentRef[];
	canEdit: boolean;
	viewerRole?: 'mentee' | 'mentor' | 'viewer';
	canWriteAnswers?: boolean;
	canEditSharers?: boolean;
	sharers?: MeetingUserRef[];
};

export function displayTopics(payload: MeetingPayload): DisplayTopicRow[] {
	const answersByTopic = new Map<string, MeetingAnswerRef>();
	for (const a of payload.answers || []) {
		const key = topicKey(a.topic);
		if (key) answersByTopic.set(key, a);
	}
	const topics = payload.topics || [];
	if (topics.length) {
		return topics.map((topic, i) => {
			const key = topicKey(topic);
			const answer = key ? answersByTopic.get(key) : undefined;
			return { topic, answer, displayNo: topic.displayNo || i + 1 };
		});
	}
	return (payload.answers || []).map((a, i) => ({
		topic: a.topic,
		answer: a,
		displayNo: a.topic?.displayNo || i + 1
	}));
}

function pairingUserUuids(payload: MeetingPayload): { mentorUuid?: string; menteeUuid?: string } {
	const p = payload.meeting?.pairing;
	return { mentorUuid: p?.mentor?.uuid, menteeUuid: p?.mentee?.uuid };
}

export function dualTopicRows(payload: MeetingPayload): DualTopicRow[] {
	const { mentorUuid, menteeUuid } = pairingUserUuids(payload);
	const slots = new Map<
		string,
		{ menteeAnswer?: MeetingAnswerRef; mentorAnswer?: MeetingAnswerRef }
	>();
	for (const a of payload.answers || []) {
		const key = topicKey(a.topic);
		if (!key) continue;
		let slot = slots.get(key);
		if (!slot) {
			slot = {};
			slots.set(key, slot);
		}
		const uid = a.user?.uuid;
		if (uid && uid === mentorUuid) {
			slot.mentorAnswer = a;
		} else if (uid && uid === menteeUuid) {
			slot.menteeAnswer = a;
		} else if (!slot.menteeAnswer) {
			slot.menteeAnswer = a;
		} else if (!slot.mentorAnswer) {
			slot.mentorAnswer = a;
		}
	}
	const topics = payload.topics || [];
	if (topics.length) {
		return topics.map((topic, i) => {
			const key = topicKey(topic);
			const slot = key ? slots.get(key) : undefined;
			return {
				topic,
				displayNo: topic.displayNo || i + 1,
				menteeAnswer: slot?.menteeAnswer,
				mentorAnswer: slot?.mentorAnswer
			};
		});
	}
	return [...slots.entries()].map(([key, slot], i) => ({
		displayNo: i + 1,
		topic: slot.menteeAnswer?.topic ?? slot.mentorAnswer?.topic ?? { uuid: key },
		menteeAnswer: slot.menteeAnswer,
		mentorAnswer: slot.mentorAnswer
	}));
}

export function ownAnswerDrafts(
	payload: MeetingPayload,
	viewerRole: 'mentee' | 'mentor' | 'viewer' | undefined
): Record<string, string> {
	const rows = dualTopicRows(payload);
	const next: Record<string, string> = {};
	for (const row of rows) {
		const key = topicKey(row.topic);
		if (!key) continue;
		const answer =
			viewerRole === 'mentor'
				? row.mentorAnswer
				: viewerRole === 'mentee'
					? row.menteeAnswer
					: undefined;
		next[key] = answer?.content || '';
	}
	return next;
}
