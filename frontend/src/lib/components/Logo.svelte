<script lang="ts">
	let { class: className = 'h-5 w-5' }: { class?: string } = $props();
</script>

<svg class={className} viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
	<path
		d="M12 2L4 6v6c0 5 3.5 9.5 8 10 4.5-.5 8-5 8-10V6l-8-4z"
		stroke="currentColor"
		stroke-width="1.5"
		stroke-linejoin="round"
	/>
	<path d="M9 12l2 2 4-4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
</svg>
