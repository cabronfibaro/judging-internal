<script lang="ts">
	import type { Stage } from '$lib/judging/stage.svelte';
	import type { GroupedMatches } from '../../../types/tournament';
	import EntrantList from '../tournament/entrant-list.svelte';
	import Matches from '../tournament/matches.svelte';
	import CardContent from '../ui/card/card-content.svelte';
	import CardHeader from '../ui/card/card-header.svelte';
	import CardTitle from '../ui/card/card-title.svelte';
	import Card from '../ui/card/card.svelte';
	import Label from '../ui/label/label.svelte';
	import Separator from '../ui/separator/separator.svelte';
	import Switch from '../ui/switch/switch.svelte';

	let {
		stage,
		matches = $bindable<GroupedMatches>([])
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
		<div class="grid grid-cols-1 gap-4">
			<div class="flex items-center justify-between space-y-3">
				<div class="flex h-6 items-center justify-between">
					<Label>Qualifying Mode</Label>
				</div>
				<Switch
					disabled
					checked={stage.isQualifying}
					onCheckedChange={(checked) => {
						stage.isQualifying = checked;
					}}
				/>
			</div>
		</div>

		{#if stage.isQualifying}
			<EntrantList {stage} />
		{/if}

		{#if !stage.isQualifying}
			<Matches {matches} />
		{/if}
	</CardContent>
</Card>
