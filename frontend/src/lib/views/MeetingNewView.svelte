<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { api } from "$lib/api";
    import { getAvatar } from "$lib/session";
    import {
        matchPairingIds,
        type MeetingPairingItem,
    } from "$lib/meeting-pair-match";
    import { expandUserPairings } from "$lib/pairing-helpers";
    import {
        buildUserDepartmentMaps,
        buildUserTeamMaps,
        departmentNameForUser,
        teamNameForUser,
    } from "$lib/planning-directory";
    import MultiCheckSelect from "$lib/components/MultiCheckSelect.svelte";
    import SearchableSelect from "$lib/components/SearchableSelect.svelte";

    const multiPickInputClass =
        "input input-bordered w-full text-left font-normal placeholder-current";

    type Template = { id: number; title?: string };
    type User = {
        id: number;
        name: string;
        email: string;
        departmentNames?: string;
    };
    type ManagedTeam = { id: number; name: string; members?: User[] };
    type CreateMode = "planned" | "pair" | "monthly";

    function todayISO() {
        const d = new Date();
        return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`;
    }
    function currentMonthISO() {
        const d = new Date();
        return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}`;
    }

    let pairings = $state<MeetingPairingItem[]>([]);
    let pairingsLoaded = $state(false);
    let pairingsLoading = $state(false);
    let selectedPairingIds = $state<number[]>([]);
    let pairSearch = $state("");
    let templates = $state<Template[]>([]);
    let users = $state<User[]>([]);
    let currentUserId = $state<number | undefined>();
    let canSpecifyPair = $state(false);
    let managedTeams = $state<ManagedTeam[]>([]);
    let mode = $state<CreateMode>("planned");
    let mentorId = $state(0);
    let menteeIds = $state<number[]>([]);
    let sharerUserIds = $state<number[]>([]);
    let pickableUsers = $derived(
        currentUserId != null
            ? users.filter((u) => u.id !== currentUserId)
            : users,
    );
    let menteeTargetPicks = $derived(
        menteeIds.map((targetId) => ({ kind: "user", targetId })),
    );
    let submitting = $state(false);
    let error = $state("");

    let form = $state({
        templateId: 0,
        implementedAt: todayISO(),
        fromMonth: currentMonthISO(),
        toMonth: currentMonthISO(),
        week: 2,
        weekday: 3,
    });

    let managedTeamPickValue = $state(0);

    let deptMaps = $derived(buildUserDepartmentMaps(users));
    let teamMaps = $derived(buildUserTeamMaps(users));
    let pairTableRows = $derived(expandUserPairings(pairings));
    let filteredPairTableRows = $derived.by(() => {
        const q = pairSearch.trim().toLowerCase();
        if (!q) return pairTableRows;
        return pairTableRows.filter((row) => {
            const dept = departmentNameForUser(row.menteeUserId, deptMaps);
            return (
                row.mentorName.toLowerCase().includes(q) ||
                row.menteeName.toLowerCase().includes(q) ||
                dept.toLowerCase().includes(q)
            );
        });
    });
    let displayedPairingIds = $derived([
        ...new Set(filteredPairTableRows.map((r) => r.pairing.id)),
    ]);
    let allPairRowsSelected = $derived(
        displayedPairingIds.length > 0 &&
            displayedPairingIds.every((id) => selectedPairingIds.includes(id)),
    );

    let matchedPairingCount = $derived(
        matchPairingIds(
            pairings,
            mentorId ? [mentorId] : [],
            menteeTargetPicks,
        ).length,
    );

    type SelectItem = { id: number };

    function userLabel(id: number) {
        const u = users.find((x) => x.id === id);
        return u ? u.name : `#${id}`;
    }

    function userLabelOf(item: SelectItem) {
        const u = users.find((x) => x.id === item.id);
        return u ? `${u.name} (${u.email})` : `#${item.id}`;
    }
    function templateLabelOf(item: SelectItem) {
        const t = templates.find((x) => x.id === item.id);
        return t?.title || "";
    }
    function managedTeamLabelOf(item: SelectItem) {
        const team = managedTeams.find((x) => x.id === item.id);
        if (!team) return `#${item.id}`;
        return `${team.name}（${team.members?.length ?? 0}人）`;
    }
    function pickManagedTeamFromSelect(item: SelectItem) {
        const team = managedTeams.find((t) => t.id === item.id);
        if (team) pickManagedTeam(team);
    }
    function pickManagedTeam(team: ManagedTeam) {
        const next = new Set(menteeIds);
        for (const m of team.members || []) {
            if (!m?.id) continue;
            if (currentUserId != null && m.id === currentUserId) continue;
            next.add(m.id);
        }
        menteeIds = [...next];
    }
    function submitLabel() {
        if (submitting) return "作成中…";
        if (mode === "planned") return "作成";
        if (mode === "pair") return "作成";
        return "定期作成";
    }
    function resolvePairingIds(): number[] {
        return matchPairingIds(
            pairings,
            mentorId ? [mentorId] : [],
            menteeTargetPicks,
        );
    }
    function toggleAllPairRows(checked: boolean) {
        selectedPairingIds = checked ? [...displayedPairingIds] : [];
    }

    async function loadPairingsIfNeeded() {
        if (pairingsLoaded || pairingsLoading) return;
        pairingsLoading = true;
        try {
            pairings = await api<MeetingPairingItem[]>(
                "GET",
                "/api/v1/pairings/meeting",
            );
            pairingsLoaded = true;
        } finally {
            pairingsLoading = false;
        }
    }

    async function onModeChange(next: CreateMode) {
        mode = next;
        if (next === "pair") {
            selectedPairingIds = [];
            await loadPairingsIfNeeded();
        } else if (next === "monthly") {
            await loadPairingsIfNeeded();
        }
    }

    onMount(async () => {
        const [t, dir, profile] = await Promise.all([
            api<Template[]>("GET", "/api/v1/templates"),
            api<{ users?: User[]; managedTeams?: ManagedTeam[] }>(
                "GET",
                "/api/v1/planning/directory",
            ),
            api<{ id?: number }>("GET", "/api/v1/setting/profile"),
        ]);
        templates = t;
        users = dir.users || [];
        managedTeams = dir.managedTeams || [];
        currentUserId = profile.id;
        canSpecifyPair = getAvatar()?.canSpecifyPair ?? false;
        form.templateId = templates[0]?.id || 0;
    });

    async function submit() {
        error = "";
        if (!form.templateId) {
            error = "テンプレートを選択してください。";
            return;
        }
        submitting = true;
        try {
            const sharerPayload = sharerUserIds.length
                ? sharerUserIds
                : undefined;

            if (mode === "planned") {
                if (!mentorId) {
                    error = "メンターを選択してください。";
                    return;
                }
                if (!menteeIds.length) {
                    error = "メンティを選択してください。";
                    return;
                }
                if (!form.implementedAt) {
                    error = "実施日を入力してください。";
                    return;
                }
                if (menteeIds.some((id) => id === mentorId)) {
                    error = "メンターと同じユーザをメンティに指定できません。";
                    return;
                }
                await api("POST", "/api/v1/meetings", {
                    mentorId,
                    menteeIds,
                    templateId: form.templateId,
                    implementedAt: form.implementedAt,
                    sharerUserIds: sharerPayload,
                });
            } else if (mode === "pair") {
                await loadPairingsIfNeeded();
                if (!pairTableRows.length) {
                    error = "ペアがありません。";
                    return;
                }
                if (!selectedPairingIds.length) {
                    error = "ペアを選択してください。";
                    return;
                }
                if (!form.implementedAt) {
                    error = "実施日を入力してください。";
                    return;
                }
                const pairingIds = [...new Set(selectedPairingIds)];
                await api("POST", "/api/v1/meetings/all", {
                    templateId: form.templateId,
                    implementedAt: form.implementedAt,
                    sharerUserIds: sharerPayload,
                    pairingIds,
                });
            } else {
                if (!mentorId) {
                    error = "メンターを選択してください。";
                    return;
                }
                if (!menteeIds.length) {
                    error = "メンティを選択してください。";
                    return;
                }
                await loadPairingsIfNeeded();
                const pairingIds = resolvePairingIds();
                if (!pairingIds.length) {
                    error = "該当するペアがありません。";
                    return;
                }
                if (!form.fromMonth || !form.toMonth) {
                    error = "期間を入力してください。";
                    return;
                }
                if (form.fromMonth > form.toMonth) {
                    error = "開始月は終了月以前にしてください。";
                    return;
                }
                if (form.week < 1 || form.week > 5) {
                    error = "第N週は 1〜5 で指定してください。";
                    return;
                }
                await api("POST", "/api/v1/meetings/monthly", {
                    templateId: form.templateId,
                    from: `${form.fromMonth}-01`,
                    to: `${form.toMonth}-01`,
                    week: form.week,
                    weekday: form.weekday,
                    sharerUserIds: sharerPayload,
                    pairingIds,
                });
            }
            await goto("/planning/meetings");
        } catch (e) {
            error = e instanceof Error ? e.message : "作成に失敗しました。";
        } finally {
            submitting = false;
        }
    }
