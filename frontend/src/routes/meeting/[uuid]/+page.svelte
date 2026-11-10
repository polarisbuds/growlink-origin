<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/state";
    import { api } from "$lib/api";
    import { showToast } from "$lib/toast";
    import MeetingAnswerSidePanel from "$lib/components/MeetingAnswerSidePanel.svelte";
    import {
        formatDate,
        dualTopicRows,
        shouldHideMeetingContentBeforeSubmit,
        isMeetingDone,
        isMeetingSubmitted,
        menteeTargetLabel,
        type MeetingPayload,
        type DualTopicRow,
        type MeetingUserRef,
    } from "$lib/meeting-helpers";
    import { statusBadge, statusLabel } from "$lib/status";
    import { goto } from "$app/navigation";

    let payload = $state<MeetingPayload | null>(null);
    let comment = $state("");
    let panelOpen = $state(true);
    let sharers = $state<MeetingUserRef[]>([]);
    let implementedDate = $state("");

    let rows = $derived(payload ? dualTopicRows(payload) : []);
    let meeting = $derived(payload?.meeting);
    let template = $derived(meeting?.template);
    let viewerRole = $derived(payload?.viewerRole);
    let submitted = $derived(isMeetingSubmitted(meeting?.status));
    let done = $derived(isMeetingDone(meeting?.status));
    let hideContentBeforeSubmit = $derived(
        shouldHideMeetingContentBeforeSubmit(viewerRole, meeting?.status),
    );
    let showCancelSubmit = $derived(submitted && viewerRole === "mentee");
    let canEditImplementedDate = $derived(
        !!payload?.canEditSharers || !!payload?.canEdit,
    );
    let uuid = $derived(page.params.uuid);

    function rowKey(row: DualTopicRow) {
        return row.topic?.uuid ?? String(row.displayNo);
    }

    async function loadPayload() {
        payload = await api<MeetingPayload>("GET", `/api/v1/meetings/${uuid}`);
        sharers = payload.sharers || [];
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

    onMount(loadPayload);

    async function postComment(isLimited: boolean) {
        if (!comment.trim() || !payload) return;
        await api("POST", `/api/v1/meetings/${uuid}/meeting-comment`, {
            content: comment,
            isLimited,
        });
        comment = "";
        showToast("投稿しました");
        await loadPayload();
    }

    async function cancelSubmit() {
        await api("POST", `/api/v1/meetings/${uuid}/draft/cancel`, null);
        showToast("提出を取り消しました");
        await goto("/");
    }

    function formatCommentDate(iso?: string) {
        if (!iso) return "";
        const d = new Date(iso);
        if (Number.isNaN(d.getTime())) return "";
        return d.toLocaleString("ja-JP");
    }
</script>

<svelte:head><title>面談閲覧 - Growlink</title></svelte:head>

{#if payload && meeting}
    <div
        class="flex max-w-7xl min-w-0 w-full flex-wrap items-start justify-between gap-4">
        <div class="min-w-0">
            <h1 class="truncate font-display text-2xl font-semibold">
                {template?.title || "面談"}
            </h1>
            <div
                class="mt-2 flex flex-wrap items-center gap-2 text-sm text-base-content/70">
                <span>実施日: {formatDate(meeting.implementedAt)}</span>
                <span class="badge badge-sm {statusBadge(meeting.status ?? 0)}"
                    >{statusLabel(meeting.status ?? 0)}</span>
            </div>
            <p class="mt-1 text-sm text-base-content/60">
                {#if hideContentBeforeSubmit}
                    メンティーの提出後に閲覧できます
                {:else if done}
                    提出済みの回答とコメントを確認できます
                {:else}
                    提出済みの回答を確認できます
                {/if}
            </p>
        </div>
        <div class="flex shrink-0 flex-wrap gap-2">
            {#if showCancelSubmit}
                <button
                    type="button"
                    class="btn btn-ghost btn-sm"
                    onclick={cancelSubmit}>提出キャンセル</button>
            {/if}
            <button
                onclick={() => {
                    goto("/");
                }}
                class="btn btn-ghost btn-sm btn-outline">戻る</button>
        </div>
    </div>
    {#if hideContentBeforeSubmit}
        <div
            class="mt-6 max-w-2xl rounded-3xl border border-dashed border-base-300 bg-base-100 p-10 text-center text-sm text-base-content/60">
            提出後に閲覧できます
        </div>
    {:else}
        <div class="mt-6 flex max-w-7xl flex-col gap-6 lg:flex-row lg:items-start">
            <div class="min-w-0 flex-1 space-y-4">
                {#each rows as row (rowKey(row))}
                    <div
                        class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm">
                        <div class="flex items-start gap-2">
                            <span
                                class="badge badge-primary badge-outline shrink-0 max-w-12"
                                >Q{row.displayNo}</span>
                            <div class="min-w-0 flex-1">
                                <p class="font-medium leading-relaxed">
                                    {row.topic?.question || "質問"}
                                </p>
                                {#if row.topic?.description}
                                    <p
                                        class="mt-1 whitespace-pre-wrap text-sm text-base-content/60">
                                        {row.topic.description}
                                    </p>
                                {/if}
                            </div>
                        </div>
                        <div class="mt-4 flex flex-row gap-4 ml-12">
                            <div class="flex-1">
                                <p
                                    class="text-xs font-semibold text-base-content/60">
                                    メンティー
                                </p>
                                {#if row.menteeAnswer?.content}
                                    <p
                                        class="mt-1 whitespace-pre-wrap text-sm text-base-content/80">
                                        {row.menteeAnswer.content}
                                    </p>
                                {:else}
                                    <p
                                        class="mt-1 text-sm italic text-base-content/40">
                                        （未回答）
                                    </p>
                                {/if}
                            </div>
                            <div class="flex-1">
                                <p
                                    class="text-xs font-semibold text-base-content/60">
                                    メンター
                                </p>
                                {#if row.mentorAnswer?.content}
                                    <p
                                        class="mt-1 whitespace-pre-wrap text-sm text-base-content/80">
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
                    </div>
                {/each}

                {#if done}
                    <div
                        class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm">
                        <h2 class="font-semibold">ミーティングコメント</h2>
                        {#if (payload.meetingComments || []).length > 0}
                            <ul class="mt-4 space-y-3">
                                {#each payload.meetingComments || [] as c (c.uuid)}
                                    <li
                                        class="rounded-2xl border border-base-300 bg-base-200/30 p-4">
                                        <div
                                            class="flex flex-wrap items-center gap-2 text-xs text-base-content/60">
                                            <span
                                                class="font-medium text-base-content"
                                                >{c.user?.name ||
                                                    "ユーザ"}</span>
                                            <span
                                                >{formatCommentDate(
                                                    c.createdAt,
                                                )}</span>
                                            {#if c.isLimited}
                                                <span
                                                    class="badge badge-warning badge-outline badge-xs"
                                                    >限定公開</span>
                                            {/if}
                                        </div>
                                        <p
                                            class="mt-2 whitespace-pre-wrap text-sm">
                                            {c.content}
                                        </p>
                                    </li>
                                {/each}
                            </ul>
                        {:else}
                            <p class="mt-3 text-sm text-base-content/50">
                                コメントはまだありません
                            </p>
                        {/if}
                        <p class="mt-4 text-xs text-base-content/60">
                            ペア以外への共有や限定公開が選べます
                        </p>
                        <textarea
                            class="textarea textarea-bordered mt-3 w-full"
                            rows="2"
                            bind:value={comment}
                            placeholder="コメント"></textarea>
                        <div class="mt-3 flex flex-col gap-2">
                            <div
                                class="flex flex-1 flex-wrap gap-2 justify-end">
                                <button
                                    type="button"
                                    class="btn btn-sm btn-primary"
                                    onclick={() => postComment(false)}
                                    >コメント投稿</button>
                                <button
                                    type="button"
                                    class="btn btn-sm btn-outline"
                                    onclick={() => postComment(true)}
                                    >限定公開コメント</button>
                            </div>
                            <div class="flex flex-1 justify-end">
                                <p class="text-xs text-base-content/50">
                                    限定公開はペア（メンター・メンティ）のみが閲覧できます
                                </p>
                            </div>
                        </div>
                    </div>
                {/if}
            </div>

            <MeetingAnswerSidePanel
                bind:open={panelOpen}
                {uuid}
                bind:implementedDate
                menteeLabel={menteeTargetLabel(meeting?.pairing)}
                bind:sharers
                canEditSharers={!!payload.canEditSharers}
                {canEditImplementedDate} />
        </div>
    {/if}
{:else}
    <p class="text-sm text-base-content/60">読み込み中…</p>
{/if}
