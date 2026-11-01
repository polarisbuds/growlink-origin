<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { page } from "$app/state";
    import { api } from "$lib/api";
    import { showToast } from "$lib/toast";
    import MeetingAnswerSidePanel from "$lib/components/MeetingAnswerSidePanel.svelte";
    import {
        dualTopicRows,
        isMeetingDone,
        menteeTargetLabel,
        ownAnswerDrafts,
        type DualTopicRow,
        type MeetingPayload,
        type MeetingUserRef,
    } from "$lib/meeting-helpers";

    let payload = $state<MeetingPayload | null>(null);
    let drafts = $state<Record<number, string>>({});
    let panelOpen = $state(true);
    let sharers = $state<MeetingUserRef[]>([]);

    let uuid = $derived(page.params.uuid);
    let rows = $derived(payload ? dualTopicRows(payload) : []);
    let viewerRole = $derived(payload?.viewerRole ?? "viewer");
    let canWrite = $derived(!!payload?.canWriteAnswers);
    let done = $derived(isMeetingDone(payload?.meeting?.status));
    let showBothReadOnly = $derived(
        done && (viewerRole === "mentor" || viewerRole === "mentee"),
    );
    let showMenteeContext = $derived(
        viewerRole === "mentor" && !showBothReadOnly,
    );

    function rowKey(row: DualTopicRow) {
        return row.topic?.id ?? row.displayNo;
    }

    async function load() {
        payload = await api<MeetingPayload>("GET", `/api/v1/meetings/${uuid}`);
        sharers = payload.sharers || [];
        drafts = ownAnswerDrafts(payload, payload.viewerRole);
    }

    onMount(load);

    function answerPayload() {
        if (!payload) return [];
        const meetingId = payload.meeting.id ?? 0;
        return rows
            .map((row) => {
                const topicId = row.topic?.id;
                if (!topicId) return null;
                const existing =
                    viewerRole === "mentor"
                        ? row.mentorAnswer
                        : row.menteeAnswer;
                return {
                    id: existing?.id ?? 0,
                    meetingId,
                    topicId,
                    content: drafts[topicId] ?? "",
                };
            })
            .filter(Boolean);
    }

    async function saveDraft() {
        await api("POST", `/api/v1/meetings/${uuid}/draft`, {
            answers: answerPayload(),
        });
        showToast("下書きを保存しました");
        await load();
    }
    async function submit() {
        await api("POST", `/api/v1/meetings/${uuid}/submit`, {
            answers: answerPayload(),
        });
        showToast("提出しました");
        await goto("/");
    }
    async function cancelSubmit() {
        await api("POST", `/api/v1/meetings/${uuid}/draft/cancel`, null);
        showToast("提出を取り消しました");
        await goto("/");
    }
</script>

<svelte:head><title>面談記入 - Growlink</title></svelte:head>

{#if payload}
    <div class="flex max-w-5xl items-center justify-between gap-4">
        <div>
            <h1 class="font-display text-2xl font-semibold">面談記入</h1>
            <p class="mt-1 text-sm text-base-content/60">
                {#if showBothReadOnly}
                    提出済みの内容を確認できます。共有者のみ変更できます。
                {:else if canWrite}
                    各質問への回答を入力し、下書き保存または提出できます
                {:else}
                    この面談への記入権限がありません
                {/if}
            </p>
        </div>
        <button
            type="button"
            class="btn btn-ghost btn-outline btn-sm shrink-0"
            onclick={() => goto("/")}>戻る</button>
    </div>

    {#if canWrite}
        <div class="mt-6 flex max-w-5xl flex-wrap gap-2 justify-end">
            <button type="button" class="btn btn-sm btn-outline" onclick={saveDraft}
                >下書き保存</button>
            <button type="button" class="btn btn-sm btn-primary" onclick={submit}
                >提出</button>
            <button type="button" class="btn btn-sm btn-ghost" onclick={cancelSubmit}
                >提出キャンセル</button>
        </div>
    {/if}

    <div class="mt-6 flex max-w-5xl flex-col gap-6 lg:flex-row lg:items-start">
        <div class="min-w-0 flex-1 space-y-4">
            {#each rows as row, i (rowKey(row))}
                {@const topicId = row.topic?.id}
                <div
                    class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm">
                    <div class="flex items-start gap-2">
                        <span class="badge badge-primary badge-outline shrink-0"
                            >Q{row.displayNo || i + 1}</span>
                        <div class="min-w-0 flex-1">
                            <p class="font-medium leading-relaxed">
                                {row.topic?.question}
                            </p>
                            {#if row.topic?.description}
                                <p
                                    class="mt-1 whitespace-pre-wrap text-sm text-base-content/60">
                                    {row.topic.description}
                                </p>
                            {/if}
                        </div>
                    </div>

                    {#if showBothReadOnly}
                        <div class="mt-4 space-y-3">
                            <div>
                                <p class="text-xs font-semibold text-base-content/60">メンティー</p>
                                {#if row.menteeAnswer?.content}
                                    <p class="mt-1 whitespace-pre-wrap text-sm">{row.menteeAnswer.content}</p>
                                {:else}
                                    <p class="mt-1 text-sm italic text-base-content/40">（未回答）</p>
                                {/if}
                            </div>
                            <div>
                                <p class="text-xs font-semibold text-base-content/60">メンター</p>
                                {#if row.mentorAnswer?.content}
                                    <p class="mt-1 whitespace-pre-wrap text-sm">{row.mentorAnswer.content}</p>
                                {:else}
                                    <p class="mt-1 text-sm italic text-base-content/40">（未回答）</p>
                                {/if}
                            </div>
                        </div>
                    {:else if showMenteeContext}
                        <div class="mt-4 space-y-3">
                            <div class="rounded-2xl bg-base-200/40 p-3">
                                <p class="text-xs font-semibold text-base-content/60">メンティー</p>
                                {#if row.menteeAnswer?.content}
                                    <p class="mt-1 whitespace-pre-wrap text-sm">{row.menteeAnswer.content}</p>
                                {:else}
                                    <p class="mt-1 text-sm italic text-base-content/40">（未回答）</p>
                                {/if}
                            </div>
                            {#if topicId != null && canWrite}
                                <div>
                                    <p class="text-xs font-semibold text-base-content/60">メンター（あなた）</p>
                                    <textarea
                                        class="textarea textarea-bordered mt-1 w-full"
                                        rows="3"
                                        value={drafts[topicId] ?? ""}
                                        oninput={(e) => {
                                            drafts[topicId] = e.currentTarget.value;
                                        }}></textarea>
                                </div>
                            {/if}
                        </div>
                    {:else if topicId != null && canWrite && viewerRole === "mentee"}
                        <textarea
                            class="textarea textarea-bordered mt-3 w-full"
                            rows="3"
                            value={drafts[topicId] ?? ""}
                            oninput={(e) => {
                                drafts[topicId] = e.currentTarget.value;
                            }}></textarea>
                    {:else if !canWrite}
                        <p class="mt-3 text-sm text-base-content/50">記入できません</p>
                    {/if}
                </div>
            {:else}
                <p class="text-sm text-base-content/60">表示できる質問がありません</p>
            {/each}
        </div>

        <MeetingAnswerSidePanel
            bind:open={panelOpen}
            {uuid}
            implementedAt={payload.meeting?.implementedAt}
            menteeLabel={menteeTargetLabel(payload.meeting?.pairing)}
            bind:sharers
            canEditSharers={!!payload.canEditSharers} />
    </div>
{:else}
    <p class="text-sm text-base-content/60">読み込み中…</p>
{/if}
