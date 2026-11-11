<script lang="ts">
    import Icon from "$lib/components/Icon.svelte";
    import {
        departmentNameForUser,
        type PlanningDirectoryUser,
        type UserDepartmentMaps,
    } from "$lib/planning-directory";
    import { avatarURL, findDirectoryUser } from "$lib/workspace-person";

    let {
        userUuid,
        name = "",
        directoryUsers,
        deptMaps,
        showDepartment = true,
    }: {
        userUuid?: string;
        name?: string;
        directoryUsers: PlanningDirectoryUser[];
        deptMaps: UserDepartmentMaps;
        showDepartment?: boolean;
    } = $props();

    let directoryUser = $derived(findDirectoryUser(directoryUsers, userUuid));
    let displayName = $derived(name || directoryUser?.name || "—");
    let department = $derived(
        departmentNameForUser(userUuid, deptMaps) || "—",
    );
</script>

<div class="min-w-0">
    <span class="flex gap-2 items-center">
        <Icon
            url={avatarURL(directoryUser)}
            text={displayName}
            width={28}
            height={28}
            alt={displayName} />
        <span class="truncate text-sm">{displayName}</span>
    </span>
    {#if showDepartment}
        <p class="mt-1 pl-9 text-xs text-base-content/60 truncate">
            {department}
        </p>
    {/if}
</div>
