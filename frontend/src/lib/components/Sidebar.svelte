<script lang="ts">
	import type { AvatarMe, NavItem } from '$lib/types';
	import { buildNavigation, settingsNav } from '$lib/navigation';
	import { ui } from '$lib/stores/ui.svelte';
	import Logo from './Logo.svelte';
	import Home from 'lucide-svelte/icons/home';
	import CalendarDays from 'lucide-svelte/icons/calendar-days';
	import Shield from 'lucide-svelte/icons/shield';
	import BarChart2 from 'lucide-svelte/icons/bar-chart-2';
	import Settings from 'lucide-svelte/icons/settings';

	let { me, pathname }: { me: AvatarMe; pathname: string } = $props();

	const icons: Record<string, typeof Home> = {
		home: Home,
		'calendar-days': CalendarDays,
		shield: Shield,
		'bar-chart-2': BarChart2,
		settings: Settings
	};

	let navigation = $derived(buildNavigation(pathname, me));
	let navSettings = $derived(settingsNav(pathname));

	function navLinkClass(active: boolean) {
		return active
			? 'bg-primary/10 text-primary'
			: 'text-base-content/60 hover:bg-base-200 hover:text-base-content';
	}
</script>

<aside
	class="hidden shrink-0 border-r border-base-300 bg-base-100 transition-[width] duration-200 md:flex md:flex-col md:min-h-0 {ui.sidebarOpen
		? 'w-64'
		: 'w-20'}"
>
	<div class="flex h-16 shrink-0 items-center gap-3 border-b border-base-300 px-5">
		<div class="grid h-8 w-8 place-items-center rounded-xl bg-primary text-primary-content">
			<Logo />
		</div>
		{#if ui.sidebarOpen}
			<div>
				<div class="font-display text-sm font-semibold">Growlink</div>
				<div class="mt-0.5 text-[10px] text-base-content/50">1on1 ミーティング</div>
			</div>
		{/if}
	</div>
	<nav class="flex-1 min-h-0 space-y-1 overflow-y-auto p-3">
		{#each navigation as item}
			{@const Icon = icons[item.icon] ?? Home}
			<a
				href={item.href}
				class="flex h-11 items-center gap-3 rounded-xl px-3 text-sm font-medium transition {navLinkClass(item.active)}"
			>
				<Icon class="h-5 w-5 shrink-0" />
				{#if ui.sidebarOpen}
					<span>{item.label}</span>
				{/if}
			</a>
		{/each}
	</nav>
	<div class="mt-auto shrink-0 border-t border-base-300 p-3">
		<a
			href={navSettings.href}
			class="flex h-11 items-center gap-3 rounded-xl px-3 text-sm font-medium transition {navLinkClass(navSettings.active)}"
		>
			<Settings class="h-5 w-5 shrink-0" />
			{#if ui.sidebarOpen}
				<span>{navSettings.label}</span>
			{/if}
		</a>
	</div>
</aside>
