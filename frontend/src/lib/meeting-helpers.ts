export function formatDate(iso?: string): string {
	if (!iso) return '—';
	const d = new Date(iso);
	if (Number.isNaN(d.getTime())) return '—';
	return d.toLocaleDateString('ja-JP');
}

export type MeetingPairingRef = {
	mentor?: { id?: number; name?: string };
	mentee?: { id?: number; name?: string };
	mentees?: { name?: string }[];
};

export type MeetingListItem = {
	id?: number;
	uuid?: string;
	implementedAt?: string;
	status?: number;
	template?: { title?: string };
	pairing?: MeetingPairingRef;
};

export type MeetingTopicRef = {
	id?: number;
	displayNo?: number;
	question?: string;
	description?: string;
};

export type MeetingAnswerRef = {
	id?: number;
	content?: string;
	topicId?: number;
	meetingId?: number;
	userId?: number;
	topic?: MeetingTopicRef;
};

export type MeetingUserRef = {
	id?: number;
	name?: string;
	email?: string;
};

export type MeetingCommentRef = {
	id?: number;
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
	item: { pairing?: MeetingPairingRef },
	userId: number | undefined
): boolean {
	if (!userId || !item.pairing) return false;
	const p = item.pairing;
	return p.mentor?.id === userId || p.mentee?.id === userId;
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
	const answersByTopic = new Map<number, MeetingAnswerRef>();
	for (const a of payload.answers || []) {
		const tid = a.topic?.id ?? a.topicId;
		if (tid) answersByTopic.set(tid, a);
	}
	const topics = payload.topics || [];
	if (topics.length) {
		return topics.map((topic, i) => {
			const answer = topic.id != null ? answersByTopic.get(topic.id) : undefined;
			return { topic, answer, displayNo: topic.displayNo || i + 1 };
		});
	}
	return (payload.answers || []).map((a, i) => ({
		topic: a.topic,
		answer: a,
		displayNo: a.topic?.displayNo || i + 1
	}));
}

function pairingUserIds(payload: MeetingPayload): { mentorId?: number; menteeId?: number } {
	const p = payload.meeting?.pairing;
	return { mentorId: p?.mentor?.id, menteeId: p?.mentee?.id };
}

export function dualTopicRows(payload: MeetingPayload): DualTopicRow[] {
	const { mentorId, menteeId } = pairingUserIds(payload);
	const slots = new Map<
		number,
		{ menteeAnswer?: MeetingAnswerRef; mentorAnswer?: MeetingAnswerRef }
	>();
	for (const a of payload.answers || []) {
		const tid = a.topic?.id ?? a.topicId;
		if (tid == null) continue;
		let slot = slots.get(tid);
		if (!slot) {
			slot = {};
			slots.set(tid, slot);
		}
		const uid = a.userId;
		if (uid != null && uid === mentorId) {
			slot.mentorAnswer = a;
		} else if (uid != null && uid === menteeId) {
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
			const slot = topic.id != null ? slots.get(topic.id) : undefined;
			return {
				topic,
				displayNo: topic.displayNo || i + 1,
				menteeAnswer: slot?.menteeAnswer,
				mentorAnswer: slot?.mentorAnswer
			};
		});
	}
	return [...slots.entries()].map(([topicId, slot], i) => ({
		displayNo: i + 1,
		topic: slot.menteeAnswer?.topic ?? slot.mentorAnswer?.topic ?? { id: topicId },
		menteeAnswer: slot.menteeAnswer,
		mentorAnswer: slot.mentorAnswer
	}));
}

export function ownAnswerDrafts(
	payload: MeetingPayload,
	viewerRole: 'mentee' | 'mentor' | 'viewer' | undefined
): Record<number, string> {
	const rows = dualTopicRows(payload);
	const next: Record<number, string> = {};
	for (const row of rows) {
		const tid = row.topic?.id;
		if (tid == null) continue;
		const answer =
			viewerRole === 'mentor'
				? row.mentorAnswer
				: viewerRole === 'mentee'
					? row.menteeAnswer
					: undefined;
		next[tid] = answer?.content || '';
	}
	return next;
}
