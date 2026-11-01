<script lang="ts">
    import { api } from "$lib/api";
    import { showToast } from "$lib/toast";
    import { formatDate, type MeetingUserRef } from "$lib/meeting-helpers";
    import UserChipSearchField from "$lib/components/UserChipSearchField.svelte";
    import PanelRightClose from "lucide-svelte/icons/panel-right-close";
    import PanelRightOpen from "lucide-svelte/icons/panel-right-open";

    type User = { id: number; name: string; email: string };

    let {
        open = $bindable(true),
        uuid,
        implementedAt,
        menteeLabel,
        sharers = $bindable([] as MeetingUserRef[]),
        canEditSharers = false,
    }: {
        open?: boolean;
        uuid: string;
        implementedAt?: string;
        menteeLabel: string;
        sharers?: MeetingUserRef[];
        canEditSharers?: boolean;
    } = $props();

    let sharerUserIds = $state<number[]>([]);
    let sharerPickOpen = $state(false);
    let sharerPickQuery = $state("");
    let candidateUsers = $state<User[]>([]);
    let savingSharers = $state(false);

    let sharerChips = $derived(
        sharerUserIds.map((id) => {
            const fromSharers = sharers.find((s) => s.id === id);
            return { id, name: fromSharers?.name || `ユーザ #${id}` };
        }),
    );

    let sharerPickFiltered = $derived.by(() => {
        const picked = new Set(sharerUserIds);
        return candidateUsers.filter((u) => !picked.has(u.id));
    });

    $effect(() => {
        sharerUserIds = (sharers || [])
            .map((s) => s.id)
            .filter((id): id is number => id != null);
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
        if (sharerUserIds.includes(u.id)) return;
        sharerUserIds = [...sharerUserIds, u.id];
        sharers = [
            ...sharers,
            { id: u.id, name: u.name, email: u.email },
        ];
        sharerPickQuery = "";
        sharerPickOpen = false;
    }

    function removeSharer(id: number) {
        sharerUserIds = sharerUserIds.filter((x) => x !== id);
        sharers = sharers.filter((s) => s.id !== id);
    }

    async function saveSharers() {
        savingSharers = true;
        try {
            const res = await api<{ sharers: MeetingUserRef[] }>(
                "PUT",
                `/api/v1/meetings/${uuid}/sharers`,
                { sharerUserIds },
            );
            sharers = res.sharers || [];
            showToast("共有者を保存しました");
        } finally {
            savingSharers = false;
        }
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
                <dt class="text-base-content/60">実施日</dt>
                <dd class="mt-0.5 font-medium">{formatDate(implementedAt)}</dd>
            </div>
            <div>
                <dt class="text-base-content/60">対象者</dt>
                <dd class="mt-0.5 font-medium">{menteeLabel}</dd>
            </div>
            <div>
                <dt class="text-base-content/60">共有者</dt>
                <dd class="mt-1">
                    {#if canEditSharers}
                        <div data-pick-root="meeting-answer-sharer">
                            <UserChipSearchField
                                chips={sharerChips}
                                placeholder={sharerChips.length
                                    ? "追加で検索"
                                    : "名前で検索して追加"}
                                bind:query={sharerPickQuery}
                                bind:open={sharerPickOpen}
                                filteredUsers={sharerPickFiltered}
                                onPick={addSharer}
                                onRemove={removeSharer} />
                        </div>
                        <button
                            type="button"
                            class="btn btn-primary btn-sm mt-3 w-full"
                            disabled={savingSharers}
                            onclick={saveSharers}>
                            {savingSharers ? "保存中…" : "共有者を保存"}
                        </button>
                    {:else if sharers.length}
                        <ul class="space-y-1">
                            {#each sharers as s (s.id)}
                                <li>{s.name || "—"}</li>
                            {/each}
                        </ul>
                    {:else}
                        <p class="text-base-content/50">未設定</p>
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
