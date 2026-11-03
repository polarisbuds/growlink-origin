<script lang="ts">
    import { statusBadge, statusLabel } from "$lib/status";
    import {
        formatDate,
        pairLabel,
        splitByMeetingPeriod,
        canUserAnswerMeeting,
        canUserViewMeeting,
        meetingRoleBadgeLabel,
        type MeetingListItem,
    } from "$lib/meeting-helpers";
    import WorkspaceMeetingTable from "$lib/components/WorkspaceMeetingTable.svelte";

    type RoleBadgeMode = "mentee-only" | "mentor-always" | "none";

    let {
        meetings,
        viewMode,
        currentUserId,
        roleBadgeMode = "none",
        labelForItem = pairLabel,
    }: {
        meetings: MeetingListItem[];
        viewMode: "card" | "table";
        currentUserId?: number;
        roleBadgeMode?: RoleBadgeMode;
        labelForItem?: (item: MeetingListItem) => string;
    } = $props();

    let periods = $derived(splitByMeetingPeriod(meetings));

    const sections = [
        {
            key: "thisMonth" as const,
            title: "今月の面談",
            empty: "今月の面談はありません",
        },
        {
            key: "past" as const,
            title: "過去の面談",
            empty: "過去の面談はありません",
        },
        {
            key: "planned" as const,
            title: "計画された面談",
            empty: "計画された面談はありません",
        },
    ];
</script>

{#each sections as section (section.key)}
    {@const items = periods[section.key]}
    <section class={section.key === "thisMonth" ? "mt-8" : "mt-10"}>
        <h2 class="text-sm font-semibold text-base-content/70">
            {section.title}
        </h2>
        {#if viewMode === "table"}
            <div class="mt-4">
                <WorkspaceMeetingTable
                    {items}
                    {currentUserId}
                    {roleBadgeMode}
                    {labelForItem} />
            </div>
        {:else}
            <div class="mt-4 grid gap-4 md:grid-cols-2 xl:grid-cols-3">
                {#each items as item (item.uuid)}
                    {@const roleBadge = meetingRoleBadgeLabel(
                        item,
                        currentUserId,
                        roleBadgeMode,
                    )}
                    <div
                        class="group rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm hover:border-primary/40 transition-colors">
                        <div class="flex items-start justify-between gap-3">
                            <div class="min-w-0">
                                {#if roleBadge}
                                    <span
                                        class="badge badge-outline badge-sm mb-2"
                                        >{roleBadge}</span>
                                {/if}
                                <p
                                    class="text-sm text-base-content/60 truncate">
                                    {labelForItem(item)}
                                </p>
                                <h3 class="mt-1 font-semibold">
                                    {item.template?.title || "1on1"}
                                </h3>
                            </div>
                            <span
                                class="badge shrink-0 {statusBadge(
                                    item.status ?? 0,
                                )}">
                                {statusLabel(item.status ?? 0)}
                            </span>
                        </div>
                        <p class="mt-3 text-sm text-base-content/70">
                            実施予定日: {formatDate(item.implementedAt)}
                        </p>
                        <div
                            class="mt-4 flex flex-wrap gap-2 items-center justify-end">
                            {#if canUserAnswerMeeting(item, currentUserId)}
                                <a
                                    href="/meeting/{item.uuid}/answer"
                                    class="btn btn-sm btn-primary">記入</a>
                            {:else if canUserViewMeeting(item, currentUserId)}
                                <a
                                    href="/meeting/{item.uuid}"
                                    class="btn btn-sm btn-ghost btn-outline"
                                    >閲覧</a>
                            {/if}
                        </div>
                    </div>
                {/each}
                {#if items.length === 0}
                    <div
                        class="rounded-3xl border border-dashed border-base-300 bg-base-100 p-10 text-center text-sm text-base-content/50 md:col-span-2 xl:col-span-3">
                        {section.empty}
                    </div>
                {/if}
            </div>
        {/if}
    </section>
{/each}
