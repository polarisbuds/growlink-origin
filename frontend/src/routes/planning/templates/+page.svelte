<script lang="ts">
    import { onMount } from "svelte";
    import { api } from "$lib/api";
    import { showToast } from "$lib/toast";
    import PlanningTabs from "$lib/components/PlanningTabs.svelte";
    import { validateTemplateForm } from "$lib/template-validation";
    import { ensureSession } from "$lib/session";

    type Topic = {
        id?: number;
        displayNo?: number;
        question?: string;
        description?: string;
        isRequire?: boolean;
    };
    type Template = {
        id: number;
        uuid?: string;
        displayNo?: number;
        title?: string;
        description?: string;
        creatorName?: string;
        creatorId?: number;
        isSample?: boolean;
        isLock?: boolean;
        topics?: Topic[];
    };

    let templates = $state<Template[]>([]);
    let currentUserId = $state<number | undefined>();
    let isSystemAdmin = $state(false);
    let modalOpen = $state(false);
    let editMode = $state(false);
    let deleteOpen = $state(false);
    let active = $state<Template | null>(null);
    let editForm = $state<{
        title: string;
        description: string;
        topics: Topic[];
    }>({
        title: "",
        description: "",
        topics: [],
    });

    async function persistOrder(list: Template[]) {
        const items = list.map((t, i) => ({ id: t.id, displayNo: i + 1 }));
        await api("PUT", "/api/v1/templates/reorder", { items });
    }

    async function reloadTemplates() {
        const list = await api<Template[]>("GET", "/api/v1/templates");
        templates = [...list].sort(
            (a, b) => (a.displayNo || 0) - (b.displayNo || 0),
        );
        if (active) {
            const next = templates.find((x) => x.id === active!.id);
            if (next) active = next;
        }
    }

    function canToggleLock(t: Template | null) {
        if (!t) return false;
        return (
            isSystemAdmin ||
            (currentUserId != null && t.creatorId === currentUserId)
        );
    }

    function canEditDelete(t: Template | null) {
        return !!t && !t.isLock;
    }

    onMount(async () => {
        const [list, profile, avatar] = await Promise.all([
            api<Template[]>("GET", "/api/v1/templates"),
            api<{ id?: number }>("GET", "/api/v1/setting/profile"),
            ensureSession(),
        ]);
        templates = [...list].sort(
            (a, b) => (a.displayNo || 0) - (b.displayNo || 0),
        );
        currentUserId = profile.id;
        isSystemAdmin = avatar?.isSystemAdmin ?? false;
    });

    function openView(t: Template) {
        active = t;
        editMode = false;
        modalOpen = true;
    }
    function closeModal() {
        modalOpen = false;
        editMode = false;
        deleteOpen = false;
    }
    function startEdit() {
        if (!active) return;
        editForm = {
            title: active.title || "",
            description: active.description || "",
            topics: (active.topics || []).map((x: Topic) => ({
                question: x.question,
                description: x.description || "",
                isRequire: x.isRequire,
                displayNo: x.displayNo,
            })),
        };
        editMode = true;
    }
    function askDelete() {
        deleteOpen = true;
    }
    function addTopic() {
        editForm.topics = [
            ...editForm.topics,
            {
                question: "",
                description: "",
                isRequire: true,
                displayNo: editForm.topics.length + 1,
            },
        ];
    }
    function removeTopic(i: number) {
        if (editForm.topics.length <= 1) return;
        editForm.topics = editForm.topics.filter(
            (_: any, idx: number) => idx !== i,
        );
    }
    async function saveEdit() {
        if (!active) return;
        const validationError = validateTemplateForm(
            editForm.title,
            editForm.description,
            editForm.topics,
        );
        if (validationError) {
            showToast(validationError, "error");
            return;
        }
        const topics = editForm.topics.map((t: Topic, i: number) => ({
            ...t,
            displayNo: i + 1,
        }));
        try {
            await api("PUT", `/api/v1/templates/${active.uuid || active.id}`, {
                title: editForm.title,
                description: editForm.description,
                topics,
            });
            showToast("保存しました");
            await reloadTemplates();
            closeModal();
        } catch (e) {
            showToast(
                e instanceof Error ? e.message : "保存に失敗しました",
                "error",
            );
        }
    }
    async function confirmDelete() {
        if (!active) return;
        try {
            await api(
                "DELETE",
                `/api/v1/templates/${active.uuid || active.id}`,
            );
            await reloadTemplates();
            closeModal();
        } catch {
            showToast("削除に失敗しました", "error");
        }
    }
    async function setTemplateLock(isLock: boolean) {
        if (!active) return;
        try {
            await api(
                "PUT",
                `/api/v1/templates/${active.uuid || active.id}/lock`,
                { isLock },
            );
            showToast(isLock ? "ロックしました" : "ロックを解除しました");
            await reloadTemplates();
        } catch (e) {
            showToast(
                e instanceof Error ? e.message : "ロックの更新に失敗しました",
                "error",
            );
        }
    }
    async function copy(t: Template) {
        try {
            await api("POST", `/api/v1/templates/${t.uuid || t.id}/copy`, {});
            showToast("コピーしました");
            await reloadTemplates();
        } catch {
            showToast("コピーに失敗しました", "error");
        }
    }
    async function moveUp(t: Template) {
        const i = templates.findIndex((x: Template) => x.id === t.id);
        if (i <= 0) return;
        const next = [...templates];
        [next[i - 1], next[i]] = [next[i], next[i - 1]];
        templates = next;
        await persistOrder(next);
    }
    async function moveDown(t: Template) {
        const i = templates.findIndex((x: Template) => x.id === t.id);
        if (i < 0 || i >= templates.length - 1) return;
        const next = [...templates];
        [next[i], next[i + 1]] = [next[i + 1], next[i]];
        templates = next;
        await persistOrder(next);
    }
