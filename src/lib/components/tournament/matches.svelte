<script lang="ts">
	/* eslint-disable svelte/require-each-key */
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Tabs, TabsContent, TabsList, TabsTrigger } from '$lib/components/ui/tabs';
	import { capitalize } from 'es-toolkit/string';
	import type { GroupedMatches, Match } from '../../../types/tournament';
	import { Stage } from '$lib/judging/stage.svelte';
	import type { Judgement } from '$lib/judging/judgement.svelte';
	import Button from '../ui/button/button.svelte';
	import type { Battles } from '$lib/api/tournaments';

	let {
		matches,
		judgement,
		stage,
		activeBattles
	}: {
		matches: GroupedMatches;
		judgement: Judgement;
		stage: Stage;
		activeBattles: Battles;
	} = $props();

	let selectedMatchId = $state<string | undefined>(undefined);

	function getLead(p1Battle: Battles[number], p2Battle: Battles[number]): 'p1' | 'p2' {
		if (p1Battle.leadStatus === -1) return 'p2';
		if (p2Battle.leadStatus === -1) return 'p1';

		return p1Battle.leadStatus === 1 ? 'p1' : 'p2';
	}

	function onStartMatch(match: Match[]) {
		stage.clear();
		judgement.clear();

		selectedMatchId = match[0].matchId;
		stage.matchId = selectedMatchId;

		const p1 = match[0];
		const p2 = match[1];

		if (!p1.participantId || !p2.participantId) {
			console.warn('missing participant id');
			return;
		}
		const battles = activeBattles.filter((battle) => battle.matchId === p1.matchId);

		if (battles.length !== 2) {
			console.warn('active battle not found or incomplete');
			return;
		}

		const p1Battle = battles.find((battle) => battle.participantId === p1.participantId);
		const p2Battle = battles.find((battle) => battle.participantId === p2.participantId);

		if (!p1Battle || !p2Battle) {
			console.warn('missing active battle participant');
			return;
		}

		const lead = getLead(p1Battle, p2Battle);

		if (!lead) {
			console.warn('no lead participant');
			return;
		}

		stage.setDrivers(p1.participantId, p2.participantId, lead);
	}

	const sortedMatches = $derived(
		[...matches]
			.sort((a, b) => {
				const roundA = a[0]?.[0]?.round ?? 0;
				const roundB = b[0]?.[0]?.round ?? 0;

				if (roundA === 0) return 1;
				if (roundB === 0) return -1;

				return roundA - roundB;
			})
			.map((round) =>
				round.map((match) => {
					const battles = activeBattles.filter((battle) => battle.matchId === match[0]?.matchId);

					const leadingBattle = battles.find((battle) => battle.leadStatus === 1);

					return {
						match,
						leadingParticipantId: leadingBattle?.participantId ?? null
					};
				})
			)
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
				{@const roundNumber = round[0]?.match[0]?.round ?? 0}

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
			{@const roundNumber = round[0]?.match[0]?.round ?? 0}
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
					{#each round as matchData}
						{@const match = matchData.match}
						{@const currentMatch = match[0]}

						<Card class={selectedMatchId === currentMatch.matchId ? 'border border-green-500' : ''}>
							<CardHeader class="pb-3">
								<div class="flex items-center justify-between">
									<CardTitle class="text-sm">
										Match {currentMatch.matchId}
									</CardTitle>

									{#if currentMatch.status !== 'open'}
										<Badge variant="secondary">
											{capitalize(currentMatch.status)}
										</Badge>
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

											{#if matchData.leadingParticipantId === entrant.participantId}
												<Badge variant="outline">Lead</Badge>
											{/if}
										</div>

										{#if entrant.winnerId !== null && entrant.winnerId === entrant.participantId}
											<Badge>Winner</Badge>
										{/if}
									</div>
								{/each}

								<div class="flex justify-end">
									<Button onclick={() => onStartMatch(match)}>Start</Button>
								</div>
							</CardContent>
						</Card>
					{/each}
				</div>
			</TabsContent>
		{/each}
	</Tabs>
</div>
