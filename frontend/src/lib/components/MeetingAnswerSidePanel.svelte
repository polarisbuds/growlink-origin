<script lang="ts">
    import { api } from "$lib/api";
    import { showToast } from "$lib/toast";
    import { formatDate, type MeetingUserRef } from "$lib/meeting-helpers";
    import SharerSelectModal, {
        type SharerPickerUser,
    } from "$lib/components/SharerSelectModal.svelte";
    import PanelRightClose from "lucide-svelte/icons/panel-right-close";
    import PanelRightOpen from "lucide-svelte/icons/panel-right-open";

    type DirectoryTeam = {
        uuid: string;
        name: string;
        departmentUuid: string;
    };

    let {
        open = $bindable(true),
        uuid,
        implementedDate = $bindable(""),
        menteeLabel,
        sharers = $bindable([] as MeetingUserRef[]),
        canEditSharers = false,
        canEditImplementedDate = false,
    }: {
        open?: boolean;
        uuid: string | undefined;
        implementedDate?: string;
        menteeLabel: string;
        sharers?: MeetingUserRef[];
        canEditSharers?: boolean;
        canEditImplementedDate?: boolean;
    } = $props();

    let sharerModalOpen = $state(false);
    let directoryLoaded = $state(false);
    let directoryUsers = $state<SharerPickerUser[]>([]);
    let departments = $state<{ uuid: string; name: string }[]>([]);
    let teams = $state<DirectoryTeam[]>([]);
    let applyingSharers = $state(false);

    let sharerUserUuids = $derived(
        (sharers || [])
            .map((s) => s.uuid)
            .filter((id): id is string => !!id),
    );

    async function loadDirectoryIfNeeded() {
        if (directoryLoaded) return;
        const dir = await api<{
            users?: SharerPickerUser[];
            departments?: { uuid: string; name: string }[];
            teams?: DirectoryTeam[];
        }>("GET", "/api/v1/planning/directory");
        directoryUsers = dir.users || [];
        departments = dir.departments || [];
        teams = (dir.teams || []).map((t) => ({
            uuid: t.uuid,
            name: t.name,
            departmentUuid: t.departmentUuid,
        }));
        directoryLoaded = true;
    }

    async function openSharerModal() {
        if (!uuid) return;
        await loadDirectoryIfNeeded();
        sharerModalOpen = true;
    }

    async function applySharers(uuids: string[]) {
        if (!uuid) return;
        applyingSharers = true;
        try {
            const res = await api<{ sharers: MeetingUserRef[] }>(
                "PUT",
                `/api/v1/meetings/${uuid}/sharers`,
                { sharerUserUuids: uuids },
            );
            sharers = res.sharers || [];
            showToast("共有者を保存しました");
            sharerModalOpen = false;
        } finally {
            applyingSharers = false;
        }
    }

    function saveImplementedDate() {
        showToast("実施日の API は今後接続されます（表示のみ更新可）", "info");
    }
</script>

<SharerSelectModal
    open={sharerModalOpen}
    users={directoryUsers}
    {departments}
    {teams}
    selectedUuids={sharerUserUuids}
    pickRootId="meeting-answer-sharer-name"
    onApply={applySharers}
    onClose={() => {
        sharerModalOpen = false;
    }} />

{#if open}
    <aside
        class="hidden w-72 shrink-0 lg:block rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm h-fit sticky top-4">
        <div class="flex items-center justify-between gap-2">
            <h2 class="text-sm font-semibold">面談情報</h2>
            <button
                type="button"
                class="btn btn-ghost btn-xs btn-square"
                aria-label="パネルを閉じる"
                onclick={() => (open = false)}>
                <PanelRightClose class="h-4 w-4" />
            </button>
        </div>
        <dl class="mt-4 space-y-4 text-sm">
            <div>
                <dt class="text-base-content/60">対象者</dt>
                <dd class="mt-0.5 font-medium">{menteeLabel}</dd>
            </div>
            <div>
                <dt class="text-base-content/60">共有者</dt>
                <dd class="mt-1">
                    {#if sharers.length}
                        <ul class="space-y-1">
                            {#each sharers as s (s.uuid)}
                                <li>{s.name || "—"}</li>
                            {/each}
                        </ul>
                    {:else}
                        <p class="text-base-content/50">未設定</p>
                    {/if}
                    {#if canEditSharers}
                        <button
                            type="button"
                            class="btn btn-outline btn-sm mt-3 w-full"
                            disabled={applyingSharers}
                            onclick={openSharerModal}>
                            {applyingSharers ? "保存中…" : "共有者を変更"}
                        </button>
                    {/if}
                </dd>
            </div>
            <div>
                <dt class="text-base-content/60">実施予定日</dt>
                <dd class="mt-1">
                    {#if canEditImplementedDate}
                        <input
                            type="date"
                            class="input input-bordered input-sm w-full"
                            bind:value={implementedDate} />
                        <button
                            type="button"
                            class="btn btn-outline btn-sm mt-2 w-full"
                            onclick={saveImplementedDate}>
                            実施予定日を保存
                        </button>
                    {:else}
                        <span class="font-medium"
                            >{formatDate(implementedDate)}</span>
                    {/if}
                </dd>
            </div>
        </dl>
    </aside>
{:else}
    <div class="hidden lg:flex shrink-0">
        <button
            type="button"
            class="btn btn-ghost btn-sm btn-square sticky top-4 h-fit"
            aria-label="面談情報パネルを開く"
            onclick={() => (open = true)}>
            <PanelRightOpen class="h-4 w-4" />
        </button>
    </div>
{/if}
