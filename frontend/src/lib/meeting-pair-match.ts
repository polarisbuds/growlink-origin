import type { PairingListItem } from '$lib/pairing-helpers';

export type MeetingTargetPick = { kind: string; targetId: number };

export type MeetingPairingItem = PairingListItem & {
	mentor?: { id?: number; name?: string };
};

function pairingMatchesUserTarget(p: MeetingPairingItem, userId: number): boolean {
	const userTargets = (p.targets || []).filter((t) => t.kind === 'user');
	if (userTargets.some((t) => t.targetId === userId)) return true;
	return p.mentee?.id === userId;
}

export function matchPairingIds(
	pairings: MeetingPairingItem[],
	mentorIds: number[],
	targets: MeetingTargetPick[]
): number[] {
	if (!mentorIds.length || !targets.length) return [];
	const userTargetIds = targets.filter((t) => t.kind === 'user').map((t) => t.targetId);
	if (!userTargetIds.length) return [];
	const mentorSet = new Set(mentorIds);
	const ids = new Set<number>();
	for (const p of pairings) {
		const mid = p.mentor?.id;
		if (!mid || !mentorSet.has(mid)) continue;
		const matched = userTargetIds.some((uid) => pairingMatchesUserTarget(p, uid));
		if (matched) ids.add(p.id);
	}
	return [...ids];
}
