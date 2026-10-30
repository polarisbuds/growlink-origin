export type PlanningDirectoryUser = {
	id: number;
	name: string;
	email: string;
	teamId?: number;
	teamNames?: string;
	departmentId?: number;
	departmentNames?: string;
};

export type PlanningDirectory = {
	users?: PlanningDirectoryUser[];
	departments?: { id: number; name: string }[];
};

export type UserTeamMaps = {
	teamIdByUserId: Map<number, number>;
	teamNameByUserId: Map<number, string>;
};

export type UserDepartmentMaps = {
	departmentIdByUserId: Map<number, number>;
	departmentNameByUserId: Map<number, string>;
};

export function buildUserTeamMaps(users: PlanningDirectoryUser[]): UserTeamMaps {
	const teamIdByUserId = new Map<number, number>();
	const teamNameByUserId = new Map<number, string>();
	for (const u of users) {
		if (u.teamId) teamIdByUserId.set(u.id, u.teamId);
		const name = (u.teamNames || '').trim();
		if (name) teamNameByUserId.set(u.id, name);
	}
	return { teamIdByUserId: teamIdByUserId, teamNameByUserId: teamNameByUserId };
}

export function buildUserDepartmentMaps(users: PlanningDirectoryUser[]): UserDepartmentMaps {
	const departmentIdByUserId = new Map<number, number>();
	const departmentNameByUserId = new Map<number, string>();
	for (const u of users) {
		if (u.departmentId) departmentIdByUserId.set(u.id, u.departmentId);
		const name = (u.departmentNames || '').trim();
		if (name) departmentNameByUserId.set(u.id, name);
	}
	return { departmentIdByUserId, departmentNameByUserId };
}

export function teamNameForUser(
	userId: number | undefined,
	maps: UserTeamMaps
): string {
	if (!userId) return '';
	return maps.teamNameByUserId.get(userId) || '';
}

export function teamIdForUser(
	userId: number | undefined,
	maps: UserTeamMaps
): number {
	if (!userId) return 0;
	return maps.teamIdByUserId.get(userId) || 0;
}

export function departmentNameForUser(
	userId: number | undefined,
	maps: UserDepartmentMaps
): string {
	if (!userId) return '';
	return maps.departmentNameByUserId.get(userId) || '';
}

export function departmentIdForUser(
	userId: number | undefined,
	maps: UserDepartmentMaps
): number {
	if (!userId) return 0;
	return maps.departmentIdByUserId.get(userId) || 0;
}
