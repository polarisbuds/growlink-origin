<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { api } from "$lib/api";
    import { showToast } from "$lib/toast";
    import UserSelectModal from "$lib/components/UserSelectModal.svelte";
    import Icon from "$lib/components/Icon.svelte";

    const MENTEE_PAGE_SIZE = 10;

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
    let mentorUuid = $state("");
    let users = $state<DirectoryUser[]>([]);
    let departments = $state<{ uuid: string; name: string }[]>([]);
    let teams = $state<DirectoryTeam[]>([]);
    let mentorModalOpen = $state(false);
    let menteeModalOpen = $state(false);
    let menteeUuids = $state<string[]>([]);
    let selectedMenteeTablePage = $state(1);

    let pickableUsers = $derived(
        mentorUuid ? users.filter((u) => u.uuid !== mentorUuid) : users,
    );

    let selectedMenteeRows = $derived.by(() => {
        const rows = menteeUuids
            .map((uuid) => users.find((u) => u.uuid === uuid))
            .filter((u): u is DirectoryUser => u != null);
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

    function avatarURL(user: DirectoryUser | undefined) {
        if (!user?.imageUri) return "";
        return `/uploads/${user.imageUri}`;
    }
    let selectedMentorUser = $derived(
        mentorUuid ? users.find((u) => u.uuid === mentorUuid) : undefined,
    );

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

    async function requestSave() {
        const userOnly = menteeUuids.map((uuid) => ({
            kind: "user",
            targetUuid: uuid,
        }));
        if (!mentorUuid || !userOnly.length) {
            showToast("メンターとメンティを選択してください", "error");
            return;
        }
        if (menteeUuids.some((id) => id === mentorUuid)) {
            showToast(
                "メンターと同じユーザをメンティに指定できません",
                "error",
            );
            return;
        }
        const body = {
            mentorUuid,
            targets: userOnly,
        };
        try {
            await api("POST", "/api/v1/pairings", body);
            showToast("作成しました");
            await goto("/planning/pairs");
        } catch {
            showToast("保存に失敗しました", "error");
        }
    }
</script>

<svelte:head>
    <title>ペア作成 - Growlink</title>
</svelte:head>

<main class="mx-auto w-full max-w-5xl min-w-0 p-4 md:p-8">
    <h1 class="font-display text-2xl font-semibold">ペア新規作成</h1>
    <p class="mt-1 text-sm text-base-content/60">
        メンターとメンティを登録します
    </p>
    <div
        class="mt-6 space-y-4 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm"
    >
        <div class="form-control space-y-3 w-full">
            <div class="flex flex-1 w-full justify-between items-end">
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
                <div class="overflow-hidden rounded-2xl border border-base-300">
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
                                <td>{selectedMentorUser.teamNames || ""}</td>
                                <td>
                                    <span class="flex items-center gap-2">
                                        <Icon
                                            url={avatarURL(selectedMentorUser)}
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
            <div class="form-control w-full">
                <div class="flex flex-1 w-full justify-between items-end">
                    <span class="label-text font-medium">メンティ</span>
                    <button
                        type="button"
                        class="btn btn-outline btn-xs mt-1 w-fit"
                        onclick={() => (menteeModalOpen = true)}
                        >メンティを選択</button
                    >
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
                                        {u.departmentNames || "所属無し"}
                                    </td>
                                    <td>
                                        {u.teamNames || ""}
                                    </td>
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

        <div class="flex flex-wrap gap-2 pt-2 justify-end">
            <button type="button" class="btn btn-primary" onclick={requestSave}
                >作成</button
            >
            <a href="/planning/pairs" class="btn btn-ghost">戻る</a>
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
    pickRootId="pairing-form-mentor-name"
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
    pickRootId="pairing-form-mentee-name"
    onApply={(uuids) => {
        menteeUuids = uuids;
        menteeModalOpen = false;
    }}
    onClose={() => {
        menteeModalOpen = false;
    }}
/>
