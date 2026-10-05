<script lang="ts">
	import type { Judgement } from '#lib/judging/judgement.svelte';
	let {
		results,
		show = false
	}: { results: InstanceType<typeof Judgement>['results']; show: boolean } = $props();
</script>

{#if show}
	<pre class="overflow-auto rounded-md bg-muted p-3 text-sm">{JSON.stringify(
			results,
			null,
			2
		)}</pre>
{/if}
