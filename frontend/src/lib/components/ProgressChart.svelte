<script lang="ts">
	import * as d3 from 'd3';
	import { onMount } from 'svelte';
	import type { ProgressBreakdownRow } from '$lib/progress-types';

	let { rows = [], groupLabel = '部署' }: { rows?: ProgressBreakdownRow[]; groupLabel?: string } =
		$props();

	let container: HTMLDivElement | undefined = $state();
	let width = $state(640);

	const margin = { top: 16, right: 16, bottom: 72, left: 48 };
	const series = [
		{ key: 'pending' as const, label: '未着手', color: '#570df8' },
		{ key: 'draft' as const, label: '下書き', color: '#fbbd23' },
		{ key: 'done' as const, label: '完了', color: '#36d399' }
	];

	function draw() {
		if (!container) return;
		const data = rows;
		container.innerHTML = '';
		if (!data.length) {
			const empty = document.createElement('p');
			empty.className = 'py-12 text-center text-sm text-base-content/50';
			empty.textContent = '表示するデータがありません';
			container.appendChild(empty);
			return;
		}

		const innerW = Math.max(280, width - margin.left - margin.right);
		const innerH = 280;
		const height = innerH + margin.top + margin.bottom;

		const svg = d3
			.select(container)
			.append('svg')
			.attr('width', width)
			.attr('height', height)
			.attr('role', 'img')
			.attr('aria-label', `${groupLabel}別の面談ステータス`);

		const g = svg.append('g').attr('transform', `translate(${margin.left},${margin.top})`);

		const x = d3
			.scaleBand()
			.domain(data.map((d) => d.name))
			.range([0, innerW])
			.padding(0.25);

		const maxTotal = d3.max(data, (d) => d.pending + d.draft + d.done) ?? 0;
		const y = d3
			.scaleLinear()
			.domain([0, Math.max(1, maxTotal)])
			.nice()
			.range([innerH, 0]);

		const stack = d3.stack<ProgressBreakdownRow>().keys(series.map((s) => s.key));
		const stacked = stack(data);

		g.append('g')
			.selectAll('g')
			.data(stacked)
			.join('g')
			.attr('fill', (_, i) => series[i]?.color ?? '#888')
			.selectAll('rect')
			.data((layer) => layer)
			.join('rect')
			.attr('x', (d) => x(d.data.name) ?? 0)
			.attr('y', (d) => y(d[1]))
			.attr('height', (d) => Math.max(0, y(d[0]) - y(d[1])))
			.attr('width', x.bandwidth())
			.attr('rx', 4);

		g.append('g')
			.attr('transform', `translate(0,${innerH})`)
			.call(d3.axisBottom(x))
			.selectAll('text')
			.attr('transform', 'rotate(-32)')
			.style('text-anchor', 'end')
			.attr('dx', '-0.4em')
			.attr('dy', '0.2em');

		g.append('g').call(d3.axisLeft(y).ticks(5).tickFormat(d3.format('d')));

		const legend = svg
			.append('g')
			.attr('transform', `translate(${margin.left}, ${height - 20})`);
		let lx = 0;
		for (const s of series) {
			const item = legend.append('g').attr('transform', `translate(${lx}, 0)`);
			item.append('rect').attr('width', 12).attr('height', 12).attr('fill', s.color).attr('rx', 2);
			item.append('text').attr('x', 16).attr('y', 10).attr('fill', 'currentColor').text(s.label);
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

<div bind:this={container} class="w-full min-h-[320px] text-base-content/80"></div>
