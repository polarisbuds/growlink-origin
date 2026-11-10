<script lang="ts">
    import { onMount } from "svelte";
    import { api } from "$lib/api";
    import {
        splitInvolvedMeetingsForWorkspace,
        mentorTabPairLabel,
        type MeetingListItem,
    } from "$lib/meeting-helpers";
    import WorkspaceMeetingPeriodViews from "$lib/components/WorkspaceMeetingPeriodViews.svelte";
    import LayoutGrid from "lucide-svelte/icons/layout-grid";
    import Table from "lucide-svelte/icons/table";

    let tab = $state<"own" | "mentor" | "shared">("own");
    let viewMode = $state<"card" | "table">("card");
    let ownMeetings = $state<MeetingListItem[]>([]);
    let mentorMeetings = $state<MeetingListItem[]>([]);
    let showMentorTab = $state(false);
    let sharedMeetings = $state<MeetingListItem[]>([]);
    let currentUserUuid = $state("");

    onMount(async () => {
        const [own, shared, profile] = await Promise.all([
            api<MeetingListItem[]>("GET", "/api/v1/meetings/involved"),
            api<MeetingListItem[]>("GET", "/api/v1/meetings/shared"),
            api<{ email?: string }>("GET", "/api/v1/setting/profile"),
        ]);
        const directory = await api<{ users?: { uuid: string; email: string }[] }>(
            "GET",
            "/api/v1/planning/directory",
        );
        currentUserUuid =
            directory.users?.find((u) => u.email === profile.email)?.uuid || "";
        const split = splitInvolvedMeetingsForWorkspace(own, currentUserUuid);
        ownMeetings = split.own;
        mentorMeetings = split.asMentor;
        showMentorTab = split.showMentorTab;
        sharedMeetings = shared;
        if (tab === "mentor" && !showMentorTab) {
            tab = "own";
        }
    });
</script>

<svelte:head><title>ワークスペース - Growlink</title></svelte:head>

<div class="flex flex-wrap items-end justify-between gap-4">
    <div>
        <h1 class="font-display text-2xl font-semibold">ワークスペース</h1>
        <p class="mt-1 text-sm text-base-content/60">
            あなたの1on1と共有された1on1
        </p>
    </div>
    <div class="join">
        <button
            type="button"
            class="btn btn-sm join-item {viewMode === 'card'
                ? 'btn-primary'
                : 'btn-ghost'}"
            onclick={() => (viewMode = "card")}
            aria-label="カード表示">
            <LayoutGrid class="h-4 w-4" />
        </button>
        <button
            type="button"
            class="btn btn-sm join-item {viewMode === 'table'
                ? 'btn-primary'
                : 'btn-ghost'}"
            onclick={() => (viewMode = "table")}
            aria-label="テーブル表示">
            <Table class="h-4 w-4" />
        </button>
    </div>
</div>

<div role="tablist" class="tabs tabs-bordered mt-6 max-w-3xl">
    <button
        type="button"
        role="tab"
        class="tab"
        class:tab-active={tab === "own"}
        onclick={() => (tab = "own")}>
        あなたの1on1
    </button>
    {#if showMentorTab}
        <button
            type="button"
            role="tab"
            class="tab"
            class:tab-active={tab === "mentor"}
            onclick={() => (tab = "mentor")}>
            メンティーとの面談
        </button>
    {/if}
    <button
        type="button"
        role="tab"
        class="tab"
        class:tab-active={tab === "shared"}
        onclick={() => (tab = "shared")}>
        共有された1on1
    </button>
</div>

{#if tab === "own"}
    <WorkspaceMeetingPeriodViews
        meetings={ownMeetings}
        {viewMode}
        {currentUserUuid}
        roleBadgeMode="mentee-only" />
{:else if tab === "mentor"}
    <WorkspaceMeetingPeriodViews
        meetings={mentorMeetings}
        {viewMode}
        {currentUserUuid}
        roleBadgeMode="mentor-always"
        labelForItem={mentorTabPairLabel} />
{:else}
    <WorkspaceMeetingPeriodViews
        meetings={sharedMeetings}
        {viewMode}
        {currentUserUuid}
        roleBadgeMode="none" />
{/if}
