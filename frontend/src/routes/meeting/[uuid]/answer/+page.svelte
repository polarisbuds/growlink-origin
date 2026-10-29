<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { page } from "$app/state";
    import { api } from "$lib/api";
    import { showToast } from "$lib/toast";
    import {
        displayTopics,
        type MeetingPayload,
        type DisplayTopicRow,
    } from "$lib/meeting-helpers";

    let payload = $state<MeetingPayload | null>(null);
    let drafts = $state<Record<number, string>>({});

    let rows = $derived(payload ? displayTopics(payload) : []);
    let uuid = $derived(page.params.uuid);

    function rowKey(row: DisplayTopicRow) {
        return row.topic?.id ?? row.displayNo;
    }

    function draftKey(row: DisplayTopicRow): number | null {
        const tid = row.topic?.id;
        return tid != null ? tid : null;
    }

    onMount(async () => {
        payload = await api<MeetingPayload>("GET", `/api/v1/meetings/${uuid}`);
        const next: Record<number, string> = {};
        for (const row of displayTopics(payload)) {
            const key = draftKey(row);
            if (key == null) continue;
            next[key] = row.answer?.content || "";
        }
        drafts = next;
        console.log(rows);
    });

    function answerPayload() {
        if (!payload) return [];
        const meetingId = payload.meeting.id ?? 0;
        return rows
            .map((row) => {
                const topicId = row.topic?.id;
                if (!topicId) return null;
                return {
                    id: row.answer?.id ?? 0,
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
    }
    async function submit() {
        await api("POST", `/api/v1/meetings/${uuid}/submit`, {
            answers: answerPayload(),
        });
        showToast("提出しました");
        await goto("/meeting");
    }
    async function cancelSubmit() {
        await api("POST", `/api/v1/meetings/${uuid}/draft/cancel`, null);
        showToast("提出を取り消しました");
    }
</script>

<svelte:head><title>面談記入 - Growlink</title></svelte:head>

<div class="flex max-w-2xl items-center justify-between gap-4">
    <div>
        <h1 class="font-display text-2xl font-semibold">面談記入</h1>
        <p class="mt-1 text-sm text-base-content/60">
            各質問への回答を入力し、下書き保存または提出できます
        </p>
    </div>
    <a href="/meeting" class="btn btn-ghost btn-sm shrink-0">戻る</a>
</div>

<div class="mt-6 flex max-w-2xl flex-wrap gap-2 justify-end">
    <button type="button" class="btn btn-outline" onclick={saveDraft}
        >下書き保存</button>
    <button type="button" class="btn btn-primary" onclick={submit}>提出</button>
    <button type="button" class="btn btn-ghost" onclick={cancelSubmit}
        >提出キャンセル</button>
</div>

<div class="mt-6 max-w-2xl space-y-4">
    {#each rows as row, i (rowKey(row))}
        {@const key = draftKey(row)}
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
            {#if key != null}
                <textarea
                    class="textarea textarea-bordered mt-3 w-full"
                    rows="3"
                    value={drafts[key] ?? ""}
                    oninput={(e) => {
                        drafts[key] = e.currentTarget.value;
                    }}></textarea>
            {/if}
        </div>
    {:else}
        <p class="text-sm text-base-content/60">表示できる質問がありません</p>
    {/each}
</div>
