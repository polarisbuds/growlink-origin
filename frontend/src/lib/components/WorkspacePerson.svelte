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
    let department = $derived(departmentNameForUser(userUuid, deptMaps) || "—");
</script>

<div class="min-w-0 gap-2 flex flex-1 w-full">
    <span class="flex items-center">
        <Icon
            url={avatarURL(directoryUser)}
            text={displayName}
            width={28}
            height={28}
            alt={displayName} />
    </span>
    <div class="flex flex-1 flex-col">
        <span class="truncate text-sm">{displayName}</span>
        {#if showDepartment}
            <p class="text-xs text-base-content/60 truncate">
                {department}
            </p>
        {/if}
    </div>
</div>
