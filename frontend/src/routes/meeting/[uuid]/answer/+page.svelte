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
        isMeetingSubmitted,
        menteeTargetLabel,
        ownAnswerDrafts,
        shouldHideMeetingContentBeforeSubmit,
        type DualTopicRow,
        type MeetingPayload,
        type MeetingUserRef,
    } from "$lib/meeting-helpers";

    let payload = $state<MeetingPayload | null>(null);
    let drafts = $state<Record<string, string>>({});
    let panelOpen = $state(true);
    let sharers = $state<MeetingUserRef[]>([]);

    let uuid = $derived(page.params.uuid);
    let rows = $derived(payload ? dualTopicRows(payload) : []);
    let viewerRole = $derived(payload?.viewerRole ?? "viewer");
    let canWriteMentee = $derived(!!payload?.canWriteAnswers);
    let canWriteMentor = $derived(!!payload?.canWriteTopicComments);
    let canWrite = $derived(canWriteMentee || canWriteMentor);
    let submitted = $derived(isMeetingSubmitted(payload?.meeting?.status));
    let done = $derived(isMeetingDone(payload?.meeting?.status));
    let showBothReadOnly = $derived(
        (submitted || done) && viewerRole === "mentee",
    );
    let showMentorAnswerForm = $derived(
        viewerRole === "mentor" && submitted && canWriteMentor,
    );
    let showMentorWaiting = $derived(
        viewerRole === "mentor" && !submitted && !done,
    );
    let showCancelSubmit = $derived(submitted && viewerRole === "mentee");
    let hideContentBeforeSubmit = $derived(
        shouldHideMeetingContentBeforeSubmit(
            viewerRole,
            payload?.meeting?.status,
        ),
    );
    let implementedDate = $state("");
    let canEditImplementedDate = $derived(
        !!payload?.canEditSharers || !!payload?.canEdit,
    );
    let templateTitle = $derived(
        payload?.meeting?.template?.title || "面談記入",
    );
    let templateDescription = $derived(
        payload?.meeting?.template?.description?.trim() || "",
    );

    function rowKey(row: DualTopicRow) {
        return row.topic?.uuid ?? String(row.displayNo);
    }

    async function load() {
        payload = await api<MeetingPayload>("GET", `/api/v1/meetings/${uuid}`);
        sharers = payload.sharers || [];
        drafts = ownAnswerDrafts(payload, payload.viewerRole);
        const raw = payload.meeting?.implementedAt;
        if (raw) {
            const d = new Date(raw);
            if (!Number.isNaN(d.getTime())) {
                implementedDate = d.toISOString().slice(0, 10);
            }
        } else {
            implementedDate = "";
        }
    }

    onMount(load);

    function answerPayload() {
        if (!payload) return [];
        return rows
            .map((row) => {
                const topicUuid = row.topic?.uuid;
                if (!topicUuid) return null;
                return {
                    uuid: row.menteeAnswer?.uuid ?? "",
                    content: drafts[topicUuid] ?? "",
                    topic: { uuid: topicUuid },
                };
            })
            .filter(Boolean);
    }

    function mentorCommentPayload() {
        return rows
            .map((row) => {
                const topicUuid = row.topic?.uuid;
                if (!topicUuid) return null;
                const answerUuid = row.menteeAnswer?.uuid;
                if (!answerUuid) return null;
                return {
                    answerUuid,
                    content: drafts[topicUuid] ?? "",
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
    async function saveMentorComplete() {
        await api("POST", `/api/v1/meetings/${uuid}/complete`, {
            comments: mentorCommentPayload(),
        });
        showToast("保存し完了しました");
        await goto("/");
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
    <div class="flex max-w-7xl items-center justify-between gap-4">
        <div>
            <h1 class="font-display text-2xl font-semibold">{templateTitle}</h1>
            {#if templateDescription}
                <p class="mt-1 text-xs text-base-content/50">
                    {templateDescription}
                </p>
            {/if}
            <p class="mt-1 text-sm text-base-content/60">
                {#if showBothReadOnly}
                    提出済みの内容を確認できます。
                {:else if showMentorAnswerForm}
                    メンティーの回答を確認し、メンターとして記入できます
                {:else if showMentorWaiting || hideContentBeforeSubmit}
                    メンティーの提出後に記入できます
                {:else if canWrite}
                    各質問への回答を入力し、下書き保存または提出できます
                {:else}
                    この面談への記入権限がありません
                {/if}
            </p>
        </div>
        <div class="flex justify-end items-center gap-2 max-w-7xl flex-wrap">
            {#if (canWrite || showCancelSubmit) && !hideContentBeforeSubmit}
                {#if canWrite}
                    {#if viewerRole === "mentor"}
                        <button
                            type="button"
                            class="btn btn-sm btn-primary"
                            onclick={saveMentorComplete}>保存し完了する</button>
                    {:else}
                        <button
                            type="button"
                            class="btn btn-sm btn-outline"
                            onclick={saveDraft}>下書き保存</button>
                        <button
                            type="button"
                            class="btn btn-sm btn-primary"
                            onclick={submit}>提出</button>
                    {/if}
                {/if}
                {#if showCancelSubmit}
                    <button
                        type="button"
                        class="btn btn-sm btn-ghost"
                        onclick={cancelSubmit}>提出キャンセル</button>
                {/if}
            {/if}
            <button
                type="button"
                class="btn btn-ghost btn-outline btn-sm shrink-0"
                onclick={() => goto("/")}>戻る</button>
        </div>
    </div>

    <div class="mt-6 flex max-w-7xl flex-col gap-6 lg:flex-row lg:items-start">
        <div class="min-w-0 flex-1 space-y-4">
            {#if hideContentBeforeSubmit}
                <div
                    class="rounded-3xl border border-dashed border-base-300 bg-base-100 p-10 text-center text-sm text-base-content/60">
                    提出後に閲覧できます
                </div>
            {:else}
                {#each rows as row, i (rowKey(row))}
                    {@const topicUuid = row.topic?.uuid}
                    <div
                        class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm">
                        <div class="flex items-start gap-2">
                            <span
                                class="badge badge-primary badge-outline shrink-0"
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
                                    <p
                                        class="text-xs font-semibold text-base-content/60">
                                        メンティー
                                    </p>
                                    {#if row.menteeAnswer?.content}
                                        <p
                                            class="mt-1 whitespace-pre-wrap text-sm">
                                            {row.menteeAnswer.content}
                                        </p>
                                    {:else}
                                        <p
                                            class="mt-1 text-sm italic text-base-content/40">
                                            （未回答）
                                        </p>
                                    {/if}
                                </div>
                                <div>
                                    <p
                                        class="text-xs font-semibold text-base-content/60">
                                        メンター
                                    </p>
                                    {#if row.mentorAnswer?.content}
                                        <p
                                            class="mt-1 whitespace-pre-wrap text-sm">
                                            {row.mentorAnswer.content}
                                        </p>
                                    {:else}
                                        <p
                                            class="mt-1 text-sm italic text-base-content/40">
                                            （未回答）
                                        </p>
                                    {/if}
                                </div>
                            </div>
                        {:else if showMentorAnswerForm}
                            <div class="mt-4 space-y-3">
                                <div class="rounded-2xl bg-base-200/40 p-3">
                                    <p
                                        class="text-xs font-semibold text-base-content/60">
                                        メンティー
                                    </p>
                                    {#if row.menteeAnswer?.content}
                                        <p
                                            class="mt-1 whitespace-pre-wrap text-sm">
                                            {row.menteeAnswer.content}
                                        </p>
                                    {:else}
                                        <p
                                            class="mt-1 text-sm italic text-base-content/40">
                                            （未回答）
                                        </p>
                                    {/if}
                                </div>
                                {#if topicUuid}
                                    <div>
                                        <p
                                            class="text-xs font-semibold text-base-content/60">
                                            メンター（あなた）
                                        </p>
                                        <textarea
                                            class="textarea textarea-bordered mt-1 w-full"
                                            rows="3"
                                            value={drafts[topicUuid] ?? ""}
                                            oninput={(e) => {
                                                drafts[topicUuid] =
                                                    e.currentTarget.value;
                                            }}></textarea>
                                    </div>
                                {/if}
                            </div>
                        {:else if showMentorWaiting}
                            <p class="mt-3 text-sm text-base-content/50">
                                メンティーの提出後に記入できます
                            </p>
                        {:else if topicUuid && canWrite && viewerRole === "mentee"}
                            <textarea
                                class="textarea textarea-bordered mt-3 w-full"
                                rows="3"
                                value={drafts[topicUuid] ?? ""}
                                oninput={(e) => {
                                    drafts[topicUuid] = e.currentTarget.value;
                                }}></textarea>
                        {:else if !canWrite}
                            <p class="mt-3 text-sm text-base-content/50">
                                記入できません
                            </p>
                        {/if}
                    </div>
                {:else}
                    <p class="text-sm text-base-content/60">
                        表示できる質問がありません
                    </p>
                {/each}
            {/if}
        </div>

        <MeetingAnswerSidePanel
            bind:open={panelOpen}
            {uuid}
            bind:implementedDate
            menteeLabel={menteeTargetLabel(payload.meeting?.pairing)}
            bind:sharers
            canEditSharers={!!payload.canEditSharers}
            {canEditImplementedDate} />
    </div>
{:else}
    <p class="text-sm text-base-content/60">読み込み中…</p>
{/if}
