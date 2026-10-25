export function formatDate(iso?: string): string {
	if (!iso) return '—';
	const d = new Date(iso);
	if (Number.isNaN(d.getTime())) return '—';
	return d.toLocaleDateString('ja-JP');
}

export type MeetingPairingRef = {
	mentor?: { name?: string };
	mentee?: { name?: string };
	mentees?: { name?: string }[];
};

export type MeetingListItem = {
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
	content?: string;
	topicId?: number;
	topic?: MeetingTopicRef;
};

export type DisplayTopicRow = {
	displayNo: number;
	topic?: MeetingTopicRef;
	answer?: MeetingAnswerRef;
};

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
	canEdit: boolean;
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
