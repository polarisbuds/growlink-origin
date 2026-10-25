<script lang="ts">
	type Chip = { id: number; name: string };
	type User = { id: number; name: string; email: string };

	let {
		chips,
		placeholder = '名前で検索して追加',
		query = $bindable(''),
		open = $bindable(false),
		filteredUsers,
		onPick,
		onRemove
	}: {
		chips: Chip[];
		placeholder?: string;
		query?: string;
		open?: boolean;
		filteredUsers: User[];
		onPick: (user: User) => void;
		onRemove: (id: number) => void;
	} = $props();
</script>

<div class="relative mt-1">
	<div
		class="flex min-h-12 w-full flex-wrap items-center gap-1.5 rounded-lg border border-base-300 bg-base-100 px-3 py-2 focus-within:border-primary"
	>
		{#each chips as chip (chip.id)}
			<span
				class="inline-flex max-w-full shrink-0 items-center gap-0.5 rounded-full bg-base-200 pl-2.5 pr-0.5 py-0.5 text-sm text-base-content"
			>
				<span class="max-w-[11rem] truncate" title={chip.name}>{chip.name}</span>
				<button
					type="button"
					class="btn btn-ghost btn-xs btn-circle h-5 min-h-5 w-5 shrink-0 opacity-70 hover:opacity-100"
					aria-label="{chip.name}を削除"
					onclick={() => onRemove(chip.id)}>×</button
				>
			</span>
		{/each}
		<input
			type="text"
			class="min-w-[7rem] flex-1 border-0 bg-transparent py-1 text-sm outline-none focus:outline-none"
			{placeholder}
			autocomplete="off"
			bind:value={query}
			onfocus={() => (open = true)}
			oninput={() => (open = true)}
		/>
	</div>
	{#if open}
		<ul
			class="absolute left-0 top-full z-20 mt-1 max-h-48 w-full overflow-y-auto rounded-2xl border border-base-300 bg-base-100 py-1 shadow-lg"
		>
			{#each filteredUsers as u (u.id)}
				<li>
					<button
						type="button"
						class="btn btn-ghost btn-sm w-full justify-start font-normal"
						onclick={() => onPick(u)}>{u.name} ({u.email})</button
					>
				</li>
			{/each}
		</ul>
	{/if}
</div>
