<script lang="ts">
	import type { PairingSection } from '$lib/pairing-helpers';

	let { active }: { active: PairingSection } = $props();
</script>

<div class="tabs tabs-boxed w-fit">
	<a href="/planning/pairs/users" class="tab" class:tab-active={active === 'users'}>ユーザー</a>
	<a href="/planning/pairs/departments" class="tab" class:tab-active={active === 'departments'}>部署</a>
	<a href="/planning/pairs/teams" class="tab" class:tab-active={active === 'teams'}>チーム</a>
</div>
