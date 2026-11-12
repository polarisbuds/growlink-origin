export type PairingListItem = {
	uuid: string;
	mentor?: { uuid?: string; name?: string; imageUri?: string; departments?: string; teams?: string };
	mentee?: { uuid?: string; name?: string; imageUri?: string; departments?: string; teams?: string };
};

export type PairingDisplayRow = {
	pairing: PairingListItem;
	rowKey: string;
	mentorUserUuid: string;
	menteeUserUuid: string;
	mentorName: string;
	menteeName: string;
	mentorDepartments: string;
	menteeDepartments: string;
	mentorTeams: string;
	menteeTeams: string;
};

export function expandUserPairings(pairings: PairingListItem[]): PairingDisplayRow[] {
	const rows: PairingDisplayRow[] = [];
	for (const p of pairings) {
		if (!p.mentee?.uuid || !p.mentor?.uuid) {
			continue;
		}

		rows.push({
			pairing: p,
			rowKey: `${p.uuid}-mentee-${p.mentee.uuid}`,
			mentorUserUuid: p.mentor.uuid || '',
			menteeUserUuid: p.mentee.uuid,
			mentorName: p.mentor.name || '—',
			menteeName: p.mentee.name || '—',
			mentorDepartments: p.mentor.departments || '',
			menteeDepartments: p.mentee.departments || '',
			mentorTeams: p.mentor.teams || '',
			menteeTeams: p.mentee.teams || '',

		});
	}
	return rows;
}
