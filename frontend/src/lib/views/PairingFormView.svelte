<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { api } from "$lib/api";
    import { showToast } from "$lib/toast";
    import UserChipSearchField from "$lib/components/UserChipSearchField.svelte";
    import SearchableSelect from "$lib/components/SearchableSelect.svelte";
    import Icon from "$lib/components/Icon.svelte";

    type DirectoryUser = {
        uuid: string;
        name: string;
        email: string;
        roleName?: string;
        departmentNames?: string;
        teamNames?: string;
        departmentUuid?: string;
        imageUri?: string;
    };
    type DirectoryTeam = {
        uuid: string;
        name: string;
        departmentUuid: string;
    };
    type Target = { kind: string; targetUuid: string; label: string };
    type PairingInitial = {
        uuid?: string;
        mentor?: { uuid?: string };
        mentee?: { uuid?: string; name?: string };
    };

    let {
        mode,
        initialPairing = null,
    }: {
        mode: "create" | "edit";
        initialPairing?: PairingInitial | null;
    } = $props();

    let mentorUuid = $state("");
    let targets = $state<Target[]>([]);
    let pairingUUID = $state("");
    let users = $state<DirectoryUser[]>([]);
    let departments = $state<{ uuid: string; name: string }[]>([]);
    let teams = $state<DirectoryTeam[]>([]);

    let nameSearch = $state("");
    let filterDeptUuid = $state("");
    let filterTeamUuid = $state("");
    let sortKey = $state<"role" | "name" | "department" | "team">("name");
    let sortDir = $state<"asc" | "desc">("asc");
    let selectedMenteeUuids = $state<Set<string>>(new Set());

    let userTargets = $derived(targets.filter((t) => t.kind === "user"));
    let menteeChips = $derived(
        userTargets.map((t) => ({
            uuid: t.targetUuid,
            name: displayNameForTarget(t),
        })),
    );

    function avatarURL(user: DirectoryUser) {
        if (!user.imageUri) return "";
        return `/uploads/${user.imageUri}`;
    }
    function displayNameForTarget(t: Target) {
        const u = users.find((x) => x.uuid === t.targetUuid);
        if (u) return u.name;
        const idx = t.label.indexOf(" (");
        if (idx > 0) return t.label.slice(0, idx);
        return t.label || t.targetUuid;
    }
    function userLabel(uuid: string) {
        const u = users.find((x) => x.uuid === uuid);
        return u ? u.name : uuid;
    }
    function userLabelOf(item: DirectoryUser) {
        return `${item.name} ${item.email}`;
    }
    function targetKey(kind: string, targetUuid: string) {
        return `${kind}:${targetUuid}`;
    }
    function addUserTarget(u: DirectoryUser) {
        const key = targetKey("user", u.uuid);
        if (targets.some((t) => targetKey(t.kind, t.targetUuid) === key))
            return;
        targets = [
            ...targets,
            {
                kind: "user",
                targetUuid: u.uuid,
                label: `${u.name} (${u.email})`,
            },
        ];
    }

    let teamOptions = $derived.by(() => {
        if (!filterDeptUuid) return teams;
        return teams.filter((t) => t.departmentUuid === filterDeptUuid);
    });

    $effect(() => {
        if (!filterTeamUuid) return;
        const ok = teamOptions.some((t) => t.uuid === filterTeamUuid);
        if (!ok) filterTeamUuid = "";
    });

    let tableRows = $derived.by(() => {
        let list = users.filter((u) => u.uuid !== mentorUuid);

        const q = nameSearch.trim().toLowerCase();
        if (q) list = list.filter((u) => u.name.toLowerCase().includes(q));

        if (filterDeptUuid) {
            const department = departments.find(
                (d) => d.uuid === filterDeptUuid,
            );
            const departmentName = department?.name || "";
            list = list.filter((u) => u.departmentNames === departmentName);
        }
        if (filterTeamUuid) {
            const team = teams.find((t) => t.uuid === filterTeamUuid);
            const teamName = team?.name || "";
            list = list.filter((u) => (u.teamNames || "").includes(teamName));
        }
        const dir = sortDir === "asc" ? 1 : -1;
        list.sort((a, b) => {
            let av = "";
            let bv = "";
            switch (sortKey) {
                case "role":
                    av = a.roleName || "";
                    bv = b.roleName || "";
                    break;
                case "department":
                    av = a.departmentNames || "";
                    bv = b.departmentNames || "";
                    break;
                case "team":
                    av = a.teamNames || "";
                    bv = b.teamNames || "";
                    break;
                default:
                    av = a.name;
                    bv = b.name;
            }
            return dir * av.localeCompare(bv, "ja");
        });
        return list;
    });

    function toggleMentee(uuid: string) {
        const next = new Set(selectedMenteeUuids);
        if (next.has(uuid)) next.delete(uuid);
        else next.add(uuid);
        selectedMenteeUuids = next;
    }

    function switchMentee(uuid: string, flag: boolean) {
        const next = new Set(selectedMenteeUuids);
        if (!flag) next.delete(uuid);
        else next.add(uuid);
        selectedMenteeUuids = next;
    }

    function toggleSort(key: "role" | "name" | "department" | "team") {
        if (sortKey === key) sortDir = sortDir === "asc" ? "desc" : "asc";
        else {
            sortKey = key;
            sortDir = "asc";
        }
    }

    function toggleAllUsers(ev: Event) {
        let target = ev.currentTarget;
        if (!(target instanceof HTMLInputElement)) return;
        tableRows.forEach((t) => {
            switchMentee(t.uuid, target.checked);
        });
    }

    function sortMark(key: string) {
        if (sortKey !== key) return "";
        return sortDir === "asc" ? " ↑" : " ↓";
    }

    $effect(() => {
        if (!initialPairing) return;
        pairingUUID = initialPairing.uuid || "";
        mentorUuid = initialPairing.mentor?.uuid || "";
        const loaded: Target[] = [];
        if (initialPairing.mentee?.uuid) {
            loaded.push({
                kind: "user",
                targetUuid: initialPairing.mentee.uuid,
                label: initialPairing.mentee.name || "",
            });
        }
        targets = loaded;
    });

    onMount(async () => {
        const data = await api<{
            users?: DirectoryUser[];
            departments?: { uuid: string; name: string }[];
            teams?: DirectoryTeam[];
        }>("GET", "/api/v1/planning/directory");
        users = data.users || [];
        departments = data.departments || [];
        teams = (data.teams || []).map((t) => ({
            uuid: t.uuid,
            name: t.name,
            departmentUuid: t.departmentUuid,
        }));
    });

    function pickUserTarget(u: DirectoryUser) {
        addUserTarget(u);
    }
    function removeUserTarget(targetUuid: string) {
        targets = targets.filter(
            (t) => !(t.kind === "user" && t.targetUuid === targetUuid),
        );
    }

    async function requestSave() {
        let userOnly: { kind: string; targetUuid: string }[];
        if (mode === "create") {
            userOnly = [...selectedMenteeUuids].map((uuid) => ({
                kind: "user",
                targetUuid: uuid,
            }));
        } else {
            userOnly = targets
                .filter((t) => t.kind === "user")
                .map((t) => ({ kind: "user", targetUuid: t.targetUuid }));
        }
        if (!mentorUuid || !userOnly.length) {
            showToast("メンターとメンティを選択してください", "error");
            return;
        }
        const body = {
            mentorUuid,
            targets: userOnly,
        };
        try {
            if (mode === "create") {
                await api("POST", "/api/v1/pairings", body);
            } else {
                await api("PUT", `/api/v1/pairings/${pairingUUID}`, body);
            }
            showToast("保存しました");
            await goto("/planning/pairs");
        } catch {
            showToast("保存に失敗しました", "error");
        }
    }
