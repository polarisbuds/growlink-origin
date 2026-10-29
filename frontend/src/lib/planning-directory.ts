export type PlanningDirectoryUser = {
	id: number;
	name: string;
	email: string;
	departmentId?: number;
	departmentNames?: string;
};

export type PlanningDirectory = {
	users?: PlanningDirectoryUser[];
	departments?: { id: number; name: string }[];
};

export type UserDepartmentMaps = {
	departmentIdByUserId: Map<number, number>;
	departmentNameByUserId: Map<number, string>;
};

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
