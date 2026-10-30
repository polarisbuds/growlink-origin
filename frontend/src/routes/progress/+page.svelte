<script lang="ts">
    import { onMount } from "svelte";
    import { api } from "$lib/api";
    import ProgressChart from "$lib/components/ProgressChart.svelte";
    import type { ProgressBreakdownRow } from "$lib/progress-types";

    type ProgressSummary = {
        pending: number;
        draft: number;
        done: number;
        total: number;
        donePercent?: number;
    };

    type BreakdownGroup = "department" | "team";

    let summary = $state<ProgressSummary>({
        pending: 0,
        draft: 0,
        done: 0,
        total: 0,
        donePercent: 0,
    });
    let breakdownGroup = $state<BreakdownGroup>("department");
    let breakdownRows = $state<ProgressBreakdownRow[]>([]);
    let breakdownLoading = $state(false);
    let dateFrom = $state("");
    let dateTo = $state("");
    let periodLabel = $state("全期間");
    let drilldownOpen = $state(false);
    let drilldownTitle = $state("");
    let drilldownUsers = $state<string[]>([]);

    let groupLabel = $derived(
        breakdownGroup === "department" ? "部署" : "チーム",
    );

    function buildQuery(extra: Record<string, string> = {}) {
        const p = new URLSearchParams();
        if (dateFrom) p.set("from", dateFrom);
        if (dateTo) p.set("to", dateTo);
        for (const [k, v] of Object.entries(extra)) p.set(k, v);
        const s = p.toString();
        return s ? `?${s}` : "";
    }

    function updatePeriodLabel() {
        if (dateFrom && dateTo) periodLabel = `${dateFrom} 〜 ${dateTo}`;
        else if (dateFrom) periodLabel = `${dateFrom} 以降`;
        else if (dateTo) periodLabel = `${dateTo} まで`;
        else periodLabel = "全期間";
    }

    async function loadSummary() {
        try {
            summary = await api<ProgressSummary>(
                "GET",
                `/api/v1/progress/summary${buildQuery()}`,
            );
        } catch {
            summary = {
                pending: 0,
                draft: 0,
                done: 0,
                total: 0,
                donePercent: 0,
            };
        }
    }

    async function loadBreakdown() {
        breakdownLoading = true;
        try {
            breakdownRows = await api<ProgressBreakdownRow[]>(
                "GET",
                `/api/v1/progress/breakdown${buildQuery({ group: breakdownGroup })}`,
            );
        } catch {
            breakdownRows = [];
        } finally {
            breakdownLoading = false;
        }
    }

    async function reloadAll() {
        updatePeriodLabel();
        await loadSummary();
        await loadBreakdown();
    }

    onMount(reloadAll);

    async function setBreakdownGroup(group: BreakdownGroup) {
        if (breakdownGroup === group) return;
        breakdownGroup = group;
        await loadBreakdown();
    }

    function onChartClick(payload: {
        row: ProgressBreakdownRow;
        status: "pending" | "draft" | "done";
        users: string[];
    }) {
        const statusLabel =
            payload.status === "pending"
                ? "未着手"
                : payload.status === "draft"
                  ? "下書き"
                  : "完了";
        drilldownTitle = `${payload.row.name}（${statusLabel}）`;
        drilldownUsers = payload.users.length
            ? payload.users
            : ["（該当ユーザなし）"];
        drilldownOpen = true;
    }
</script>

<svelte:head><title>進捗 - Growlink</title></svelte:head>

<h1 class="font-display text-2xl font-semibold">進捗</h1>
<p class="mt-1 text-sm text-base-content/60">
    面談ステータス別の件数（表示期間: {periodLabel}）
</p>

