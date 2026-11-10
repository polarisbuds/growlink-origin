export type PairingListItem = {
	uuid: string;
	mentor?: { uuid?: string; name?: string; imageUri?: string };
	mentee?: { uuid?: string; name?: string; imageUri?: string };
};

export type PairingDisplayRow = {
	pairing: PairingListItem;
	rowKey: string;
	mentorUserUuid: string;
	menteeUserUuid: string;
	mentorName: string;
	menteeName: string;
};

export function expandUserPairings(pairings: PairingListItem[]): PairingDisplayRow[] {
	const rows: PairingDisplayRow[] = [];
	for (const p of pairings) {
		const mentorUserUuid = p.mentor?.uuid || '';
		const mentorName = p.mentor?.name || '—';
		if (!p.mentee?.uuid) {
			continue;
		}
		rows.push({
			pairing: p,
			rowKey: `${p.uuid}-mentee-${p.mentee.uuid}`,
			mentorUserUuid,
			menteeUserUuid: p.mentee.uuid,
			mentorName,
			menteeName: p.mentee.name || '—'
		});
	}
	return rows;
}
