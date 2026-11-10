import type { PairingListItem } from '$lib/pairing-helpers';

export type MeetingTargetPick = { kind: string; targetUuid: string };

export type MeetingPairingItem = PairingListItem & {
	mentor?: { uuid?: string; name?: string; imageUri?: string };
};

function pairingMatchesUserTarget(p: MeetingPairingItem, userUuid: string): boolean {
	return p.mentee?.uuid === userUuid;
}

export function matchPairingUuids(
	pairings: MeetingPairingItem[],
	mentorUuids: string[],
	targets: MeetingTargetPick[]
): string[] {
	if (!mentorUuids.length || !targets.length) return [];
	const userTargetUuids = targets.filter((t) => t.kind === 'user').map((t) => t.targetUuid);
	if (!userTargetUuids.length) return [];
	const mentorSet = new Set(mentorUuids);
	const uuids = new Set<string>();
	for (const p of pairings) {
		const mid = p.mentor?.uuid;
		if (!mid || !mentorSet.has(mid)) continue;
		const matched = userTargetUuids.some((uid) => pairingMatchesUserTarget(p, uid));
		if (matched) uuids.add(p.uuid);
	}
	return [...uuids];
}