</script>

<svelte:head><title>1on1作成 - Growlink</title></svelte:head>

<main class="mx-auto w-full max-w-6xl min-w-0 p-4 md:p-8">
    <h1 class="font-display text-2xl font-semibold">1on1ミーティング作成</h1>
    <p class="mt-1 text-sm text-base-content/60">
        作成方法を選び、メンター・メンティ・テンプレートを指定します
    </p>
    <div class="mt-4 flex flex-wrap gap-2">
        <button
            type="button"
            class="btn btn-sm"
            class:btn-primary={mode === "planned"}
            class:btn-outline={mode !== "planned"}
            onclick={() => onModeChange("planned")}>作成</button>
        {#if canSpecifyPair}
            <button
                type="button"
                class="btn btn-sm"
                class:btn-primary={mode === "pair"}
                class:btn-outline={mode !== "pair"}
                onclick={() => onModeChange("pair")}>ペア指定</button>
        {/if}
        <button
            type="button"
            class="btn btn-sm"
            class:btn-primary={mode === "monthly"}
            class:btn-outline={mode !== "monthly"}
            onclick={() => onModeChange("monthly")}>定期面談作成</button>
    </div>
    <div
        class="mt-6 space-y-4 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm">
        {#if error}
            <div class="alert alert-error text-sm">{error}</div>
        {/if}
        <div class="form-control">
            <span class="label-text font-medium">テンプレート</span>
            <div class="mt-1">
                <SearchableSelect
                    pickRootId="meeting-new-template"
                    items={templates}
                    bind:value={form.templateId}
                    labelOf={templateLabelOf}
                    placeholder="名前で検索して選択"
                    allowClear={true} />
            </div>
        </div>
        {#if mode === "pair"}
            <div>
                <span class="label-text font-medium">ペアを選択</span>
                <input
                    type="search"
                    class="input input-bordered mt-2 w-full"
                    placeholder="メンター・メンティ・部署で検索"
                    bind:value={pairSearch} />
                <div
                    class="mt-2 overflow-hidden rounded-2xl border border-base-300">
                    {#if pairingsLoading}
                        <div class="p-6 space-y-3">
                            {#each Array(4) as _, i (i)}
                                <div class="skeleton h-10 w-full rounded-xl">
                                </div>
                            {/each}
                        </div>
                    {:else if pairTableRows.length === 0}
                        <div
                            class="p-8 text-center text-sm text-base-content/50">
                            ペアがありません
                        </div>
                    {:else if filteredPairTableRows.length === 0}
                        <div
                            class="p-8 text-center text-sm text-base-content/50">
                            該当するペアがありません
                        </div>
                    {:else}
                        <div class="max-h-80 overflow-y-auto">
                            <table class="table table-pin-rows">
                                <thead>
                                    <tr class="text-base-content/70">
                                        <th class="w-10">
                                            <input
                                                type="checkbox"
                                                class="checkbox checkbox-sm"
                                                checked={allPairRowsSelected}
                                                onchange={(e) =>
                                                    toggleAllPairRows(
                                                        (
                                                            e.currentTarget as HTMLInputElement
                                                        ).checked,
                                                    )} />
                                        </th>
                                        <th>メンター</th>
                                        <th>メンティ</th>
                                        <th>部署</th>
                                        <th>チーム</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {#each filteredPairTableRows as row (row.rowKey)}
                                        <tr class="hover:bg-base-200/50">
                                            <td>
                                                <input
                                                    type="checkbox"
                                                    class="checkbox checkbox-sm"
                                                    value={row.pairing.id}
                                                    bind:group={
                                                        selectedPairingIds
                                                    } />
                                            </td>
                                            <td>{row.mentorName}</td>
                                            <td>{row.menteeName}</td>
                                            <td
                                                >{departmentNameForUser(
                                                    row.menteeUserId,
                                                    deptMaps,
                                                ) || "—"}</td>
                                            <td
                                                >{teamNameForUser(
                                                    row.menteeUserId,
                                                    teamMaps,
                                                )}</td>
                                        </tr>
                                    {/each}
                                </tbody>
                            </table>
                        </div>
                    {/if}
                </div>
            </div>
        {:else}
            <div class="grid items-end gap-4 md:grid-cols-3">
                <div class="form-control">
                    <span class="label-text font-medium">メンター</span>
                    <SearchableSelect
                        pickRootId="meeting-new-mentor"
                        items={users}
                        bind:value={mentorId}
                        labelOf={userLabelOf}
                        selectedLabel={userLabel}
                        placeholder="名前で検索" />
                </div>
                <div class="md:col-span-2">
                    <label class="form-control">
                        <span class="label-text font-medium">メンティ</span>
                        <MultiCheckSelect
                            pickRootId="meeting-new-mentee"
                            items={pickableUsers}
                            bind:selectedIds={menteeIds}
                            labelOf={userLabelOf}
                            placeholder="選択してください"
                            inputClass={multiPickInputClass} />
                    </label>
                    {#if managedTeams.length > 0}
                        <div class="form-control mt-3">
                            <span class="label-text font-medium"
                                >管理チームからまとめて追加</span>
                            <div class="mt-1">
                                <SearchableSelect
                                    pickRootId="meeting-new-managed-team"
                                    items={managedTeams}
                                    bind:value={managedTeamPickValue}
                                    labelOf={managedTeamLabelOf}
                                    placeholder="チーム名で検索"
                                    clearOnSelect={true}
                                    onSelect={pickManagedTeamFromSelect} />
                            </div>
                        </div>
                    {/if}
                </div>
            </div>
            {#if mode === "monthly" && mentorId && menteeIds.length}
                <p class="text-sm text-base-content/60">
                    一致するペア: <span class="font-medium text-base-content"
                        >{matchedPairingCount}</span> 件（各ペアに面談を1件作成）
                </p>
            {/if}
        {/if}
        <label class="form-control">
            <span class="label-text font-medium">共有者（任意・複数）</span>
            <MultiCheckSelect
                pickRootId="meeting-new-sharer"
                items={pickableUsers}
                bind:selectedIds={sharerUserIds}
                labelOf={userLabelOf}
                placeholder="選択してください"
                inputClass={multiPickInputClass} />
        </label>
        {#if mode !== "monthly"}
            <label class="form-control">
                <span class="label-text font-medium">実施日</span>
                <input
                    class="input input-bordered"
                    type="date"
                    bind:value={form.implementedAt} />
            </label>
        {:else}
            <div class="flex justify-between gap-2 space-y-4 items-baseline">
                <label class="form-control w-2/6">
                    <span class="label-text font-medium">年月（から）</span>
                    <input
                        class="input input-bordered"
                        type="month"
                        bind:value={form.fromMonth} />
                </label>
                <label class="form-control w-2/6">
                    <span class="label-text font-medium">年月（まで）</span>
                    <input
                        class="input input-bordered"
                        type="month"
                        bind:value={form.toMonth} />
                </label>
                <label class="form-control w-1/6">
                    <span class="label-text font-medium">実施週</span>
                    <input
                        class="input input-bordered"
                        type="number"
                        min="1"
                        max="5"
                        bind:value={form.week} />
                </label>
                <label class="form-control w-1/6">
                    <span class="label-text font-medium">曜日</span>
                    <select
                        class="select select-bordered"
                        bind:value={form.weekday}>
                        <option value={0}>日曜</option>
                        <option value={1}>月曜</option>
                        <option value={2}>火曜</option>
                        <option value={3}>水曜</option>
                        <option value={4}>木曜</option>
                        <option value={5}>金曜</option>
                        <option value={6}>土曜</option>
                    </select>
                </label>
            </div>
        {/if}
        <div class="flex flex-wrap gap-2 pt-2">
            <button
                type="button"
                class="btn btn-primary"
                disabled={submitting}
                onclick={submit}>{submitLabel()}</button>
            <a href="/planning/meetings" class="btn btn-ghost">戻る</a>
        </div>
    </div>
</main>
