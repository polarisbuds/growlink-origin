<script lang="ts">
    import { onMount } from "svelte";
    import { api } from "$lib/api";
    import ProgressChart from "$lib/components/ProgressChart.svelte";
    import type {
        ProgressDepartmentRow,
        ProgressDoneMonthPoint,
    } from "$lib/progress-types";

    type ProgressViewTab = "departments" | "trend";

    let viewTab = $state<ProgressViewTab>("departments");
    let departmentRows = $state<ProgressDepartmentRow[]>([]);
    let doneTrendPoints = $state<ProgressDoneMonthPoint[]>([]);
    let departmentLoading = $state(false);
    let doneTrendLoading = $state(false);
    let dateFrom = $state("");
    let dateTo = $state("");
    let periodLabel = $state("全期間");

    let trendRangeReady = $derived(Boolean(dateFrom && dateTo));

    let chartPoints = $derived(
        doneTrendPoints.map((p) => ({ label: p.label, done: p.done })),
    );

    async function setTomonthRange() {
        const now = new Date();

        const start = new Date(now.getFullYear(), now.getMonth(), 1);
        const end = new Date(now.getFullYear(), now.getMonth() + 1, 0);

        const formatDate = (dateObj: Date) => {
            const y = dateObj.getFullYear();
            const m = String(dateObj.getMonth() + 1).padStart(2, "0");
            const d = String(dateObj.getDate()).padStart(2, "0");
            return `${y}-${m}-${d}`;
        };
        dateFrom = formatDate(start);
        dateTo = formatDate(end);
    }

    function buildQuery() {
        const p = new URLSearchParams();
        if (dateFrom) p.set("from", dateFrom);
        if (dateTo) p.set("to", dateTo);
        const s = p.toString();
        return s ? `?${s}` : "";
    }

    function updatePeriodLabel() {
        if (dateFrom && dateTo) periodLabel = `${dateFrom} 〜 ${dateTo}`;
        else if (dateFrom) periodLabel = `${dateFrom} 以降`;
        else if (dateTo) periodLabel = `${dateTo} まで`;
        else periodLabel = "全期間";
    }

    async function loadDepartmentTable() {
        departmentLoading = true;
        try {
            const res = await api<{ rows: ProgressDepartmentRow[] }>(
                "GET",
                `/api/v1/progress/departments${buildQuery()}`,
            );
            departmentRows = res.rows ?? [];
        } catch {
            departmentRows = [];
        } finally {
            departmentLoading = false;
        }
    }

    async function loadDoneTrend() {
        if (!dateFrom || !dateTo) {
            doneTrendPoints = [];
            return;
        }
        doneTrendLoading = true;
        try {
            doneTrendPoints = await api<ProgressDoneMonthPoint[]>(
                "GET",
                `/api/v1/progress/done-trend${buildQuery()}`,
            );
        } catch {
            doneTrendPoints = [];
        } finally {
            doneTrendLoading = false;
        }
    }

    async function reloadAll() {
        updatePeriodLabel();
        await Promise.all([loadDepartmentTable(), loadDoneTrend()]);
    }

    onMount(async () => {
        await setTomonthRange();
        await reloadAll();
    });

    function switchViewTab(tab: ProgressViewTab) {
        viewTab = tab;
    }

    function deptIndent(depth: number) {
        return depth > 0 ? `${"　".repeat(depth)}` : "";
    }
</script>

<svelte:head><title>進捗 - Growlink</title></svelte:head>

<h1 class="font-display text-2xl font-semibold">進捗</h1>

<div role="tablist" class="tabs tabs-bordered mt-6 max-w-3xl">
    <button
        type="button"
        role="tab"
        class="tab"
        class:tab-active={viewTab === "departments"}
        onclick={() => switchViewTab("departments")}>
        部署ごとの実施状況
    </button>
    <button
        type="button"
        role="tab"
        class="tab"
        class:tab-active={viewTab === "trend"}
        onclick={() => switchViewTab("trend")}>
        実施数の遷移
    </button>