</script>

<svelte:head><title>テンプレート - Growlink</title></svelte:head>

<svelte:window
    onkeydown={(e: KeyboardEvent) => e.key === "Escape" && closeModal()} />

<div class="flex flex-wrap items-center justify-between gap-2">
    <div>
        <h1 class="font-display text-2xl font-semibold">テンプレート一覧</h1>
        <p class="mt-1 text-sm text-base-content/60">
            行をクリックして内容を確認できます
        </p>
    </div>
    <dir class="flex flex-1 justify-end items-end">
        <a href="/planning/templates/new" class="btn btn-primary btn-sm"
            >新規作成</a>
    </dir>
</div>

<PlanningTabs active="templates" />

<div class="mt-6 space-y-2">
    {#if templates.length === 0}
        <div
            class="rounded-3xl border border-base-300 bg-base-100 p-12 text-center text-sm text-base-content/50">
            テンプレートがありません
        </div>
    {:else}
        {#each templates.entries() as [idx, t] (t.id)}
            <div
                role="button"
                tabindex="0"
                class="group cursor-pointer rounded-2xl border border-base-300 bg-base-100 p-4 shadow-sm transition-all hover:border-primary/50 hover:shadow-md"
                onclick={() => openView(t)}
                onkeydown={(e) => e.key === "Enter" && openView(t)}>
                <div class="flex min-w-0 items-center justify-between gap-3">
                    <div class="flex min-w-0 items-center gap-2">
                        <span class="badge badge-ghost shrink-0">#{idx}</span>
                        <span class="truncate font-medium">{t.title}</span>
                        {#if t.isSample}
                            <span class="badge badge-outline shrink-0"
                                >サンプル</span>
                        {/if}
                        {#if t.isLock}
                            <span class="badge badge-warning badge-outline shrink-0"
                                >ロック</span>
                        {/if}
                    </div>
                    <div
                        class="flex shrink-0 items-center gap-2"
                        onclick={(e) => e.stopPropagation()}
                        role="presentation">
                        <span class="text-xs text-base-content/50"
                            >作成者: {t.creatorName || "—"}</span>
                        <button
                            type="button"
                            class="btn btn-xs btn-outline"
                            onclick={() => copy(t)}>コピー</button>
                        <button
                            type="button"
                            class="btn btn-xs btn-ghost"
                            onclick={() => moveUp(t)}>↑</button>
                        <button
                            type="button"
                            class="btn btn-xs btn-ghost"
                            onclick={() => moveDown(t)}>↓</button>
                    </div>
                </div>
                <p
                    class="mt-1 text-sm text-base-content/60 whitespace-pre-wrap line-clamp-3">
                    {t.description?.trim() || "—"}
                </p>
            </div>
        {/each}
    {/if}
</div>

<dialog class="modal" class:modal-open={modalOpen}>
    <div class="modal-box flex max-h-[85vh] max-w-2xl flex-col p-0">
        <div
            class="flex shrink-0 items-start justify-between gap-2 border-b border-base-200 p-6 pb-2">
            <div class="min-w-0">
                {#if editMode}
                    <h3 class="truncate text-lg font-bold">テンプレート編集</h3>
                {:else}
                    <h3 class="truncate text-lg font-bold">{active?.title}</h3>
                    {#if active?.description}
                        <p
                            class="mt-1 text-sm text-base-content/60 whitespace-pre-wrap">
                            {active.description}
                        </p>
                    {/if}
                {/if}
            </div>
            <div class="flex shrink-0 gap-1">
                {#if !editMode && canToggleLock(active)}
                    {#if active?.isLock}
                        <button
                            type="button"
                            class="btn btn-sm btn-warning btn-outline"
                            onclick={() => setTemplateLock(false)}
                            >ロック解除</button>
                    {:else}
                        <button
                            type="button"
                            class="btn btn-sm btn-warning btn-outline"
                            onclick={() => setTemplateLock(true)}
                            >ロック</button>
                    {/if}
                {/if}
                {#if !editMode && canEditDelete(active)}
                    <button
                        type="button"
                        class="btn btn-sm btn-outline"
                        onclick={startEdit}>編集</button>
                {/if}
                {#if canEditDelete(active)}
                    <button
                        type="button"
                        class="btn btn-sm btn-error btn-outline"
                        onclick={askDelete}>削除</button>
                {/if}
                <button
                    type="button"
                    class="btn btn-sm btn-circle btn-ghost"
                    aria-label="閉じる"
                    onclick={closeModal}>×</button>
            </div>
        </div>
        {#if !editMode}
            <div class="flex-1 overflow-y-auto px-6 py-4">
                {#each active?.topics || [] as topic (topic.id || topic.displayNo)}
                    <div
                        class="mb-3 rounded-2xl border border-base-300 bg-base-200/30 p-4">
                        <div class="flex items-start gap-2">
                            <span
                                class="badge badge-primary badge-outline shrink-0"
                                >Q{topic.displayNo}</span>
                            <div class="min-w-0 flex-1">
                                <p class="text-sm leading-relaxed">
                                    {topic.question}
                                </p>
                                {#if topic.description}
                                    <p
                                        class="mt-1 whitespace-pre-wrap text-xs text-base-content/60">
                                        {topic.description}
                                    </p>
                                {/if}
                            </div>
                            {#if topic.isRequire}
                                <span
                                    class="badge badge-error badge-outline badge-sm shrink-0"
                                    >必須</span>
                            {/if}
                        </div>
                    </div>
                {/each}
            </div>
        {:else}
            <div class="flex-1 space-y-3 overflow-y-auto px-6 py-4">
                <label class="form-control">
                    <span class="label-text font-medium">タイトル</span>
                    <input
                        class="input input-bordered"
                        bind:value={editForm.title} />
                </label>
                <label class="form-control">
                    <span class="label-text font-medium">説明</span>
                    <textarea
                        class="textarea textarea-bordered min-h-[4rem]"
                        rows="3"
                        bind:value={editForm.description}
                        placeholder="このテンプレートの目的や使い方"></textarea>
                </label>
                {#each editForm.topics as topic, i (i)}
                    <div
                        class="flex items-center gap-2 rounded-2xl border border-base-300 p-2">
                        <span class="badge badge-ghost shrink-0">Q{i + 1}</span>
                        <div class="min-w-0 flex-1 space-y-1">
                            <input
                                class="input input-bordered input-sm w-full placeholder-current"
                                bind:value={topic.question}
                                placeholder="質問文" />
                            <textarea
                                class="textarea textarea-bordered textarea-sm min-h-[2.5rem] w-full placeholder-current"
                                rows="2"
                                bind:value={topic.description}
                                placeholder="説明（改行可能）"></textarea>
                        </div>
                        <label class="label shrink-0 cursor-pointer gap-1 py-0">
                            <input
                                type="checkbox"
                                class="checkbox checkbox-sm"
                                bind:checked={topic.isRequire} />
                            <span class="label-text text-xs">必須</span>
                        </label>
                        <button
                            type="button"
                            class="btn btn-ghost btn-xs btn-square text-error"
                            title="削除"
                            disabled={editForm.topics.length <= 1}
                            onclick={() => removeTopic(i)}>×</button>
                    </div>
                {/each}
                <button
                    type="button"
                    class="btn btn-outline btn-sm w-full"
                    onclick={addTopic}>質問を追加</button>
            </div>
        {/if}
        <div class="modal-action shrink-0 border-t border-base-200 p-4">
            <button type="button" class="btn btn-ghost" onclick={closeModal}
                >閉じる</button>
            {#if editMode}
                <button type="button" class="btn btn-primary" onclick={saveEdit}
                    >保存</button>
            {/if}
        </div>
    </div>
    <form method="dialog" class="modal-backdrop">
        <button type="button" onclick={closeModal}>close</button>
    </form>
</dialog>

<dialog class="modal" class:modal-open={deleteOpen}>
    <div class="modal-box max-h-[85vh]">
        <h3 class="font-bold">削除の確認</h3>
        <p class="py-4 text-sm text-base-content/80">
            このテンプレートを削除しますか？
        </p>
        <div class="modal-action">
            <button
                type="button"
                class="btn btn-ghost"
                onclick={() => (deleteOpen = false)}>キャンセル</button>
            <button type="button" class="btn btn-error" onclick={confirmDelete}
                >削除確定</button>
        </div>
    </div>
    <form method="dialog" class="modal-backdrop">
        <button type="button" onclick={() => (deleteOpen = false)}
            >close</button>
    </form>
</dialog>
