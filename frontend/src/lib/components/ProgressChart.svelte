<script lang="ts">
    import { max } from "d3-array";
    import { axisBottom, axisLeft } from "d3-axis";
    import { format } from "d3-format";
    import { scaleBand, scaleLinear } from "d3-scale";
    import { select } from "d3-selection";
    import "d3-transition";
    import { onMount } from "svelte";

    export type ProgressDoneChartPoint = {
        label: string;
        done: number;
    };

    let {
        points = [],
    }: {
        points?: ProgressDoneChartPoint[];
    } = $props();

    let container: HTMLDivElement | undefined = $state();
    let width = $state(640);
    let tooltip = $state<{ x: number; y: number; text: string } | null>(null);

    const margin = { top: 16, right: 16, bottom: 72, left: 48 };
    const doneColor = "oklch(var(--su))";

    function draw() {
        if (!container) return;
        const data = points;
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
            .attr("aria-label", "月別の完了件数");

        const g = svg
            .append("g")
            .attr("transform", `translate(${margin.left},${margin.top})`);

        const x = scaleBand()
            .domain(data.map((d) => d.label))
            .range([0, innerW])
            .padding(0.25);

        const maxDone = max(data, (d) => d.done) ?? 0;
        const y = scaleLinear()
            .domain([0, Math.max(1, maxDone)])
            .nice()
            .range([innerH, 0]);

        g.selectAll("rect.bar")
            .data(data)
            .join("rect")
            .attr("class", "bar")
            .attr("x", (d) => x(d.label) ?? 0)
            .attr("width", x.bandwidth())
            .attr("rx", 4)
            .attr("fill", doneColor)
            .attr("y", innerH)
            .attr("height", 0)
            .transition()
            .duration(600)
            .attr("y", (d) => y(d.done))
            .attr("height", (d) => Math.max(0, innerH - y(d.done)))
            .selection()
            .on("mouseenter", (ev, d) => {
                tooltip = {
                    x: ev.clientX,
                    y: ev.clientY,
                    text: `${d.label} / 完了: ${d.done}件`,
                };
            })
            .on("mouseleave", () => {
                tooltip = null;
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
