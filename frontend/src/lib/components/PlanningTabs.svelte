<script lang="ts">
	let { active }: { active: 'meetings' | 'pairs' | 'templates' } = $props();
</script>

<div class="tabs tabs-boxed w-fit">
	<a href="/planning/meetings" class="tab" class:tab-active={active === 'meetings'}>面談</a>
	<a href="/planning/pairs" class="tab" class:tab-active={active === 'pairs'}>ペア</a>
	<a href="/planning/templates" class="tab" class:tab-active={active === 'templates'}>テンプレート</a>
</div>
