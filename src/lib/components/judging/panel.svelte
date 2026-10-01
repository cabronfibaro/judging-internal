<script lang="ts">
	import { Judgement } from '$lib/judging/judgement.svelte';
	import { createMutation, useQueryClient } from '@tanstack/svelte-query';
	import type { CriteriaCatrgory } from '../../../types';
	import type { Entrant, Tournament } from '../../../types/tournament';
	import Button from '../ui/button/button.svelte';
	import CardContent from '../ui/card/card-content.svelte';
	import CardHeader from '../ui/card/card-header.svelte';
	import CardTitle from '../ui/card/card-title.svelte';
	import Card from '../ui/card/card.svelte';
	import CommoncriteriaCard from './commoncriteria-card.svelte';
	import CriteriaCard from './criteria-card.svelte';
	import DebugJudgement from './debug-judgement.svelte';
	import { updateTournament } from '$lib/api/tournaments';

	let {
		tournament = $bindable(),
		entrants,
		categories,
		isQualifying
	}: {
		tournament: Tournament;
		entrants: Entrant[];
		categories: [string, CriteriaCatrgory][];
		isQualifying: boolean;
	} = $props();

	let common = $derived(categories.find(([id]) => id === 'common'));
	let lead = $derived(categories.find(([id]) => id === 'lead'));
	let chase = $derived(categories.find(([id]) => id === 'chase'));

	let judgement = new Judgement();

	const queryClient = useQueryClient();

	const updateTournamentMutation = createMutation(() => ({
		mutationFn: async ({
			tournamentId,
			qualifyingOrder
		}: {
			tournamentId: string;
			qualifyingOrder: string[];
		}) =>
			updateTournament(tournamentId, {
				qualifyingOrder: qualifyingOrder.join(', ')
			}),

		onSuccess: (_, variables) => {
			tournament.qualifyingOrder = variables.qualifyingOrder.join(', ');
		}
	}));

	const updateEntrantMutation = createMutation(() => ({
		mutationFn: async ({
			entrantId,
			score,
			status
		}: {
			entrantId: string;
			score: number;
			status: Entrant['status'];
		}) => {
			const response = await fetch(`/api/tournaments/${tournament.id}/entrants/${entrantId}`, {
				method: 'PATCH',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ score, status })
			});

			const data = await response.json();

			if (!response.ok) {
				throw new Error(data.message ?? 'Failed to update entrant score.');
			}

			return data;
		},

		onSuccess: async () => {
			await queryClient.invalidateQueries({
				queryKey: ['entrants', tournament.id]
			});
		}
	}));

	const orderedEntrants = $derived.by(() => {
		const order = tournament.qualifyingOrder.split(',').map((id) => id.trim());
		const byParticipantId = new Map(entrants.map((entrant) => [entrant.participantId, entrant]));
		return order
			.map((id) => byParticipantId.get(id))
			.filter((entrant): entrant is Entrant => entrant !== undefined);
	});

	async function onSubmitJudgement() {
		const { lead } = judgement.finalizedScore();

		if (tournament.status === 'elimination') {
			return;
		}

		const front = orderedEntrants[0];

		await updateEntrantMutation.mutateAsync({
			entrantId: front.id,
			score: lead,
			status: 'scored'
		});

		const order = tournament.qualifyingOrder.split(',').map((id) => id.trim());

		const first = order.shift();
		if (first) {
			order.push(first);
		}

		await updateTournamentMutation.mutateAsync({
			tournamentId: tournament.id,
			qualifyingOrder: order
		});

		judgement.clear();
	}
</script>

<div class="grid grid-cols-2 gap-4">
	{#if common}
		{@const [, category] = common}
		<div class:col-span-2={!isQualifying}>
			<CommoncriteriaCard {category} {judgement} {isQualifying} />
		</div>
	{/if}

	{#if lead}
		{@const [, category] = lead}
		<CriteriaCard {category} role="lead" {judgement} />
	{/if}

	{#if chase && !isQualifying}
		{@const [, category] = chase}
		<CriteriaCard {category} role="chase" {judgement} />
	{/if}
</div>

<Card>
	<CardHeader class="pb-3">
		<CardTitle>Judgement</CardTitle>
	</CardHeader>

	<CardContent class="space-y-3">
		<DebugJudgement show results={judgement.results} />
		<Button class="w-full" onclick={onSubmitJudgement}>Submit Judgment</Button>
		<Button variant="destructive" type="reset" class="w-full" onclick={() => judgement.clear()}
			>Reset</Button
		>
	</CardContent>
</Card>
