<script lang="ts">
    import { api } from "$lib/api";
    import { showToast } from "$lib/toast";
    import Icon from "$lib/components/Icon.svelte";
    import SharerSelectModal, {
        type SharerPickerUser,
    } from "$lib/components/SharerSelectModal.svelte";
    import {
        lastPageForTotal,
        visiblePageNumbers,
    } from "$lib/planning-pairings-pager";

    const TABLE_PAGE_SIZE = 10;

    export type PairShareTarget = {
        pairingUuid: string;
        mentorUserUuid: string;
        menteeUserUuid: string;
        mentorName: string;
        menteeName: string;
    };

    type Dept = { uuid: string; name: string };
    type Team = { uuid: string; name: string; departmentUuid: string };

    let {
        open = false,
        targets = [],
        users,
        departments = [],
        teams = [],
        onClose,
    }: {
        open?: boolean;
        targets?: PairShareTarget[];
        users: SharerPickerUser[];
        departments?: Dept[];
        teams?: Team[];
        onClose: () => void;
    } = $props();

    let dialogEl = $state<HTMLDialogElement | null>(null);
    let sharerUserUuids = $state<string[]>([]);
    let pickerOpen = $state(false);
    let tablePage = $state(1);
    let loading = $state(false);
    let saving = $state(false);
    let loadGen = 0;

    let excludedUserUuids = $derived.by(() => {
        const set = new Set<string>();
        for (let i = 0; i < targets.length; i++) {
            const t = targets[i];
            if (t.mentorUserUuid) set.add(t.mentorUserUuid);
            if (t.menteeUserUuid) set.add(t.menteeUserUuid);
        }
        return set;
    });

    let pickableUsers = $derived(
        users.filter((u) => !excludedUserUuids.has(u.uuid)),
    );

    let subtitle = $derived.by(() => {
        if (targets.length === 1) {
            return `${targets[0].mentorName} × ${targets[0].menteeName}`;
        }
        if (targets.length > 1) {
            return `選択した ${targets.length} 件`;
        }
        return "";
    });

    let selectedRows = $derived.by(() => {
        const rows = sharerUserUuids
            .map((uuid) => users.find((u) => u.uuid === uuid))
            .filter((u): u is SharerPickerUser => u != null);
        rows.sort((a, b) => {
            const c = (a.departmentNames || "").localeCompare(
                b.departmentNames || "",
                "ja",
            );
            if (c !== 0) return c;
            return a.name.localeCompare(b.name, "ja");
        });
        return rows;
    });

    let tableLastPage = $derived(
        lastPageForTotal(selectedRows.length, TABLE_PAGE_SIZE),
    );
    let tablePagerPages = $derived(
        visiblePageNumbers(tablePage, tableLastPage),
    );
    let pageRows = $derived(
        selectedRows.slice(
            (tablePage - 1) * TABLE_PAGE_SIZE,
            tablePage * TABLE_PAGE_SIZE,
        ),
    );

    $effect(() => {
        if (tablePage > tableLastPage) tablePage = tableLastPage;
    });

    $effect(() => {
        sharerUserUuids.length;
        tablePage = 1;
    });

    $effect(() => {
        if (!open || targets.length === 0) return;
        const gen = ++loadGen;
        loading = targets.length === 1;
        sharerUserUuids = [];
        tablePage = 1;
        if (targets.length !== 1) {
            loading = false;
            return;
        }
        const pairingUuid = targets[0].pairingUuid;
        void (async () => {
            try {
                const res = await api<{
                    sharers?: { uuid?: string }[];
                }>("GET", `/api/v1/pairings/${pairingUuid}/sharers`);
                if (gen !== loadGen) return;
                const uuids: string[] = [];
                const list = res.sharers || [];
                for (let idx = 0; idx < list.length; idx++) {
                    const u = list[idx]?.uuid;
                    if (u) uuids[idx] = u;
                }
                sharerUserUuids = uuids.filter(Boolean);
            } catch {
                if (gen === loadGen) {
                    showToast("共有者の取得に失敗しました", "error");
                    onClose();
                }
            } finally {
                if (gen === loadGen) loading = false;
            }
        })();
    });

    function avatarURL(u: SharerPickerUser | undefined) {
        if (!u?.imageUri) return "";
        return `/uploads/${u.imageUri}`;
    }

    async function save() {
        if (targets.length === 0 || saving) return;
        saving = true;
        try {
            for (let i = 0; i < targets.length; i++) {
                const pairingUuid = targets[i].pairingUuid;
                await api("PUT", `/api/v1/pairings/${pairingUuid}/sharers`, {
                    sharerUserUuids,
                });
            }
            showToast(
                targets.length === 1
                    ? "共有者を保存しました"
                    : `${targets.length} 件のペアに共有者を保存しました`,
            );
            onClose();
        } catch {
            showToast("保存に失敗しました", "error");
        } finally {
            saving = false;
        }
    }

    function discardClose() {
        onClose();
    }

    function onWindowKeydown(e: KeyboardEvent) {
        if (!open || e.key !== "Escape") return;
        if (pickerOpen) return;
        if (!dialogEl) return;
        if (document.querySelector("[data-pick-panel]")) return;
        discardClose();
    }
