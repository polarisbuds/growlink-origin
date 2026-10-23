export function formatDate(iso?: string): string {
	if (!iso) return '—';
	const d = new Date(iso);
	if (Number.isNaN(d.getTime())) return '—';
	return d.toLocaleDateString('ja-JP');
}

export function pairLabel(item: {
	pairing?: {
		mentor?: { name?: string };
		mentee?: { name?: string };
		mentees?: { name?: string }[];
	};
}): string {
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
	meeting: Record<string, unknown>;
	answers: Record<string, unknown>[];
	topics: Record<string, unknown>[];
	canEdit: boolean;
};

export function displayTopics(payload: MeetingPayload) {
	const answersByTopic = new Map<number, Record<string, unknown>>();
	for (const a of payload.answers || []) {
		const topic = a.topic as { id?: number } | undefined;
		const tid = topic?.id ?? (a.topicId as number);
		if (tid) answersByTopic.set(tid, a);
	}
	const topics = payload.topics || [];
	if (topics.length) {
		return topics.map((topic, i) => {
			const t = topic as { id?: number; displayNo?: number };
			const answer = answersByTopic.get(t.id!);
			return { topic, answer, displayNo: t.displayNo || i + 1 };
		});
	}
	return (payload.answers || []).map((a, i) => {
		const topic = a.topic as { displayNo?: number };
		return { topic: a.topic, answer: a, displayNo: topic?.displayNo || i + 1 };
	});
}
