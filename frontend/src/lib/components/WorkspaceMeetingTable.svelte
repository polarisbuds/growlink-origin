<script lang="ts">
    import { statusBadge, statusLabel } from "$lib/status";
    import {
        formatDate,
        pairLabel,
        canUserAnswerMeeting,
        canUserViewMeeting,
        meetingRoleBadgeLabel,
        type MeetingListItem,
    } from "$lib/meeting-helpers";

    type RoleBadgeMode = "mentee-only" | "mentor-always" | "none";

    let {
        items,
        currentUserUuid,
        roleBadgeMode = "none",
        labelForItem = pairLabel,
        emptyMessage = "面談はありません",
    }: {
        items: MeetingListItem[];
        currentUserUuid?: string;
        roleBadgeMode?: RoleBadgeMode;
        labelForItem?: (item: MeetingListItem) => string;
        emptyMessage?: string;
    } = $props();

    let showRoleColumn = $derived(roleBadgeMode !== "none");
</script>

<div
    class="overflow-x-auto rounded-3xl border border-base-300 bg-base-100 shadow-sm">
    <table class="table">
        <thead>
            <tr class="text-base-content/70">
                {#if showRoleColumn}
                    <th>区分</th>
                {/if}
                <th>ペア</th>
                <th>テンプレート</th>
                <th>状態</th>
                <th>実施予定日</th>
                <th class="text-right">操作</th>
            </tr>
        </thead>
        <tbody>
            {#each items as item (item.uuid)}
                {@const roleBadge = meetingRoleBadgeLabel(
                    item,
                    currentUserUuid,
                    roleBadgeMode,
                )}
                <tr class="hover:bg-base-200/40">
                    {#if showRoleColumn}
                        <td class="max-w-[11rem] text-xs">
                            {#if roleBadge}
                                <span
                                    class="badge badge-outline badge-sm whitespace-normal text-left"
                                    >{roleBadge}</span>
                            {/if}
                        </td>
                    {/if}
                    <td class="whitespace-nowrap">{labelForItem(item)}</td>
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
