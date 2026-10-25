<script lang="ts">
	import { themeOptions, applyTheme, THEME_KEY, DEFAULT_THEME } from '$lib/themes';
	import Palette from 'lucide-svelte/icons/palette';
	import ChevronDown from 'lucide-svelte/icons/chevron-down';

	let open = $state(false);
	let theme = $state(
		typeof localStorage !== 'undefined'
			? localStorage.getItem(THEME_KEY) || DEFAULT_THEME
			: DEFAULT_THEME
	);

	function pick(value: string) {
		theme = value;
		applyTheme(value);
		open = false;
	}
</script>

<div class="dropdown dropdown-end">
	<button
		type="button"
		tabindex="0"
		class="btn btn-ghost btn-sm gap-1"
		onclick={() => (open = !open)}
		aria-label="テーマ"
	>
		<Palette class="h-4 w-4" />
		<ChevronDown class="h-3 w-3 opacity-60" />
	</button>
	{#if open}
		<ul
			class="menu dropdown-content z-50 mt-2 max-h-72 w-52 overflow-y-auto rounded-box border border-base-300 bg-base-100 p-2 shadow-lg"
		>
			{#each themeOptions as opt}
				<li>
					<button type="button" class={theme === opt.value ? 'active' : ''} onclick={() => pick(opt.value)}>
						{opt.label}
					</button>
				</li>
			{/each}
		</ul>
	{/if}
</div>
