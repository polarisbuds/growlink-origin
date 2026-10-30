<script lang="ts">
	import type { AvatarMe } from '$lib/types';
	import Sidebar from './Sidebar.svelte';
	import Header from './Header.svelte';

	let {
		me,
		pathname,
		breadcrumb,
		children
	}: {
		me: AvatarMe;
		pathname: string;
		breadcrumb?: string;
		children: import('svelte').Snippet;
	} = $props();
</script>

<div class="flex h-screen overflow-hidden">
	<Sidebar {me} {pathname} />
	<div class="flex min-w-0 flex-1 min-h-0 flex-col">
		<Header {me} {breadcrumb} />
		<main class="flex-1 min-h-0 overflow-y-auto p-4 md:p-8">
			{@render children()}
		</main>
	</div>
</div>
