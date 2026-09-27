<script lang="ts">
	import type { Stage } from '$lib/judging/stage.svelte';
	import EntrantList from '../tournament/entrant-list.svelte';
	import Badge from '../ui/badge/badge.svelte';
	import CardContent from '../ui/card/card-content.svelte';
	import CardHeader from '../ui/card/card-header.svelte';
	import CardTitle from '../ui/card/card-title.svelte';
	import Card from '../ui/card/card.svelte';
	import Input from '../ui/input/input.svelte';
	import Label from '../ui/label/label.svelte';
	import Separator from '../ui/separator/separator.svelte';
	import Switch from '../ui/switch/switch.svelte';

	let {
		stage
	}: {
		stage: Stage;
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
			<div class="grid grid-cols-2 items-start gap-4">
				<!-- P1 -->
				<div class="space-y-3">
					<div class="flex h-6 items-center justify-between">
						<Label>P1</Label>

						<div class:invisible={stage.state.p1.state !== 'lead'}>
							<Badge variant="default">lead</Badge>
						</div>
					</div>

					<Input
						placeholder="Driver name"
						value={stage.state.p1.name}
						oninput={(event) => stage.updateDriver('p1', 'name', event.currentTarget.value)}
					/>
				</div>

				<div class="space-y-3">
					<div class="flex h-6 items-center justify-between">
						<Label>P2</Label>

						<div class:invisible={stage.state.p2.state !== 'lead'}>
							<Badge variant="default">lead</Badge>
						</div>
					</div>

					<Input
						placeholder="Driver name"
						value={stage.state.p2.name}
						oninput={(event) => stage.updateDriver('p2', 'name', event.currentTarget.value)}
					/>
				</div>
			</div>
		{/if}
	</CardContent>
</Card>
