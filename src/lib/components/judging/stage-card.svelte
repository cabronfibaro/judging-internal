<script lang="ts">
	import type { Snippet } from 'svelte';
	import CardContent from '../ui/card/card-content.svelte';
	import CardHeader from '../ui/card/card-header.svelte';
	import CardTitle from '../ui/card/card-title.svelte';
	import Card from '../ui/card/card.svelte';
	import Separator from '../ui/separator/separator.svelte';

	let {
		children
	}: {
		children: Snippet<[]>;
	} = $props();
</script>

<Card class="gap-0 border-border/70 shadow-none">
	<CardHeader class="pb-4">
		<CardTitle>Stage</CardTitle>
	</CardHeader>
	<Separator />
	<CardContent class="space-y-4">
		{@render children()}
	</CardContent>
</Card>
