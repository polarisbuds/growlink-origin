export type PairingListItem = {
	id: number;
	uuid: string;
	mentor?: { id?: number; name?: string };
	mentee?: { id?: number; name?: string };
	targets?: { kind: string; targetId: number; label?: string }[];
};

export type PairingDisplayRow = {
	pairing: PairingListItem;
	rowKey: string;
	mentorUserId: number;
	menteeUserId: number;
	mentorName: string;
	menteeName: string;
};

export function expandUserPairings(pairings: PairingListItem[]): PairingDisplayRow[] {
	const rows: PairingDisplayRow[] = [];
	for (const p of pairings) {
		const mentorUserId = p.mentor?.id || 0;
		const mentorName = p.mentor?.name || '—';
		const userTargets = (p.targets || []).filter((t) => t.kind === 'user');
		if (userTargets.length) {
			for (const t of userTargets) {
				rows.push({
					pairing: p,
					rowKey: `${p.id}-user-${t.targetId}`,
					mentorUserId,
					menteeUserId: t.targetId,
					mentorName,
					menteeName: t.label || '—'
				});
			}
		} else if (p.mentee?.id) {
			rows.push({
				pairing: p,
				rowKey: `${p.id}-mentee-${p.mentee.id}`,
				mentorUserId,
				menteeUserId: p.mentee.id,
				mentorName,
				menteeName: p.mentee.name || '—'
			});
		}
	}
	return rows;
}
