<script lang="ts">
    import { api } from "$lib/api";
    import { showToast } from "$lib/toast";
    import { formatDate, type MeetingUserRef } from "$lib/meeting-helpers";
    import UserChipSearchField from "$lib/components/UserChipSearchField.svelte";
    import PanelRightClose from "lucide-svelte/icons/panel-right-close";
    import PanelRightOpen from "lucide-svelte/icons/panel-right-open";

    type User = { uuid: string; name?: string; email?: string };

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

    let sharerUserUuids = $state<string[]>([]);
    let sharerPickQuery = $state("");
    let candidateUsers = $state<User[]>([]);
    let savingSharers = $state(false);

    let sharerChips = $derived(
        sharerUserUuids.map((id) => {
            const fromSharers = sharers.find((s) => s.uuid === id);
            return { uuid: id, name: fromSharers?.name || `ユーザ` };
        }),
    );

    $effect(() => {
        sharerUserUuids = (sharers || [])
            .map((s) => s.uuid)
            .filter((id): id is string => !!id);
    });

    let searchTimer: ReturnType<typeof setTimeout> | undefined;
    $effect(() => {
        if (!canEditSharers || !open) return;
        const q = sharerPickQuery;
        clearTimeout(searchTimer);
        searchTimer = setTimeout(async () => {
            candidateUsers = await api<User[]>(
                "GET",
                `/api/v1/meetings/${uuid}/sharer-candidates?q=${encodeURIComponent(q)}`,
            );
        }, 200);
        return () => clearTimeout(searchTimer);
    });

    function addSharer(u: User) {
        if (sharerUserUuids.includes(u.uuid)) return;
        sharerUserUuids = [...sharerUserUuids, u.uuid];
        sharers = [
            ...sharers,
            { uuid: u.uuid, name: u.name, email: u.email },
        ];
    }

    function removeSharer(uuid: string) {
        sharerUserUuids = sharerUserUuids.filter((x) => x !== uuid);
        sharers = sharers.filter((s) => s.uuid !== uuid);
    }

    async function saveSharers() {
        savingSharers = true;
        try {
            const res = await api<{ sharers: MeetingUserRef[] }>(
                "PUT",
                `/api/v1/meetings/${uuid}/sharers`,
                { sharerUserUuids },
            );
            sharers = res.sharers || [];
            showToast("共有者を保存しました");
        } finally {
            savingSharers = false;
        }
    }

    function saveImplementedDate() {
        showToast("実施日の API は今後接続されます（表示のみ更新可）", "info");
    }
</script>

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
                    {#if canEditSharers}
                        <UserChipSearchField
                            pickRootId="meeting-answer-sharer"
                            chips={sharerChips}
                            items={candidateUsers}
                            bind:query={sharerPickQuery}
                            placeholder={sharerChips.length
                                ? "追加で検索"
                                : "名前で検索して追加"}
                            onPick={addSharer}
                            onRemove={removeSharer} />
                        <button
                            type="button"
                            class="btn btn-primary btn-sm mt-3 w-full"
                            disabled={savingSharers}
                            onclick={saveSharers}>
                            {savingSharers ? "保存中…" : "共有者を保存"}
                        </button>
                    {:else if sharers.length}
                        <ul class="space-y-1">
                            {#each sharers as s (s.uuid)}
                                <li>{s.name || "—"}</li>
                            {/each}
                        </ul>
                    {:else}
                        <p class="text-base-content/50">未設定</p>
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
