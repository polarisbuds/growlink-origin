<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/state";
    import { api } from "$lib/api";
    import { showToast } from "$lib/toast";
    import {
        formatDate,
        displayTopics,
        type MeetingPayload,
        type DisplayTopicRow,
    } from "$lib/meeting-helpers";
    import { statusBadge, statusLabel } from "$lib/status";
    import { goto } from "$app/navigation";

    let payload = $state<MeetingPayload | null>(null);
    let comment = $state("");

    let rows = $derived(payload ? displayTopics(payload) : []);
    let meeting = $derived(payload?.meeting);
    let template = $derived(meeting?.template);
    let canEdit = $derived(!!payload?.canEdit);
    let uuid = $derived(page.params.uuid);

    function rowKey(row: DisplayTopicRow) {
        return row.topic?.id ?? row.displayNo;
    }

    async function loadPayload() {
        payload = await api<MeetingPayload>("GET", `/api/v1/meetings/${uuid}`);
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
        class="flex max-w-2xl min-w-0 w-full flex-wrap items-start justify-between gap-4">
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
                提出済みの回答とコメントを確認できます
            </p>
        </div>
        <div class="flex shrink-0 flex-wrap gap-2">
            {#if canEdit}
                <button
                    onclick={() => {
                        goto("/meeting/" + uuid + "/setting");
                    }}
                    class="btn btn-outline btn-sm">設定</button>
            {/if}
            <button
                onclick={() => {
                    goto("/meeting/" + uuid + "/answer");
                }}
                class="btn btn-outline btn-sm">更新</button>
            <button
                onclick={() => {
                    goto("/");
                }}
                class="btn btn-ghost btn-sm">戻る</button>
        </div>
    </div>
    <div class="mt-6 max-w-2xl space-y-4">
        {#each rows as row (rowKey(row))}
            <div
                class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm">
                <div class="flex items-start gap-2">
                    <span class="badge badge-primary badge-outline shrink-0"
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
                {#if row.answer?.content}
                    <p
                        class="mt-3 whitespace-pre-wrap text-sm text-base-content/80">
                        {row.answer.content}
                    </p>
                {:else}
                    <p class="mt-3 text-sm italic text-base-content/40">
                        （未回答）
                    </p>
                {/if}
            </div>
        {/each}
    </div>
    <div
        class="mt-6 max-w-2xl rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm">
        <h2 class="font-semibold">ミーティングコメント</h2>
        {#if (payload.meetingComments || []).length > 0}
            <ul class="mt-4 space-y-3">
                {#each payload.meetingComments || [] as c (c.id)}
                    <li
                        class="rounded-2xl border border-base-300 bg-base-200/30 p-4">
                        <div
                            class="flex flex-wrap items-center gap-2 text-xs text-base-content/60">
                            <span class="font-medium text-base-content"
                                >{c.user?.name || "ユーザ"}</span>
                            <span>{formatCommentDate(c.createdAt)}</span>
                            {#if c.isLimited}
                                <span
                                    class="badge badge-warning badge-outline badge-xs"
                                    >限定公開</span>
                            {/if}
                        </div>
                        <p class="mt-2 whitespace-pre-wrap text-sm">
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
        <div class="mt-2 flex flex-wrap gap-2">
            <button
                type="button"
                class="btn btn-sm btn-primary"
                onclick={() => postComment(false)}>コメント投稿</button>
            <button
                type="button"
                class="btn btn-sm btn-outline"
                onclick={() => postComment(true)}>限定公開コメント</button>
        </div>
        <p class="mt-2 text-xs text-base-content/50">
            限定公開はペア（メンター・メンティ）のみが閲覧できます
        </p>
    </div>
{:else}
    <p class="text-sm text-base-content/60">読み込み中…</p>
{/if}
