<script lang="ts">
	import type { Stage } from '$lib/judging/stage.svelte';
	import type { GroupedMatches } from '../../../types/tournament';
	import EntrantList from '../tournament/entrant-list.svelte';
	import Matches from '../tournament/matches.svelte';
	import CardContent from '../ui/card/card-content.svelte';
	import CardHeader from '../ui/card/card-header.svelte';
	import CardTitle from '../ui/card/card-title.svelte';
	import Card from '../ui/card/card.svelte';
	import Separator from '../ui/separator/separator.svelte';

	let {
		stage,
		matches = $bindable()
	}: {
		stage: Stage;
		matches: GroupedMatches;
	} = $props();
</script>

<Card>
	<CardHeader>
		<CardTitle>Stage</CardTitle>
	</CardHeader>
	<Separator />
	<CardContent class="space-y-4">
		{#if stage.isQualifying}
			<EntrantList {stage} bind:matches />
		{:else}
			<Matches {matches} />
		{/if}
	</CardContent>
</Card>
