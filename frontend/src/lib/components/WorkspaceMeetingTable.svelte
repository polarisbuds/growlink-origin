<script lang="ts">
    import { statusBadge, statusLabel } from "$lib/status";
    import {
        formatDate,
        canUserAnswerMeeting,
        canUserViewMeeting,
        type MeetingListItem,
    } from "$lib/meeting-helpers";
    import {
        buildUserDepartmentMaps,
        departmentNameForUser,
        type PlanningDirectoryUser,
    } from "$lib/planning-directory";
    import WorkspacePerson from "$lib/components/WorkspacePerson.svelte";

    export type WorkspacePersonMode = "mentor" | "mentee" | "both";

    let {
        items,
        currentUserUuid,
        directoryUsers,
        personMode,
        emptyMessage = "面談はありません",
    }: {
        items: MeetingListItem[];
        currentUserUuid?: string;
        directoryUsers: PlanningDirectoryUser[];
        personMode: WorkspacePersonMode;
        emptyMessage?: string;
    } = $props();

    let deptMaps = $derived(buildUserDepartmentMaps(directoryUsers));
</script>

<div
    class="overflow-x-auto rounded-3xl border border-base-300 bg-base-100 shadow-sm">
    <table class="table">
        <thead>
            <tr class="text-base-content/70">
                {#if personMode === "both"}
                    <th>メンター部署</th>
                    <th>メンター</th>
                    <th>メンティー部署</th>
                    <th>メンティー</th>
                {:else if personMode === "mentor"}
                    <th>部署</th>
                    <th>メンター</th>
                {:else}
                    <th>部署</th>
                    <th>メンティー</th>
                {/if}
                <th>テンプレート</th>
                <th>状態</th>
                <th>実施予定日</th>
                <th class="text-right">操作</th>
            </tr>
        </thead>
        <tbody>
            {#each items as item (item.uuid)}
                <tr class="hover:bg-base-200/40">
                    {#if personMode === "both"}
                        <td class="text-sm whitespace-nowrap"
                            >{departmentNameForUser(
                                item.pairing?.mentor?.uuid,
                                deptMaps,
                            ) || "—"}</td>
                        <td>
                            <WorkspacePerson
                                userUuid={item.pairing?.mentor?.uuid}
                                name={item.pairing?.mentor?.name}
                                {directoryUsers}
                                {deptMaps}
                                showDepartment={false} />
                        </td>
                        <td class="text-sm whitespace-nowrap"
                            >{departmentNameForUser(
                                item.pairing?.mentee?.uuid,
                                deptMaps,
                            ) || "—"}</td>
                        <td>
                            <WorkspacePerson
                                userUuid={item.pairing?.mentee?.uuid}
                                name={item.pairing?.mentee?.name}
                                {directoryUsers}
                                {deptMaps}
                                showDepartment={false} />
                        </td>
                    {:else if personMode === "mentor"}
                        <td class="text-sm whitespace-nowrap"
                            >{departmentNameForUser(
                                item.pairing?.mentor?.uuid,
                                deptMaps,
                            ) || "—"}</td>
                        <td>
                            <WorkspacePerson
                                userUuid={item.pairing?.mentor?.uuid}
                                name={item.pairing?.mentor?.name}
                                {directoryUsers}
                                {deptMaps}
                                showDepartment={false} />
                        </td>
                    {:else}
                        <td class="text-sm whitespace-nowrap"
                            >{departmentNameForUser(
                                item.pairing?.mentee?.uuid,
                                deptMaps,
                            ) || "—"}</td>
                        <td>
                            <WorkspacePerson
                                userUuid={item.pairing?.mentee?.uuid}
                                name={item.pairing?.mentee?.name}
                                {directoryUsers}
                                {deptMaps}
                                showDepartment={false} />
                        </td>
                    {/if}
                    <td>{item.template?.title || "1on1"}</td>
                    <td>
                        <span
                            class="badge badge-sm {statusBadge(
                                item.status ?? 0,
                            )}">
                            {statusLabel(item.status ?? 0)}
                        </span>
                    </td>
                    <td class="whitespace-nowrap"
                        >{formatDate(item.implementedAt)}</td>
                    <td class="text-right">
                        <div class="flex flex-wrap justify-end gap-1">
                            {#if canUserAnswerMeeting(item, currentUserUuid)}
                                <a
                                    href="/meeting/{item.uuid}/answer"
                                    class="btn btn-xs btn-primary">記入</a>
                            {:else if canUserViewMeeting(item, currentUserUuid)}
                                <a
                                    href="/meeting/{item.uuid}"
                                    class="btn btn-xs btn-ghost">閲覧</a>
                            {/if}
                        </div>
                    </td>
                </tr>
            {/each}
        </tbody>
    </table>
    {#if items.length === 0}
        <p class="p-8 text-center text-sm text-base-content/50">
            {emptyMessage}
        </p>
    {/if}
</div>
