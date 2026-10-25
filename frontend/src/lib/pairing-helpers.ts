export type PairingSection = 'users' | 'departments' | 'teams';

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
	mentorName: string;
	counterpartyName: string;
};

const SECTION_TARGET_KIND: Record<Exclude<PairingSection, 'users'>, string> = {
	departments: 'department',
	teams: 'team'
};

export function sectionSecondColumnLabel(section: PairingSection): string {
	if (section === 'users') return 'メンティ';
	if (section === 'departments') return '部署';
	return 'チーム';
}

export function sectionListHref(section: PairingSection): string {
	return `/planning/pairs/${section}`;
}

export function targetKindForSection(section: PairingSection): string {
	if (section === 'users') return 'user';
	if (section === 'departments') return 'department';
	return 'team';
}

export function sectionTargetFieldLabel(section: PairingSection): string {
	if (section === 'users') return 'メンティ';
	if (section === 'departments') return '部署';
	return 'チーム';
}

export function parsePairingSection(value: string | null | undefined): PairingSection {
	if (value === 'departments' || value === 'teams') return value;
	return 'users';
}

export function expandPairingsForSection(
	pairings: PairingListItem[],
	section: PairingSection
): PairingDisplayRow[] {
	const rows: PairingDisplayRow[] = [];
	for (const p of pairings) {
		const mentorName = p.mentor?.name || '—';
		if (section === 'users') {
			const userTargets = (p.targets || []).filter((t) => t.kind === 'user');
			if (userTargets.length) {
				for (const t of userTargets) {
					rows.push({
						pairing: p,
						rowKey: `${p.id}-user-${t.targetId}`,
						mentorName,
						counterpartyName: t.label || '—'
					});
				}
			} else if (p.mentee?.name) {
				rows.push({
					pairing: p,
					rowKey: `${p.id}-mentee-${p.mentee.id ?? 0}`,
					mentorName,
					counterpartyName: p.mentee.name
				});
			}
			continue;
		}
		const kind = SECTION_TARGET_KIND[section];
		for (const t of (p.targets || []).filter((x) => x.kind === kind)) {
			rows.push({
				pairing: p,
				rowKey: `${p.id}-${kind}-${t.targetId}`,
				mentorName,
				counterpartyName: t.label || '—'
			});
		}
	}
	return rows;
}

export function rowSearchHaystack(row: PairingDisplayRow): string {
	return [row.mentorName, row.counterpartyName].filter(Boolean).join(' ').toLowerCase();
}
