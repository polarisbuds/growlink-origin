<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { api } from "$lib/api";
    import { getAvatar } from "$lib/session";
    import { type MeetingPairingItem } from "$lib/meeting-pair-match";
    import { expandUserPairings } from "$lib/pairing-helpers";
    import {
        buildUserDepartmentMaps,
        buildUserTeamMaps,
        departmentNameForUser,
        teamNameForUser,
    } from "$lib/planning-directory";
    import SearchableSelect from "$lib/components/SearchableSelect.svelte";
    import UserSelectModal from "$lib/components/UserSelectModal.svelte";
    import PairSelectModal from "$lib/components/PairSelectModal.svelte";
    import Icon from "$lib/components/Icon.svelte";

    const MENTEE_PAGE_SIZE = 10;

    type Template = { uuid: string; title?: string; isDefault?: boolean };
    type User = {
        uuid: string;
        name: string;
        email: string;
        departmentNames?: string;
        teamNames?: string;
        imageUri?: string;
    };
    type ManagedTeam = { uuid: string; name: string; members?: User[] };
    type DirectoryTeam = {
        uuid: string;
        name: string;
        departmentUuid: string;
    };
    type CreateTarget = "manual" | "pair";
    type ScheduleKind = "date" | "monthly";

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
    let selectedPairingUuids = $state<string[]>([]);
    let pairModalOpen = $state(false);
    let selectedPairTablePage = $state(1);
    let templates = $state<Template[]>([]);
    let users = $state<User[]>([]);
    let departments = $state<{ uuid: string; name: string }[]>([]);
    let teams = $state<DirectoryTeam[]>([]);
    let mentorModalOpen = $state(false);
    let menteeModalOpen = $state(false);
    let managedTeamModalOpen = $state(false);
    let sharerModalOpen = $state(false);
    let selectedMenteeTablePage = $state(1);
    let selectedSharerTablePage = $state(1);
    let currentUserUuid = $state("");
    let canSpecifyPair = $state(false);
    let managedTeams = $state<ManagedTeam[]>([]);
    let target = $state<CreateTarget>("manual");
    let schedule = $state<ScheduleKind>("date");
    let mentorUuid = $state("");
    let menteeUuids = $state<string[]>([]);
    let sharerUserUuids = $state<string[]>([]);
    let pickableUsers = $derived(
        mentorUuid ? users.filter((u) => u.uuid !== mentorUuid) : users,
    );
    let selectedMenteeRows = $derived.by(() => {
        const rows = menteeUuids
            .map((uuid) => users.find((u) => u.uuid === uuid))
            .filter((u): u is User => u != null);
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
    let selectedMenteePageCount = $derived(
        Math.max(1, Math.ceil(selectedMenteeRows.length / MENTEE_PAGE_SIZE)),
    );
    let selectedMenteePageRows = $derived(
        selectedMenteeRows.slice(
            (selectedMenteeTablePage - 1) * MENTEE_PAGE_SIZE,
            selectedMenteeTablePage * MENTEE_PAGE_SIZE,
        ),
    );

    $effect(() => {
        if (selectedMenteeTablePage > selectedMenteePageCount) {
            selectedMenteeTablePage = selectedMenteePageCount;
        }
    });

    $effect(() => {
        menteeUuids.length;
        selectedMenteeTablePage = 1;
    });

    let selectedSharerRows = $derived.by(() => {
        const rows = sharerUserUuids
            .map((uuid) => users.find((u) => u.uuid === uuid))
            .filter((u): u is User => u != null);
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
    let selectedSharerPageCount = $derived(
        Math.max(1, Math.ceil(selectedSharerRows.length / MENTEE_PAGE_SIZE)),
    );
    let selectedSharerPageRows = $derived(
        selectedSharerRows.slice(
            (selectedSharerTablePage - 1) * MENTEE_PAGE_SIZE,
            selectedSharerTablePage * MENTEE_PAGE_SIZE,
        ),
    );

    $effect(() => {
        if (selectedSharerTablePage > selectedSharerPageCount) {
            selectedSharerTablePage = selectedSharerPageCount;
        }
    });

    $effect(() => {
        sharerUserUuids.length;
        selectedSharerTablePage = 1;
    });

    let submitting = $state(false);
    let error = $state("");

    let form = $state({
        templateUuid: "",
        implementedAt: todayISO(),
        fromMonth: currentMonthISO(),
        toMonth: currentMonthISO(),
        week: 2,
        weekday: 3,
    });

    let managedTeamPickValue = $state("");

    let deptMaps = $derived(buildUserDepartmentMaps(users));
    let teamMaps = $derived(buildUserTeamMaps(users));
    let pairTableRows = $derived(expandUserPairings(pairings));

    let selectedPairRows = $derived.by(() => {
        const rows = selectedPairingUuids
            .map((uuid) => pairTableRows.find((r) => r.pairing.uuid === uuid))
            .filter((r): r is NonNullable<typeof r> => r != null);
        rows.sort((a, b) => {
            const da = departmentNameForUser(a.menteeUserUuid, deptMaps) || "";
            const db = departmentNameForUser(b.menteeUserUuid, deptMaps) || "";
            const c = da.localeCompare(db, "ja");
            if (c !== 0) return c;
            return a.menteeName.localeCompare(b.menteeName, "ja");
        });
        return rows;
    });
    let selectedPairPageCount = $derived(
        Math.max(1, Math.ceil(selectedPairRows.length / MENTEE_PAGE_SIZE)),
    );
    let selectedPairPageRows = $derived(
        selectedPairRows.slice(
            (selectedPairTablePage - 1) * MENTEE_PAGE_SIZE,
            selectedPairTablePage * MENTEE_PAGE_SIZE,
        ),
    );

    $effect(() => {
        if (selectedPairTablePage > selectedPairPageCount) {
            selectedPairTablePage = selectedPairPageCount;
        }
    });

    $effect(() => {
        selectedPairingUuids.length;
        selectedPairTablePage = 1;
    });

    function avatarURL(user: User | undefined) {
        if (!user?.imageUri) return "";
        return `/uploads/${user.imageUri}`;
    }

    let selectedMentorUser = $derived(
        mentorUuid ? users.find((u) => u.uuid === mentorUuid) : undefined,
    );

    function templateLabelOf(item: Template) {
        return item.title || "";
    }
    function managedTeamLabelOf(item: ManagedTeam) {
        return `${item.name}（${item.members?.length ?? 0}人）`;
    }
    function pickManagedTeamFromSelect(item: ManagedTeam) {
        pickManagedTeam(item);
        managedTeamModalOpen = false;
    }
    function pickManagedTeam(team: ManagedTeam) {
        const next = new Set(menteeUuids);
        for (const m of team.members || []) {
            if (!m?.uuid) continue;
            if (currentUserUuid && m.uuid === currentUserUuid) continue;
            next.add(m.uuid);
        }
        menteeUuids = [...next];
    }
    function submitLabel() {
        if (submitting) return "作成中…";
        return "作成";
    }

    function monthlyPayloadBase() {
        return {
            templateUuid: form.templateUuid,
            from: `${form.fromMonth}-01`,
            to: `${form.toMonth}-01`,
            week: form.week,
            weekday: form.weekday,
        };
    }

    function validateMonthlySchedule(): boolean {
        if (!form.fromMonth || !form.toMonth) {
            error = "期間を入力してください。";
            return false;
        }
        if (form.fromMonth > form.toMonth) {
            error = "開始月は終了月以前にしてください。";
            return false;
        }
        if (form.week < 1 || form.week > 5) {
            error = "第N週は 1〜5 で指定してください。";
            return false;
        }
        return true;
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

    async function onTargetChange(next: CreateTarget) {
        target = next;
        if (next === "pair") {
            selectedPairingUuids = [];
            await loadPairingsIfNeeded();
        }
    }

    function onScheduleChange(next: ScheduleKind) {
        schedule = next;
    }

    onMount(async () => {
        const [t, dir, profile] = await Promise.all([
            api<Template[]>("GET", "/api/v1/templates"),
            api<{
                users?: User[];
                managedTeams?: ManagedTeam[];
                departments?: { uuid: string; name: string }[];
                teams?: DirectoryTeam[];
            }>("GET", "/api/v1/planning/directory"),
            api<{ email?: string }>("GET", "/api/v1/setting/profile"),
        ]);
        templates = t;
        users = dir.users || [];
        departments = dir.departments || [];
        teams = (dir.teams || []).map((t) => ({
            uuid: t.uuid,
            name: t.name,
            departmentUuid: t.departmentUuid,
        }));
        managedTeams = dir.managedTeams || [];
        currentUserUuid =
            users.find((u) => u.email === profile.email)?.uuid || "";
        canSpecifyPair = getAvatar()?.canSpecifyPair ?? false;
        const defaultTpl = templates.find((x) => x.isDefault) ?? templates[0];
        form.templateUuid = defaultTpl?.uuid || "";
    });

    async function openPairModal() {
        await loadPairingsIfNeeded();
        pairModalOpen = true;
    }

    async function submit() {
        error = "";
        if (!form.templateUuid) {
            error = "テンプレートを選択してください。";
            return;
        }
        submitting = true;
        try {
            const sharerPayload = sharerUserUuids.length
                ? sharerUserUuids
                : undefined;

            if (target === "manual" && schedule === "date") {
                if (!mentorUuid) {
                    error = "メンターを選択してください。";
                    return;
                }
                if (!menteeUuids.length) {
                    error = "メンティを選択してください。";
                    return;
                }
                if (!form.implementedAt) {
                    error = "実施日を入力してください。";
                    return;
                }
                if (menteeUuids.some((id) => id === mentorUuid)) {
                    error = "メンターと同じユーザをメンティに指定できません。";
                    return;
                }
                await api("POST", "/api/v1/meetings", {
                    mentorUuid,
                    menteeUuids,
                    templateUuid: form.templateUuid,
                    implementedAt: form.implementedAt,
                    sharerUserUuids: sharerPayload,
                });
            } else if (target === "manual" && schedule === "monthly") {
                if (!mentorUuid) {
                    error = "メンターを選択してください。";
                    return;
                }
                if (!menteeUuids.length) {
                    error = "メンティを選択してください。";
                    return;
                }
                if (menteeUuids.some((id) => id === mentorUuid)) {
                    error = "メンターと同じユーザをメンティに指定できません。";
                    return;
                }
                if (!validateMonthlySchedule()) return;
                await api("POST", "/api/v1/meetings/monthly", {
                    mentorUuid,
                    menteeUuids,
                    ...monthlyPayloadBase(),
                    sharerUserUuids: sharerPayload,
                });
            } else if (target === "pair" && schedule === "date") {
                await loadPairingsIfNeeded();
                if (!pairTableRows.length) {
                    error = "ペアがありません。";
                    return;
                }
                if (!selectedPairingUuids.length) {
                    error = "ペアを選択してください。";
                    return;
                }
                if (!form.implementedAt) {
                    error = "実施日を入力してください。";
                    return;
                }
                const pairingUuids = [...new Set(selectedPairingUuids)];
                await api("POST", "/api/v1/meetings/all", {
                    templateUuid: form.templateUuid,
                    implementedAt: form.implementedAt,
                    sharerUserUuids: sharerPayload,
                    pairingUuids,
                });
            } else {
                await loadPairingsIfNeeded();
                if (!pairTableRows.length) {
                    error = "ペアがありません。";
                    return;
                }
                if (!selectedPairingUuids.length) {
                    error = "ペアを選択してください。";
                    return;
                }
                if (!validateMonthlySchedule()) return;
                const pairingUuids = [...new Set(selectedPairingUuids)];
                await api("POST", "/api/v1/meetings/monthly/pairings", {
                    ...monthlyPayloadBase(),
                    sharerUserUuids: sharerPayload,
                    pairingUuids,
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
            class:btn-primary={target === "manual"}
            class:btn-outline={target !== "manual"}
            onclick={() => onTargetChange("manual")}>作成</button
        >
        {#if canSpecifyPair}
            <button
                type="button"
                class="btn btn-sm"
                class:btn-primary={target === "pair"}
                class:btn-outline={target !== "pair"}
                onclick={() => onTargetChange("pair")}>ペア指定</button
            >
        {/if}
    </div>
    <div
        class="mt-6 space-y-4 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm"
    >
        {#if error}
            <div class="alert alert-error text-sm">{error}</div>
        {/if}
        <div class="form-control">
            <span class="label-text font-medium">テンプレート</span>
            <div class="mt-1">
                <SearchableSelect
                    pickRootId="meeting-new-template"
                    items={templates}
                    bind:value={form.templateUuid}
                    labelOf={templateLabelOf}
                    placeholder="名前で検索して選択"
                    allowClear={true}
                />
            </div>
        </div>
        {#if target === "pair"}
            <div class="form-control space-y-3">
                <div class="flex flex-1 w-full justify-between">
                    <span class="label-text font-medium">ペアを選択</span>
                    <button
                        type="button"
                        class="btn btn-outline btn-xs mt-1 w-fit"
                        disabled={pairingsLoading}
                        onclick={openPairModal}
                    >
                        {#if pairingsLoading}読み込み中…{:else}ペアを選択{/if}
                    </button>
                </div>
                {#if pairingsLoading && !pairingsLoaded}
                    <div class="space-y-3">
                        {#each Array(3) as _, i (i)}
                            <div class="skeleton h-10 w-full rounded-xl"></div>
                        {/each}
                    </div>
                {:else if pairingsLoaded && pairTableRows.length === 0}
                    <p class="text-sm text-base-content/50">ペアがありません</p>
                {:else if selectedPairRows.length === 0}
                    <p class="text-sm text-base-content/50">未選択</p>
                {:else}
                    <div
                        class="overflow-x-auto overflow-hidden rounded-2xl border border-base-300"
                    >
                        <table class="table table-sm">
                            <thead>
                                <tr>
                                    <th>メンター</th>
                                    <th>部署</th>
                                    <th>メンティ</th>
                                    <th>部署</th>
                                    <th>チーム</th>
                                </tr>
                            </thead>
                            <tbody>
                                {#each selectedPairPageRows as row (row.rowKey)}
                                    {@const mentorUser = users.find(
                                        (u) => u.uuid === row.mentorUserUuid,
                                    )}
                                    {@const menteeUser = users.find(
                                        (u) => u.uuid === row.menteeUserUuid,
                                    )}
                                    <tr>
                                        <td>
                                            <span
                                                class="flex items-center gap-2"
                                            >
                                                <Icon
                                                    url={avatarURL(mentorUser)}
                                                    text={row.mentorName}
                                                    width={28}
                                                    height={28}
                                                />
                                                {row.mentorName}
                                            </span>
                                        </td>
                                        <td
                                            >{departmentNameForUser(
                                                row.mentorUserUuid,
                                                deptMaps,
                                            ) || "—"}</td
                                        >
                                        <td>
                                            <span
                                                class="flex items-center gap-2"
                                            >
                                                <Icon
                                                    url={avatarURL(menteeUser)}
                                                    text={row.menteeName}
                                                    width={28}
                                                    height={28}
                                                />
                                                {row.menteeName}
                                            </span>
                                        </td>
                                        <td
                                            >{departmentNameForUser(
                                                row.menteeUserUuid,
                                                deptMaps,
                                            ) || "—"}</td
                                        >
                                        <td
                                            >{teamNameForUser(
                                                row.menteeUserUuid,
                                                teamMaps,
                                            ) || "—"}</td
                                        >
                                    </tr>
                                {/each}
                            </tbody>
                        </table>
                    </div>
                    {#if selectedPairRows.length > MENTEE_PAGE_SIZE}
                        <div class="join justify-end w-full flex">
                            <button
                                type="button"
                                class="btn btn-sm join-item"
                                disabled={selectedPairTablePage <= 1}
                                onclick={() =>
                                    (selectedPairTablePage = Math.max(
                                        1,
                                        selectedPairTablePage - 1,
                                    ))}>前へ</button
                            >
                            <span
                                class="btn btn-sm join-item btn-disabled pointer-events-none"
                                >{selectedPairTablePage} /
                                {selectedPairPageCount}</span
                            >
                            <button
                                type="button"
                                class="btn btn-sm join-item"
                                disabled={selectedPairTablePage >=
                                    selectedPairPageCount}
                                onclick={() =>
                                    (selectedPairTablePage = Math.min(
                                        selectedPairPageCount,
                                        selectedPairTablePage + 1,
                                    ))}>次へ</button
                            >
                        </div>
                    {/if}
                    <p class="text-xs text-base-content/60">
                        選択済み: {selectedPairRows.length} ペア
                    </p>
                {/if}
            </div>
        {:else}
            <div class="grid items-end gap-4 grid-cols-1">
                <div class="form-control space-y-3">
                    <div class="flex flex-1 justify-between w-full items-end">
                        <span class="label-text font-medium">メンター</span>
                        <button
                            type="button"
                            class="btn btn-outline btn-xs mt-1 w-fit"
                            onclick={() => (mentorModalOpen = true)}
                            >メンターを選択</button
                        >
                    </div>
                    {#if !selectedMentorUser}
                        <p class="text-sm text-base-content/50">未選択</p>
                    {:else}
                        <div
                            class="overflow-hidden rounded-2xl border border-base-300"
                        >
                            <table class="table table-sm">
                                <thead>
                                    <tr>
                                        <th>部署</th>
                                        <th>チーム</th>
                                        <th>氏名</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr>
                                        <td>
                                            {selectedMentorUser.departmentNames ||
                                                "所属無し"}
                                        </td>
                                        <td>
                                            {selectedMentorUser.teamNames || ""}
                                        </td>
                                        <td>
                                            <span
                                                class="flex items-center gap-2"
                                            >
                                                <Icon
                                                    url={avatarURL(
                                                        selectedMentorUser,
                                                    )}
                                                    text={selectedMentorUser.name}
                                                    width={28}
                                                    height={28}
                                                />
                                                {selectedMentorUser.name}
                                            </span>
                                        </td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    {/if}
                </div>
                <div class="space-y-3">
                    <div class="form-control">
                        <div class="flex items-center justify-between gap-2">
                            <div class="flex">
                                <span class="label-text font-medium"
                                    >メンティ</span
                                >
                            </div>
                            <div class="flex gap-2 items-end">
                                <button
                                    type="button"
                                    class="btn btn-outline btn-xs mt-1 w-fit"
                                    onclick={() => (menteeModalOpen = true)}
                                    >メンティを選択</button
                                >
                                {#if managedTeams.length > 0}
                                    <button
                                        type="button"
                                        class="btn btn-outline btn-xs"
                                        onclick={() =>
                                            (managedTeamModalOpen = true)}
                                        >管理チーム選択</button
                                    >
                                {/if}
                            </div>
                        </div>
                    </div>
                </div>
                <div class="flex flex-1 w-full">
                    {#if selectedMenteeRows.length === 0}
                        <p class="text-sm text-base-content/50">未選択</p>
                    {:else}
                        <div
                            class="flex flex-1 w-full overflow-hidden rounded-2xl border border-base-300"
                        >
                            <table class="table table-sm">
                                <thead>
                                    <tr>
                                        <th>部署</th>
                                        <th>チーム</th>
                                        <th>氏名</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {#each selectedMenteePageRows as u (u.uuid)}
                                        <tr>
                                            <td>
                                                {u.departmentNames ||
                                                    "所属無し"}
                                            </td>
                                            <td>
                                                {u.teamNames || ""}
                                            </td>
                                            <td>
                                                <span
                                                    class="flex items-center gap-2"
                                                >
                                                    <Icon
                                                        url={avatarURL(u)}
                                                        text={u.name}
                                                        width={28}
                                                        height={28}
                                                    />
                                                    {u.name}
                                                </span>
                                            </td>
                                        </tr>
                                    {/each}
                                </tbody>
                            </table>
                        </div>
                        {#if selectedMenteeRows.length > MENTEE_PAGE_SIZE}
                            <div class="join justify-end w-full flex">
                                <button
                                    type="button"
                                    class="btn btn-sm join-item"
                                    disabled={selectedMenteeTablePage <= 1}
                                    onclick={() =>
                                        (selectedMenteeTablePage = Math.max(
                                            1,
                                            selectedMenteeTablePage - 1,
                                        ))}>前へ</button
                                >
                                <span
                                    class="btn btn-sm join-item btn-disabled pointer-events-none"
                                    >{selectedMenteeTablePage} /
                                    {selectedMenteePageCount}</span
                                >
                                <button
                                    type="button"
                                    class="btn btn-sm join-item"
                                    disabled={selectedMenteeTablePage >=
                                        selectedMenteePageCount}
                                    onclick={() =>
                                        (selectedMenteeTablePage = Math.min(
                                            selectedMenteePageCount,
                                            selectedMenteeTablePage + 1,
                                        ))}>次へ</button
                                >
                            </div>
                        {/if}
                    {/if}
                </div>
            </div>
        {/if}
        <div class="form-control space-y-3">
            <div class="flex flex-1 justify-between w-full items-end">
                <span class="label-text font-medium">共有者（任意）</span>
                <button
                    type="button"
                    class="btn btn-outline btn-xs mt-1 w-fit"
                    onclick={() => (sharerModalOpen = true)}
                    >共有者を選択</button
                >
            </div>
            {#if selectedSharerRows.length === 0}
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
                            {#each selectedSharerPageRows as u (u.uuid)}
                                <tr>
                                    <td>{u.departmentNames || "所属無し"}</td>
                                    <td>
                                        <span class="flex items-center gap-2">
                                            <Icon
                                                url={avatarURL(u)}
                                                text={u.name}
                                                width={28}
                                                height={28}
                                            />
                                            {u.name}
                                        </span>
                                    </td>
                                </tr>
                            {/each}
                        </tbody>
                    </table>
                </div>
                {#if selectedSharerRows.length > MENTEE_PAGE_SIZE}
                    <div class="join justify-end w-full flex">
                        <button
                            type="button"
                            class="btn btn-sm join-item"
                            disabled={selectedSharerTablePage <= 1}
                            onclick={() =>
                                (selectedSharerTablePage = Math.max(
                                    1,
                                    selectedSharerTablePage - 1,
                                ))}>前へ</button
                        >
                        <span
                            class="btn btn-sm join-item btn-disabled pointer-events-none"
                            >{selectedSharerTablePage} /
                            {selectedSharerPageCount}</span
                        >
                        <button
                            type="button"
                            class="btn btn-sm join-item"
                            disabled={selectedSharerTablePage >=
                                selectedSharerPageCount}
                            onclick={() =>
                                (selectedSharerTablePage = Math.min(
                                    selectedSharerPageCount,
                                    selectedSharerTablePage + 1,
                                ))}>次へ</button
                        >
                    </div>
                {/if}
                <p class="text-xs text-base-content/60">
                    選択済み: {selectedSharerRows.length} 人
                </p>
            {/if}
        </div>
        <section
            class="rounded-2xl border border-base-300 bg-base-200/30 p-4 space-y-4"
        >
            <div
                class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
            >
                <span class="label-text font-medium">日程</span>
                <div class="join">
                    <button
                        type="button"
                        class="btn btn-sm join-item"
                        class:btn-primary={schedule === "date"}
                        class:btn-outline={schedule !== "date"}
                        onclick={() => onScheduleChange("date")}
                        >日付指定</button
                    >
                    <button
                        type="button"
                        class="btn btn-sm join-item"
                        class:btn-primary={schedule === "monthly"}
                        class:btn-outline={schedule !== "monthly"}
                        onclick={() => onScheduleChange("monthly")}
                        >繰り返し</button
                    >
                </div>
            </div>

            {#if schedule === "date"}
                <label class="form-control w-full">
                    <span class="label-text font-medium">実施日</span>
                    <input
                        class="input input-bordered mt-1"
                        type="date"
                        bind:value={form.implementedAt}
                    />
                </label>
            {:else}
                <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
                    <label class="form-control">
                        <span class="label-text font-medium">年月（から）</span>
                        <input
                            class="input input-bordered mt-1"
                            type="month"
                            bind:value={form.fromMonth}
                        />
                    </label>
                    <label class="form-control">
                        <span class="label-text font-medium">年月（まで）</span>
                        <input
                            class="input input-bordered mt-1"
                            type="month"
                            bind:value={form.toMonth}
                        />
                    </label>
                    <label class="form-control">
                        <span class="label-text font-medium">実施週</span>
                        <select
                            class="select select-bordered mt-1"
                            bind:value={form.week}
                        >
                            <option value={1}>第1週</option>
                            <option value={2}>第2週</option>
                            <option value={3}>第3週</option>
                            <option value={4}>第4週</option>
                            <option value={5}>第5週</option>
                        </select>
                    </label>
                    <label class="form-control">
                        <span class="label-text font-medium">曜日</span>
                        <select
                            class="select select-bordered mt-1"
                            bind:value={form.weekday}
                        >
                            <option value={1}>月曜</option>
                            <option value={2}>火曜</option>
                            <option value={3}>水曜</option>
                            <option value={4}>木曜</option>
                            <option value={5}>金曜</option>
                        </select>
                    </label>
                </div>
            {/if}
        </section>
        <div class="flex flex-wrap justify-end gap-2 pt-2">
            <button
                type="button"
                class="btn btn-primary"
                disabled={submitting}
                onclick={submit}>{submitLabel()}</button
            >
            <a href="/planning/meetings" class="btn btn-ghost　btn-outline"
                >戻る</a
            >
        </div>
    </div>
</main>

<UserSelectModal
    open={mentorModalOpen}
    title="メンターを選択"
    multiple={false}
    {users}
    {departments}
    {teams}
    selectedUuids={mentorUuid ? [mentorUuid] : []}
    pickRootId="meeting-new-mentor-name"
    onApply={(uuids) => {
        mentorUuid = uuids[0] ?? "";
        mentorModalOpen = false;
    }}
    onClose={() => {
        mentorModalOpen = false;
    }}
/>

<UserSelectModal
    open={menteeModalOpen}
    title="メンティを選択"
    users={pickableUsers}
    {departments}
    {teams}
    selectedUuids={menteeUuids}
    pickRootId="meeting-new-mentee-name"
    onApply={(uuids) => {
        menteeUuids = uuids;
        menteeModalOpen = false;
    }}
    onClose={() => {
        menteeModalOpen = false;
    }}
/>

<UserSelectModal
    open={sharerModalOpen}
    title="共有者を選択"
    users={target === "pair" ? users : pickableUsers}
    {departments}
    {teams}
    selectedUuids={sharerUserUuids}
    pickRootId="meeting-new-sharer-name"
    onApply={(uuids) => {
        sharerUserUuids = uuids;
        sharerModalOpen = false;
    }}
    onClose={() => {
        sharerModalOpen = false;
    }}
/>

<PairSelectModal
    open={pairModalOpen}
    rows={pairTableRows}
    {users}
    {departments}
    {teams}
    selectedUuids={selectedPairingUuids}
    pickRootId="meeting-new-pair-name"
    onApply={(uuids) => {
        selectedPairingUuids = uuids;
        pairModalOpen = false;
    }}
    onClose={() => {
        pairModalOpen = false;
    }}
/>

<dialog class="modal" class:modal-open={managedTeamModalOpen}>
    <div class="modal-box max-w-lg">
        <h3 class="text-lg font-bold">管理チーム選択</h3>
        <p class="mt-1 text-sm text-base-content/60">
            チームを選ぶと、そのメンバーがメンティに追加されます。
        </p>
        <div class="mt-4">
            <SearchableSelect
                pickRootId="meeting-new-managed-team"
                items={managedTeams}
                bind:value={managedTeamPickValue}
                labelOf={managedTeamLabelOf}
                placeholder="チーム名で検索"
                clearOnSelect={true}
                onSelect={pickManagedTeamFromSelect}
            />
        </div>
        <div class="modal-action">
            <button
                type="button"
                class="btn"
                onclick={() => (managedTeamModalOpen = false)}>閉じる</button
            >
        </div>
    </div>
    <form method="dialog" class="modal-backdrop">
        <button
            type="button"
            aria-label="閉じる"
            onclick={() => (managedTeamModalOpen = false)}
        ></button>
    </form>
</dialog>
