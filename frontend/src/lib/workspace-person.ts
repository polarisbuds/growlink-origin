import type { PlanningDirectoryUser } from "$lib/planning-directory";

export function avatarURL(user: PlanningDirectoryUser | undefined): string {
    if (!user?.imageUri) return "";
    return `/uploads/${user.imageUri}`;
}

export function findDirectoryUser(
    users: PlanningDirectoryUser[],
    uuid: string | undefined,
): PlanningDirectoryUser | undefined {
    if (!uuid) return undefined;
    return users.find((u) => u.uuid === uuid);
}
