<script lang="ts">
    import { onMount } from "svelte";
    import { api } from "$lib/api";
    import {
        splitInvolvedMeetingsForWorkspace,
        type MeetingListItem,
    } from "$lib/meeting-helpers";
    import type { PlanningDirectoryUser } from "$lib/planning-directory";
    import WorkspaceMeetingPeriodViews from "$lib/components/WorkspaceMeetingPeriodViews.svelte";
    import LayoutGrid from "lucide-svelte/icons/layout-grid";
    import Table from "lucide-svelte/icons/table";

    let tab = $state<"own" | "shared">("own");
    let viewMode = $state<"card" | "table">("card");
    let ownMeetings = $state<MeetingListItem[]>([]);
    let mentorMeetings = $state<MeetingListItem[]>([]);
    let hiddenMentorMeeting = $state(false);
    let sharedMeetings = $state<MeetingListItem[]>([]);
    let currentUserUuid = $state("");
    let directoryUsers = $state<PlanningDirectoryUser[]>([]);

    let ownVisibleMeetings = $derived(
        hiddenMentorMeeting ? ownMeetings : [...ownMeetings, ...mentorMeetings],
    );

    let canIncludeMenteeMeetings = $derived(mentorMeetings.length > 0);

    onMount(async () => {
        const [own, shared, profile, directory] = await Promise.all([
            api<MeetingListItem[]>("GET", "/api/v1/meetings/involved"),
            api<MeetingListItem[]>("GET", "/api/v1/meetings/shared"),
            api<{ email?: string }>("GET", "/api/v1/setting/profile"),
            api<{ users?: PlanningDirectoryUser[] }>(
                "GET",
                "/api/v1/planning/directory",
            ),
        ]);
        directoryUsers = directory.users || [];
        currentUserUuid =
            directory.users?.find((u) => u.email === profile.email)?.uuid || "";
        const split = splitInvolvedMeetingsForWorkspace(own, currentUserUuid);
        ownMeetings = split.own;
        mentorMeetings = split.asMentor;
        sharedMeetings = shared;
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
        meetings={ownVisibleMeetings}
        {viewMode}
        {currentUserUuid}
        {directoryUsers}
        {canIncludeMenteeMeetings}
        bind:hiddenMentorMeeting />
{:else}
    <WorkspaceMeetingPeriodViews
        meetings={sharedMeetings}
        {viewMode}
        {currentUserUuid}
        {directoryUsers} />
{/if}
