<script lang="ts">
	import { goto } from '$app/navigation';
	import CalendarDays from 'lucide-svelte/icons/calendar-days';
	import Users from 'lucide-svelte/icons/users';
	import FileText from 'lucide-svelte/icons/file-text';

	let { active }: { active: 'meetings' | 'pairs' | 'templates' } = $props();

	function go(href: string) {
		goto(href);
	}
</script>

<div role="tablist" class="tabs tabs-bordered mt-6 max-w-3xl">
	<button
		type="button"
		role="tab"
		class="tab gap-2"
		class:tab-active={active === 'meetings'}
		onclick={() => go('/planning/meetings')}
	>
		<CalendarDays class="h-4 w-4" />面談
	</button>
	<button
		type="button"
		role="tab"
		class="tab gap-2"
		class:tab-active={active === 'pairs'}
		onclick={() => go('/planning/pairs/users')}
	>
		<Users class="h-4 w-4" />ペア
	</button>
	<button
		type="button"
		role="tab"
		class="tab gap-2"
		class:tab-active={active === 'templates'}
		onclick={() => go('/planning/templates')}
	>
		<FileText class="h-4 w-4" />テンプレート
	</button>
</div>