</script>

<svelte:head>
    <title>{mode === "create" ? "ペア作成" : "ペア編集"} - Growlink</title>
</svelte:head>

<main class="mx-auto w-full max-w-5xl min-w-0 p-4 md:p-8">
    <h1 class="font-display text-2xl font-semibold">
        {mode === "create" ? "ペア新規作成" : "ペア編集"}
    </h1>
    <p class="mt-1 text-sm text-base-content/60">
        メンターとメンティを{mode === "create" ? "登録" : "変更"}します
    </p>
    <div
        class="mt-6 space-y-4 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm">
        <div class="form-control max-w-md">
            <span class="label-text font-medium">メンター</span>
            <SearchableSelect
                pickRootId="pairing-form-mentor"
                items={users}
                bind:value={mentorUuid}
                labelOf={userLabelOf}
                selectedLabel={userLabel}
                placeholder="名前で検索" />
        </div>

        {#if mode === "create"}
            <div class="flex flex-wrap justify-between gap-3 items-end">
                <span class="label-text font-medium w-1/2 min-w-48"
                    >メンティー
                    <label
                        class="input input-bordered flex items-center w-full">
                        <input
                            type="search"
                            class="grow"
                            placeholder="氏名で検索"
                            bind:value={nameSearch} />
                    </label>
                </span>
                <div class="flex gap-2 justify-end">
                    <select
                        class="select select-bordered w-full"
                        bind:value={filterDeptUuid}>
                        <option value="">部署: すべて</option>
                        {#each departments as d (d.uuid)}
                            <option value={d.uuid}>{d.name}</option>
                        {/each}
                    </select>
                    <select
                        class="select select-bordered w-full"
                        bind:value={filterTeamUuid}>
                        <option value="">チーム: すべて</option>
                        {#each teamOptions as t (t.uuid)}
                            <option value={t.uuid}>{t.name}</option>
                        {/each}
                    </select>
                </div>
            </div>
            <div class="overflow-x-auto rounded-2xl border border-base-300">
                <table class="table table-sm">
                    <thead>
                        <tr>
                            <th>
                                <input
                                    class="checkbox"
                                    type="checkbox"
                                    onclick={(e: Event) => toggleAllUsers(e)} />
                            </th>
                            <th>
                                <button
                                    type="button"
                                    class="font-semibold"
                                    onclick={() => toggleSort("department")}
                                    >部署{sortMark("department")}</button>
                            </th>
                            <th>
                                <button
                                    type="button"
                                    class="font-semibold"
                                    onclick={() => toggleSort("role")}
                                    >役職{sortMark("role")}</button>
                            </th>
                            <th>
                                <button
                                    type="button"
                                    class="font-semibold"
                                    onclick={() => toggleSort("name")}
                                    >氏名{sortMark("name")}</button>
                            </th>
                            <th>
                                <button
                                    type="button"
                                    class="font-semibold"
                                    onclick={() => toggleSort("team")}
                                    >チーム{sortMark("team")}</button>
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        {#each tableRows as u (u.uuid)}
                            <tr
                                class="hover:bg-base-200/40 cursor-pointer"
                                onclick={() => toggleMentee(u.uuid)}>
                                <td>
                                    <input
                                        type="checkbox"
                                        class="checkbox checkbox-sm"
                                        checked={selectedMenteeUuids.has(
                                            u.uuid,
                                        )} />
                                </td>
                                <td>{u.departmentNames || "所属無し"}</td>
                                <td>{u.roleName || ""}</td>
                                <td>
                                    <span class="flex gap-2 items-center">
                                        <Icon
                                            url={avatarURL(u)}
                                            text={u.name ?? ""}
                                            width={28}
                                            height={28} />
                                        {u.name}
                                    </span>
                                </td>
                                <td>{u.teamNames || ""}</td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
            <p class="text-xs text-base-content/60">
                選択中: {selectedMenteeUuids.size} 人
            </p>
        {:else}
            <div class="form-control">
                <span class="label-text font-medium">メンティ</span>
                <UserChipSearchField
                    pickRootId="pairing-form-mentee"
                    chips={menteeChips}
                    items={users}
                    labelOf={userLabelOf}
                    placeholder={menteeChips.length
                        ? "追加で検索"
                        : "名前で検索して追加"}
                    onPick={pickUserTarget}
                    onRemove={removeUserTarget} />
            </div>
        {/if}
        <div class="flex flex-wrap gap-2 pt-2 justify-end">
            <button type="button" class="btn btn-primary" onclick={requestSave}
                >{mode === "create" ? "作成" : "保存"}</button>
            <a href="/planning/pairs" class="btn btn-ghost">戻る</a>
        </div>
    </div>
</main>