</div>

<div class="flex flex-1 justify-between mt-5 items-center">
    <p class="text-sm text-base-content/60">
        表示期間: {periodLabel}
    </p>

    <div class="flex flex-wrap items-end gap-3 justify-end">
        <label class="form-control">
            <span class="label-text text-xs">開始日</span>
            <input
                type="date"
                class="input input-bordered input-sm"
                bind:value={dateFrom} />
        </label>
        <label class="form-control">
            <span class="label-text text-xs">終了日</span>
            <input
                type="date"
                class="input input-bordered input-sm"
                bind:value={dateTo} />
        </label>
        <button type="button" class="btn btn-primary btn-sm" onclick={reloadAll}
            >期間を適用</button>
    </div>
</div>

{#if viewTab === "departments"}
    <section
        class="mt-3 overflow-hidden rounded-3xl border border-base-300 bg-base-100 shadow-sm">
        {#if departmentLoading}
            <p class="p-12 text-center text-sm text-base-content/50">
                読み込み中…
            </p>
        {:else if departmentRows.length === 0}
            <p class="p-12 text-center text-sm text-base-content/50">
                表示するデータがありません
            </p>
        {:else}
            <div class="overflow-x-auto">
                <table class="table table-sm">
                    <thead>
                        <tr class="border-b border-base-300 bg-base-200/40">
                            <th>部署</th>
                            <th class="text-right">ユーザー数</th>
                            <th class="text-right">ペア数</th>
                            <th class="text-right">未実施[件]</th>
                            <th class="text-right">完了[件]</th>
                            <th class="text-right">実施率</th>
                        </tr>
                    </thead>
                    <tbody>
                        {#each departmentRows as row (row.isTotal ? "total" : row.uuid + row.name)}
                            <tr
                                class="border-b {row.isTotal
                                    ? 'font-semibold bg-base-200/20'
                                    : 'hover:bg-base-300/10'}">
                                <td>{deptIndent(row.depth)}{row.name}</td>
                                <td class="text-right">{row.userCount}</td>
                                <td class="text-right">{row.pairCount}</td>
                                <td class="text-right">{row.notDoneCount}</td>
                                <td class="text-right">{row.doneCount}</td>
                                <td class="text-right"
                                    >{row.donePercent.toFixed(1)}%</td>
                            </tr>
                            {#each row.teams ?? [] as team (team.uuid + team.name)}
                                <tr
                                    class="border-b border-base-200/80 hover:bg-base-300/10">
                                    <td
                                        >{deptIndent(
                                            row.depth + 1,
                                        )}{team.name}</td>
                                    <td class="text-right">{team.userCount}</td>
                                    <td class="text-right">{team.pairCount}</td>
                                    <td class="text-right"
                                        >{team.notDoneCount}</td>
                                    <td class="text-right">{team.doneCount}</td>
                                    <td class="text-right"
                                        >{team.donePercent.toFixed(1)}%</td>
                                </tr>
                            {/each}
                        {/each}
                    </tbody>
                </table>
            </div>
        {/if}
    </section>
{:else}
    <section
        class="mt-10 rounded-3xl border border-base-300 bg-base-100 p-6 shadow-sm">
        <div>
            <h2 class="text-lg font-semibold">完了件数</h2>
            <p class="mt-1 text-sm text-base-content/60">
                表示期間の全部署合計（月別）
            </p>
        </div>
        <div class="mt-6">
            {#if !trendRangeReady}
                <p class="py-12 text-center text-sm text-base-content/50">
                    開始日と終了日を指定してください
                </p>
            {:else if doneTrendLoading}
                <p class="py-12 text-center text-sm text-base-content/50">
                    読み込み中…
                </p>
            {:else}
                <ProgressChart points={chartPoints} />
            {/if}
        </div>
    </section>
{/if}
