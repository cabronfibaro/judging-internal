<script lang="ts">
	import type { Stage } from '$lib/judging/stage.svelte';
	import Badge from '../ui/badge/badge.svelte';
	import CardContent from '../ui/card/card-content.svelte';
	import CardHeader from '../ui/card/card-header.svelte';
	import CardTitle from '../ui/card/card-title.svelte';
	import Card from '../ui/card/card.svelte';
	import Input from '../ui/input/input.svelte';
	import Label from '../ui/label/label.svelte';

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

	<CardContent class="space-y-4">
		<div class="grid grid-cols-2 gap-4">
			<!-- P1 -->
			<div class="space-y-3">
				<div class="flex items-center justify-between">
					<Label>P1</Label>

					<Badge variant={stage.state.p1.state === 'lead' ? 'default' : 'secondary'}>
						{stage.state.p1.state}
					</Badge>
				</div>
				<Input
					placeholder="Driver name"
					value={stage.state.p1.name}
					oninput={(event) => stage.updateDriver('p1', 'name', event.currentTarget.value)}
				/>
			</div>

			<!-- P2 -->
			<div class="space-y-3">
				<div class="flex items-center justify-between">
					<Label>P2</Label>

					<Badge variant={stage.state.p2.state === 'lead' ? 'default' : 'secondary'}>
						{stage.state.p2.state}
					</Badge>
				</div>

				<Input
					placeholder="Driver name"
					value={stage.state.p2.name}
					oninput={(event) => stage.updateDriver('p2', 'name', event.currentTarget.value)}
				/>
			</div>
		</div>
	</CardContent>
</Card>
