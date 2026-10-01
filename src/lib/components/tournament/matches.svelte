<script lang="ts">
	/* eslint-disable svelte/require-each-key */
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Tabs, TabsContent, TabsList, TabsTrigger } from '$lib/components/ui/tabs';
	import { capitalize } from 'es-toolkit/string';
	import type { Entrant, GroupedMatches, Match } from '../../../types/tournament';
	import { Stage } from '$lib/judging/stage.svelte';
	import type { Judgement } from '$lib/judging/judgement.svelte';

	let {
		matches,
		judgement,
		entrants,
		stage
	}: { matches: GroupedMatches; judgement: Judgement; entrants: Entrant[]; stage: Stage } =
		$props();

	let selectedMatchId = $state<string | undefined>(undefined);

	function onSelectMatch(match: Match[]) {
		stage.clear();
		judgement.clear();

		selectedMatchId = match[0].matchId;

		const p1 = match[0];
		const p2 = match[1];

		// only select match that already have 2 participants
		if (!p1.participantId || !p2.participantId) {
			console.warn('missing participant id');
			return;
		}

		const p1Entrant = entrants.find((entrant) => entrant.participantId === p1.participantId);
		const p2Entrant = entrants.find((entrant) => entrant.participantId === p2.participantId);

		if (!p1Entrant || !p2Entrant) {
			console.warn('missing participant score');
			return;
		}

		const lead = p1Entrant.score > p2Entrant.score ? 'p1' : 'p2';

		stage.setDrivers(p1.participantId, p2.participantId, lead);
	}

	const sortedMatches = $derived(
		[...matches].sort((a, b) => {
			const roundA = a[0]?.[0]?.round ?? 0;
			const roundB = b[0]?.[0]?.round ?? 0;

			if (roundA === 0) return 1;
			if (roundB === 0) return -1;

			return roundA - roundB;
		})
	);
</script>

<div class="space-y-4">
	<div>
		<h3 class="font-semibold">Bracket</h3>
		<p class="text-sm text-muted-foreground">Drift battle matches by round</p>
	</div>

	<Tabs value="round-1">
		<TabsList
			class="grid w-full"
			style={`grid-template-columns: repeat(${sortedMatches.length}, minmax(0, 1fr))`}
		>
			{#each sortedMatches as round}
				{@const roundNumber = round[0]?.[0]?.round ?? 0}

				{#if roundNumber === 0}
					<TabsTrigger value="third-place">3rd Place</TabsTrigger>
				{:else}
					<TabsTrigger value={`round-${roundNumber}`}>
						Round {roundNumber}
					</TabsTrigger>
				{/if}
			{/each}
		</TabsList>

		{#each sortedMatches as round}
			{@const roundNumber = round[0]?.[0]?.round ?? 0}
			{@const tabValue = roundNumber === 0 ? 'third-place' : `round-${roundNumber}`}

			<TabsContent value={tabValue} class="mt-4 space-y-3">
				<div>
					<h4 class="font-medium">
						{roundNumber === 0 ? '3rd Place' : `Round ${roundNumber}`}
					</h4>

					<p class="text-sm text-muted-foreground">
						{round.length}
						{round.length === 1 ? 'battle' : 'battles'}
					</p>
				</div>

				<div class="grid gap-3">
					{#each round as match}
						{@const currentMatch = match[0]}
						<Card
							onclick={() => onSelectMatch(match)}
							class={selectedMatchId === currentMatch.matchId ? 'border border-green-500' : ''}
						>
							<CardHeader class="pb-3">
								<div class="flex items-center justify-between">
									<CardTitle class="text-sm">
										Match {match[0]?.matchId}
									</CardTitle>

									{#if currentMatch?.status !== 'open'}
										<Badge variant="secondary">{capitalize(currentMatch.status)}</Badge>
									{/if}
								</div>
							</CardHeader>

							<CardContent class="space-y-2">
								{#each match as entrant}
									<div
										class={[
											'flex items-center justify-between rounded-md border px-4 py-3',
											entrant.winnerId === entrant.participantId
												? 'bg-accent text-accent-foreground'
												: ''
										]}
									>
										<div class="flex items-center gap-3">
											<span class="font-medium">
												{entrant.participant_name ?? 'TBD'}
											</span>
										</div>

										{#if entrant.winnerId !== null && entrant.winnerId === entrant.participantId}
											<Badge>Winner</Badge>
										{/if}
									</div>
								{/each}
							</CardContent>
						</Card>
					{/each}
				</div>
			</TabsContent>
		{/each}
	</Tabs>
</div>
