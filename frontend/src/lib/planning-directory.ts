export type PlanningDirectoryUser = {
	uuid: string;
	name: string;
	email: string;
	imageUri?: string;
	teamUuid?: string;
	teamNames?: string;
	departmentUuid?: string;
	departmentNames?: string;
};

export type PlanningDirectory = {
	users?: PlanningDirectoryUser[];
	departments?: { uuid: string; name: string }[];
};

export type UserTeamMaps = {
	teamUuidByUserUuid: Map<string, string>;
	teamNameByUserUuid: Map<string, string>;
};

export type UserDepartmentMaps = {
	departmentUuidByUserUuid: Map<string, string>;
	departmentNameByUserUuid: Map<string, string>;
};

export function buildUserTeamMaps(users: PlanningDirectoryUser[]): UserTeamMaps {
	const teamUuidByUserUuid = new Map<string, string>();
	const teamNameByUserUuid = new Map<string, string>();
	for (const u of users) {
		if (u.teamUuid) teamUuidByUserUuid.set(u.uuid, u.teamUuid);
		const name = (u.teamNames || '').trim();
		if (name) teamNameByUserUuid.set(u.uuid, name);
	}
	return { teamUuidByUserUuid, teamNameByUserUuid };
}

export function buildUserDepartmentMaps(users: PlanningDirectoryUser[]): UserDepartmentMaps {
	const departmentUuidByUserUuid = new Map<string, string>();
	const departmentNameByUserUuid = new Map<string, string>();
	for (const u of users) {
		if (u.departmentUuid) departmentUuidByUserUuid.set(u.uuid, u.departmentUuid);
		const name = (u.departmentNames || '').trim();
		if (name) departmentNameByUserUuid.set(u.uuid, name);
	}
	return { departmentUuidByUserUuid, departmentNameByUserUuid };
}

export function teamNameForUser(
	userUuid: string | undefined,
	maps: UserTeamMaps
): string {
	if (!userUuid) return '';
	return maps.teamNameByUserUuid.get(userUuid) || '';
}

export function teamUuidForUser(
	userUuid: string | undefined,
	maps: UserTeamMaps
): string {
	if (!userUuid) return '';
	return maps.teamUuidByUserUuid.get(userUuid) || '';
}

export function departmentNameForUser(
	userUuid: string | undefined,
	maps: UserDepartmentMaps
): string {
	if (!userUuid) return '';
	return maps.departmentNameByUserUuid.get(userUuid) || '';
}

export function departmentUuidForUser(
	userUuid: string | undefined,
	maps: UserDepartmentMaps
): string {
	if (!userUuid) return '';
	return maps.departmentUuidByUserUuid.get(userUuid) || '';
}
