<script lang="ts">
    import { onMount } from "svelte";
    import { api } from "$lib/api";
    import { statusBadge, statusLabel } from "$lib/status";
    import { formatDate, type MeetingListItem } from "$lib/meeting-helpers";

    let meetings = $state<MeetingListItem[]>([]);
    let search = $state("");

    let filtered = $derived.by(() => {
        const q = search.trim().toLowerCase();
        if (!q) return meetings;
        return meetings.filter((m) => {
            const hay = [
                m.template?.title,
                m.pairing?.mentor?.name,
                m.pairing?.mentee?.name,
                statusLabel(m.status ?? 0),
            ]
                .filter(Boolean)
                .join(" ")
                .toLowerCase();
            return hay.includes(q);
        });
    });

    onMount(async () => {
        meetings = await api<MeetingListItem[]>(
            "GET",
            "/api/v1/meetings/involved",
        );
    });
</script>

<svelte:head><title>面談一覧 - Growlink</title></svelte:head>

<h1 class="font-display text-2xl font-semibold">面談一覧</h1>
<input
    type="search"
    class="input input-bordered mt-4 w-full max-w-md placeholder-current"
    placeholder="検索"
    bind:value={search}
/>
<div
    class="mt-6 overflow-x-auto rounded-2xl border border-base-300 bg-base-100"
>
    <table class="table">
        <thead>
            <tr>
                <th>テンプレート</th>
                <th>ペア</th>
                <th>実施日</th>
                <th>状態</th>
                <th></th>
            </tr>
        </thead>
        <tbody>
            {#each filtered as m (m.uuid)}
                <tr>
                    <td>{m.template?.title || "—"}</td>
                    <td>
                        {m.pairing?.mentor?.name || "—"} /
                        {m.pairing?.mentee?.name || "—"}
                    </td>
                    <td>{formatDate(m.implementedAt)}</td>
                    <td>
                        <span class="badge {statusBadge(m.status ?? 0)}"
                            >{statusLabel(m.status ?? 0)}</span
                        >
                    </td>
                    <td class="text-right">
                        <a href="/meeting/{m.uuid}" class="btn btn-ghost btn-xs"
                            >開く</a
                        >
                    </td>
                </tr>
            {/each}
        </tbody>
    </table>
</div>
