<script lang="ts">
    import { max } from "d3-array";
    import { axisBottom, axisLeft } from "d3-axis";
    import { format } from "d3-format";
    import { scaleBand, scaleLinear } from "d3-scale";
    import { select } from "d3-selection";
    import { stack, type Series } from "d3-shape";
    import "d3-transition";
    import { onMount } from "svelte";
    import type { ProgressBreakdownRow } from "$lib/progress-types";

    let {
        rows = [],
        groupLabel = "部署",
        onSegmentClick,
    }: {
        rows?: ProgressBreakdownRow[];
        groupLabel?: string;
        onSegmentClick?: (payload: {
            row: ProgressBreakdownRow;
            status: "pending" | "draft" | "done";
            users: string[];
        }) => void;
    } = $props();

    let container: HTMLDivElement | undefined = $state();
    let width = $state(640);
    let tooltip = $state<{ x: number; y: number; text: string } | null>(null);

    const margin = { top: 16, right: 16, bottom: 72, left: 48 };
    const series = [
        { key: "pending" as const, label: "未着手", color: "oklch(var(--p))" },
        { key: "draft" as const, label: "下書き", color: "oklch(var(--wa))" },
        { key: "done" as const, label: "完了", color: "oklch(var(--su))" },
    ];

    function usersFor(
        row: ProgressBreakdownRow,
        key: "pending" | "draft" | "done",
    ) {
        if (key === "pending") return row.usersPending || [];
        if (key === "draft") return row.usersDraft || [];
        return row.usersDone || [];
    }

    function draw() {
        if (!container) return;
        const data = rows;
        container.innerHTML = "";
        tooltip = null;
        if (!data.length) {
            const empty = document.createElement("p");
            empty.className = "py-12 text-center text-sm text-base-content/50";
            empty.textContent = "表示するデータがありません";
            container.appendChild(empty);
            return;
        }

        const innerW = Math.max(280, width - margin.left - margin.right);
        const innerH = 280;
        const height = innerH + margin.top + margin.bottom;

        const svg = select(container)
            .append("svg")
            .attr("width", width)
            .attr("height", height)
            .attr("role", "img")
            .attr("aria-label", `${groupLabel}別の面談ステータス`);

        const g = svg
            .append("g")
            .attr("transform", `translate(${margin.left},${margin.top})`);

        const x = scaleBand()
            .domain(data.map((d) => d.name))
            .range([0, innerW])
            .padding(0.25);

        const maxTotal = max(data, (d) => d.pending + d.draft + d.done) ?? 0;
        const y = scaleLinear()
            .domain([0, Math.max(1, maxTotal)])
            .nice()
            .range([innerH, 0]);

        const stackGen = stack<ProgressBreakdownRow>().keys(
            series.map((s) => s.key),
        );
        const stacked = stackGen(data);

        g.selectAll("g.layer")
            .data(stacked)
            .join("g")
            .attr("class", "layer")
            .attr(
                "fill",
                (layer) =>
                    series.find((s) => s.key === layer.key)?.color ??
                    "hsl(var(--bc))",
            )
            .selectAll("rect")
            .data((layer) => layer)
            .join("rect")
            .attr("x", (d) => x(d.data.name) ?? 0)
            .attr("width", x.bandwidth())
            .attr("rx", 4)
            .attr("y", innerH)
            .attr("height", 0)
            .transition()
            .duration(600)
            .attr("y", (d) => y(d[1]))
            .attr("height", (d) => Math.max(0, y(d[0]) - y(d[1])))
            .selection()
            .on("mouseenter", (ev, d) => {
                const parent = (ev.currentTarget as SVGRectElement)
                    .parentElement;
                const layerKey =
                    parent instanceof SVGGElement
                        ? (parent as SVGGElement & { __layerKey?: string })
                              .__layerKey
                        : undefined;
                const key = (layerKey || "pending") as
                    | "pending"
                    | "draft"
                    | "done";
                const s = series.find((x) => x.key === key);
                const count = d.data[key];
                const total = d.data.pending + d.data.draft + d.data.done;
                const pct = total > 0 ? Math.round((count / total) * 100) : 0;
                tooltip = {
                    x: ev.clientX,
                    y: ev.clientY,
                    text: `${d.data.name} / ${s?.label ?? key}: ${count}件 (${pct}%)`,
                };
            })
            .on("mouseleave", () => {
                tooltip = null;
            })
            .on("click", (ev, d) => {
                const parent = (ev.currentTarget as SVGRectElement)
                    .parentElement;
                const layerKey =
                    parent instanceof SVGGElement
                        ? (parent as SVGGElement & { __layerKey?: string })
                              .__layerKey
                        : undefined;
                const key = (layerKey || "pending") as
                    | "pending"
                    | "draft"
                    | "done";
                onSegmentClick?.({
                    row: d.data,
                    status: key,
                    users: usersFor(d.data, key),
                });
            });

        g.selectAll<SVGGElement, Series<ProgressBreakdownRow, string>>(
            "g.layer",
        ).each(function (layer) {
            (this as SVGGElement & { __layerKey?: string }).__layerKey =
                layer.key;
        });

        g.append("g")
            .attr("transform", `translate(0,${innerH})`)
            .call(axisBottom(x))
            .selectAll("text")
            .attr("transform", "rotate(-32)")
            .style("text-anchor", "end")
            .attr("dx", "-0.4em")
            .attr("dy", "0.2em");

        g.append("g").call(axisLeft(y).ticks(5).tickFormat(format("d")));

        const legend = svg
            .append("g")
            .attr("transform", `translate(${margin.left}, ${height - 20})`);
        let lx = 0;
        for (const s of series) {
            const item = legend
                .append("g")
                .attr("transform", `translate(${lx}, 0)`);
            item.append("rect")
                .attr("width", 12)
                .attr("height", 12)
                .attr("fill", s.color)
                .attr("rx", 2);
            item.append("text")
                .attr("x", 16)
                .attr("y", 10)
                .attr("fill", "currentColor")
                .text(s.label);
            lx += s.label.length * 10 + 36;
        }
    }

    $effect(() => {
        draw();
    });

    onMount(() => {
        if (!container) return;
        const ro = new ResizeObserver((entries) => {
            const entry = entries[0];
            if (entry) width = Math.floor(entry.contentRect.width);
        });
        ro.observe(container);
        return () => ro.disconnect();
    });
</script>

<div class="relative w-full min-h-[320px] text-base-content/80">
    <div bind:this={container} class="w-full"></div>
    {#if tooltip}
        <div
            class="pointer-events-none fixed z-50 rounded-lg border border-base-300 bg-base-100 px-3 py-2 text-xs shadow-lg"
            style="left: {tooltip.x + 12}px; top: {tooltip.y + 12}px;">
            {tooltip.text}
        </div>
    {/if}
</div>