</script>

<svelte:window onkeydown={onWindowKeydown} />

<dialog bind:this={dialogEl} class="modal" class:modal-open={open}>
    <div
        class="modal-box flex max-h-[85vh] max-w-2xl flex-col !overflow-hidden p-0">
        <div
            class="flex shrink-0 items-center justify-between gap-2 border-b border-base-200 px-6 py-4">
            <div class="min-w-0">
                <h3 class="text-lg font-bold">ペア共有</h3>
                <p class="mt-1 truncate text-sm text-base-content/60">
                    {subtitle}
                </p>
            </div>
            <button
                type="button"
                class="btn btn-sm btn-circle btn-ghost"
                aria-label="閉じる"
                onclick={discardClose}>×</button>
        </div>

        <div class="min-h-0 flex-1 overflow-y-auto px-6 py-4 space-y-4">
            <div class="flex items-end justify-between gap-2">
                <span class="text-sm font-medium">共有者</span>
                <button
                    type="button"
                    class="btn btn-outline btn-xs"
                    disabled={loading}
                    onclick={() => (pickerOpen = true)}>共有者を追加</button>
            </div>

            {#if loading}
                <p class="text-sm text-base-content/50">読み込み中…</p>
            {:else if selectedRows.length === 0}
                <p class="text-sm text-base-content/50">未選択</p>
            {:else}
                <div class="overflow-hidden rounded-2xl border border-base-300">
                    <table class="table table-sm">
                        <thead>
                            <tr>
                                <th>部署</th>
                                <th>氏名</th>
                            </tr>
                        </thead>
                        <tbody>
                            {#each pageRows as u (u.uuid)}
                                <tr>
                                    <td>{u.departmentNames || "所属無し"}</td>
                                    <td>
                                        <span class="flex items-center gap-2">
                                            <Icon
                                                url={avatarURL(u)}
                                                text={u.name}
                                                width={28}
                                                height={28} />
                                            {u.name}
                                        </span>
                                    </td>
                                </tr>
                            {/each}
                        </tbody>
                    </table>
                </div>
                {#if selectedRows.length > TABLE_PAGE_SIZE}
                    <div
                        class="flex flex-wrap items-center justify-center gap-2 text-sm">
                        <button
                            type="button"
                            class="btn btn-sm btn-ghost"
                            disabled={tablePage <= 1}
                            onclick={() =>
                                (tablePage = Math.max(1, tablePage - 1))}
                            >前へ</button>
                        {#each tablePagerPages.entries() as [idx, p] (typeof p === "number" ? p : `ellipsis-${idx}`)}
                            {#if p === "ellipsis"}
                                <span class="w-8 px-1 text-base-content/50"
                                    >…</span>
                            {:else}
                                <button
                                    type="button"
                                    class="btn btn-sm {p === tablePage
                                        ? 'btn-primary'
                                        : 'btn-ghost'}"
                                    onclick={() => (tablePage = p)}>{p}</button>
                            {/if}
                        {/each}
                        <button
                            type="button"
                            class="btn btn-sm btn-ghost"
                            disabled={tablePage >= tableLastPage}
                            onclick={() =>
                                (tablePage = Math.min(
                                    tableLastPage,
                                    tablePage + 1,
                                ))}>次へ</button>
                        <span class="ml-2 text-base-content/60">
                            全 {selectedRows.length} 件
                        </span>
                    </div>
                {:else}
                    <p class="text-xs text-base-content/60">
                        選択済み: {selectedRows.length} 人
                    </p>
                {/if}
            {/if}
        </div>

        <div
            class="modal-action shrink-0 border-t border-base-200 p-4 flex justify-end gap-2">
            <button
                type="button"
                class="btn btn-ghost"
                disabled={saving}
                onclick={discardClose}>キャンセル</button>
            <button
                type="button"
                class="btn btn-primary"
                disabled={loading || saving}
                onclick={save}>{saving ? "保存中…" : "保存"}</button>
        </div>
    </div>
    <form method="dialog" class="modal-backdrop">
        <button type="button" onclick={discardClose}>close</button>
    </form>
</dialog>

<SharerSelectModal
    open={pickerOpen}
    users={pickableUsers}
    {departments}
    {teams}
    selectedUuids={sharerUserUuids}
    pickRootId="pair-sharer-select"
    onApply={(uuids) => {
        sharerUserUuids = uuids;
        pickerOpen = false;
    }}
    onClose={() => {
        pickerOpen = false;
    }} />