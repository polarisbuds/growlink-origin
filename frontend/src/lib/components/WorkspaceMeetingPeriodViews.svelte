<script lang="ts">
    import { statusBadge, statusLabel } from "$lib/status";
    import {
        workspace,
        type WorkspaceTab,
    } from "$lib/stores/workspaceTab.svelte";
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
    type PeriodKey = Exclude<WorkspaceTab, "all">;

    let {
        meetings,
        viewMode,
        currentUserUuid,
        roleBadgeMode = "none",
        labelForItem = pairLabel,
    }: {
        meetings: MeetingListItem[];
        viewMode: "card" | "table";
        currentUserUuid?: string;
        roleBadgeMode?: RoleBadgeMode;
        labelForItem?: (item: MeetingListItem) => string;
    } = $props();

    let periodTab = $derived(workspace.workspaceTab);

    let periods = $derived(splitByMeetingPeriod(meetings));

    const sections: {
        key: PeriodKey;
        tabLabel: string;
        sectionTitle: string;
        empty: string;
    }[] = [
        {
            key: "thisMonth",
            tabLabel: "今月",
            sectionTitle: "今月の面談",
            empty: "今月の面談はありません",
        },
        {
            key: "past",
            tabLabel: "これまで",
            sectionTitle: "これまでの面談",
            empty: "面談の履歴はありません",
        },
        {
            key: "planned",
            tabLabel: "これから",
            sectionTitle: "計画済みの面談",
            empty: "計画済みの面談はありません",
        },
    ];

    let activeSection = $derived(
        periodTab === "all"
            ? null
            : (sections.find((s) => s.key === periodTab) ?? sections[0]),
    );
    let singleItems = $derived(
        periodTab === "all" ? [] : periods[periodTab as PeriodKey],
    );
</script>

<div role="tablist" class="tabs tabs-bordered mt-8 max-w-3xl">
    <button
        type="button"
        role="tab"
        class="tab"
        class:tab-active={periodTab === "all"}
        onclick={() => workspace.changeTab("all")}>
        すべて ({meetings.length})
    </button>
    {#each sections as section (section.key)}
        <button
            type="button"
            role="tab"
            class="tab"
            class:tab-active={periodTab === section.key}
            onclick={() => workspace.changeTab(section.key)}>
            {section.tabLabel} ({periods[section.key].length})
        </button>
    {/each}
</div>

{#snippet meetingCards(items: MeetingListItem[], empty: string)}
    <div class="mt-4 grid gap-4 md:grid-cols-2 xl:grid-cols-3">
        {#each items as item (item.uuid)}
            {@const roleBadge = meetingRoleBadgeLabel(
                item,
                currentUserUuid,
                roleBadgeMode,
            )}
            <div
                class="group rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm hover:border-primary/40 transition-colors">
                <div class="flex items-start justify-between gap-3">
                    <div class="min-w-0">
                        {#if roleBadge}
                            <span class="badge badge-outline badge-sm mb-2"
                                >{roleBadge}</span>
                        {/if}
                        <p class="text-sm text-base-content/60 truncate">
                            {labelForItem(item)}
                        </p>
                        <h3 class="mt-1 font-semibold">
                            {item.template?.title || "1on1"}
                        </h3>
                    </div>
                    <span
                        class="badge shrink-0 {statusBadge(item.status ?? 0)}">
                        {statusLabel(item.status ?? 0)}
                    </span>
                </div>
                <p class="mt-3 text-sm text-base-content/70">
                    実施予定日: {formatDate(item.implementedAt)}
                </p>
                <div class="mt-4 flex flex-wrap gap-2 items-center justify-end">
                    {#if canUserAnswerMeeting(item, currentUserUuid)}
                        <a
                            href="/meeting/{item.uuid}/answer"
                            class="btn btn-sm btn-primary">記入</a>
                    {:else if canUserViewMeeting(item, currentUserUuid)}
                        <a
                            href="/meeting/{item.uuid}"
                            class="btn btn-sm btn-ghost btn-outline">閲覧</a>
                    {/if}
                </div>
            </div>
        {/each}
        {#if items.length === 0}
            <div
                class="rounded-3xl border border-dashed border-base-300 bg-base-100 p-10 text-center text-sm text-base-content/50 md:col-span-2 xl:col-span-3">
                {empty}
            </div>
        {/if}
    </div>
{/snippet}

{#if periodTab === "all"}
    {#each sections as section, i (section.key)}
        {@const items = periods[section.key]}
        <section class={i === 0 ? "mt-8" : "mt-10"}>
            <h2 class="text-sm font-semibold text-base-content/70">
                {section.sectionTitle}
            </h2>
            {#if viewMode === "table"}
                <div class="mt-4">
                    <WorkspaceMeetingTable
                        {items}
                        {currentUserUuid}
                        {roleBadgeMode}
                        {labelForItem}
                        emptyMessage={section.empty} />
                </div>
            {:else}
                {@render meetingCards(items, section.empty)}
            {/if}
        </section>
    {/each}
{:else if viewMode === "table"}
    <div class="mt-4">
        <WorkspaceMeetingTable
            items={singleItems}
            {currentUserUuid}
            {roleBadgeMode}
            {labelForItem}
            emptyMessage={activeSection?.empty ?? "面談はありません"} />
    </div>
{:else}
    {@render meetingCards(
        singleItems,
        activeSection?.empty ?? "面談はありません",
    )}
{/if}
