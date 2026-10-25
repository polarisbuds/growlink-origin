<script lang="ts">
	import type { AvatarMe } from '$lib/types';
	import { ui } from '$lib/stores/ui.svelte';
	import ThemeSelect from './ThemeSelect.svelte';
	import PanelLeft from 'lucide-svelte/icons/panel-left';
	import ChevronRight from 'lucide-svelte/icons/chevron-right';
	import ChevronDown from 'lucide-svelte/icons/chevron-down';
	import User from 'lucide-svelte/icons/user';
	import LogOut from 'lucide-svelte/icons/log-out';

	let { me, breadcrumb }: { me: AvatarMe; breadcrumb?: string } = $props();

	let avatarUrl = $derived(me.imageUri ? `/uploads/${me.imageUri}` : null);
	let initial = $derived((me.name || '?').slice(0, 1));
</script>

<header class="sticky top-0 z-30 flex h-16 items-center gap-3 border-b border-base-300 bg-base-100 px-4 md:px-6">
	<button
		type="button"
		class="btn btn-ghost btn-sm btn-square border border-base-300"
		onclick={() => ui.toggleSidebar()}
		aria-label="サイドバーを開閉"
	>
		<PanelLeft class="h-5 w-5" />
	</button>
	<nav class="flex items-center gap-2 text-sm text-base-content/60">
		<a href="/" class="hover:text-base-content">ホーム</a>
		{#if breadcrumb}
			<ChevronRight class="h-4 w-4" />
			<span class="font-medium text-base-content">{breadcrumb}</span>
		{/if}
	</nav>
	<div class="ml-auto flex items-center gap-2">
		<ThemeSelect />
		<div class="dropdown dropdown-end">
			<button type="button" tabindex="0" class="btn btn-ghost btn-sm gap-2">
				{#if avatarUrl}
					<img src={avatarUrl} alt="" class="h-7 w-7 rounded-full object-cover" />
				{:else}
					<span
						class="grid h-7 w-7 place-items-center rounded-full bg-primary text-xs font-semibold text-primary-content"
					>
						{initial}
					</span>
				{/if}
				<span class="hidden sm:inline">{me.name}</span>
				<ChevronDown class="h-4 w-4 opacity-60" />
			</button>
			<ul
				role="menu"
				tabindex="0"
				class="menu dropdown-content z-50 mt-2 w-52 rounded-box border border-base-300 bg-base-100 p-2 shadow-lg"
			>
				<li>
					<a href="/setting"><User class="h-4 w-4 inline" /> 設定</a>
				</li>
				<li>
					<a href="/logout" data-sveltekit-reload><LogOut class="h-4 w-4 inline" /> ログアウト</a>
				</li>
			</ul>
		</div>
	</div>
</header>