<div class="mt-4 flex flex-wrap items-end gap-3">
    <label class="form-control">
        <span class="label-text text-xs">開始日</span>
        <input
            type="date"
            class="input input-bordered input-sm"
            bind:value={dateFrom}
        />
    </label>
    <label class="form-control">
        <span class="label-text text-xs">終了日</span>
        <input
            type="date"
            class="input input-bordered input-sm"
            bind:value={dateTo}
        />
    </label>
    <button type="button" class="btn btn-primary btn-sm" onclick={reloadAll}
        >期間を適用</button
    >
</div>

<div class="mt-8 grid gap-4 md:grid-cols-2 xl:grid-cols-4">
    <div class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm">
        <p class="text-sm text-base-content/60">未着手</p>
        <p class="mt-2 text-3xl font-semibold text-primary">
            {summary.pending}
        </p>
    </div>
    <div class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm">
        <p class="text-sm text-base-content/60">下書き</p>
        <p class="mt-2 text-3xl font-semibold text-warning">{summary.draft}</p>
    </div>
    <div class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm">
        <p class="text-sm text-base-content/60">完了</p>
        <p class="mt-2 text-3xl font-semibold text-success">{summary.done}</p>
    </div>
    <div class="rounded-3xl border border-base-300 bg-base-100 p-5 shadow-sm">
        <p class="text-sm text-base-content/60">合計 / 完了率</p>
        <p class="mt-2 text-3xl font-semibold">{summary.total}</p>
        <p class="text-sm text-base-content/60">
            完了率 {(summary.donePercent ?? 0).toFixed(1)}%
        </p>
    </div>
</div>

<section
    class="mt-10 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm"
>
    <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
            <h2 class="text-lg font-semibold">ステータス内訳</h2>
            <p class="mt-1 text-sm text-base-content/60">
                {groupLabel}ごとの面談件数
            </p>
        </div>
        <div role="tablist" class="tabs tabs-bordered">
            <button
                type="button"
                role="tab"
                class="tab"
                class:tab-active={breakdownGroup === "department"}
                onclick={() => setBreakdownGroup("department")}>部署</button
            >
            <button
                type="button"
                role="tab"
                class="tab"
                class:tab-active={breakdownGroup === "team"}
                onclick={() => setBreakdownGroup("team")}>チーム</button
            >
        </div>
    </div>
    <div class="mt-6">
        {#if breakdownLoading}
            <p class="py-12 text-center text-sm text-base-content/50">
                読み込み中…
            </p>
        {:else}
            <ProgressChart
                rows={breakdownRows}
                {groupLabel}
                onSegmentClick={onChartClick}
            />
        {/if}
    </div>
    {#if breakdownRows.length > 0}
        <div class="mt-8 overflow-x-auto">
            <table class="table table-sm">
                <thead>
                    <tr>
                        <th>{groupLabel}</th>
                        <th>未着手</th>
                        <th>下書き</th>
                        <th>完了</th>
                        <th>合計</th>
                        <th>完了率</th>
                    </tr>
                </thead>
                <tbody>
                    {#each breakdownRows as row (row.id + row.name)}
                        <tr>
                            <td>{row.name}</td>
                            <td>{row.pending}</td>
                            <td>{row.draft}</td>
                            <td>{row.done}</td>
                            <td>{row.pending + row.draft + row.done}</td>
                            <td>{(row.donePercent ?? 0).toFixed(1)}%</td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    {/if}
</section>

<dialog class="modal" class:modal-open={drilldownOpen}>
    <div class="modal-box">
        <h3 class="font-bold text-lg">{drilldownTitle}</h3>
        <ul class="mt-4 list-disc space-y-1 pl-5 text-sm">
            {#each drilldownUsers as name, i (i)}
                <li>{name}</li>
            {/each}
        </ul>
        <div class="modal-action">
            <button
                type="button"
                class="btn"
                onclick={() => (drilldownOpen = false)}>閉じる</button
            >
        </div>
    </div>
    <form method="dialog" class="modal-backdrop">
        <button type="button" onclick={() => (drilldownOpen = false)}
            >close</button
        >
    </form>
</dialog>
