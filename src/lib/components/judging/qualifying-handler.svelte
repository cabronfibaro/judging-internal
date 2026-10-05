<script lang="ts">
	import type { Judgement } from '#lib/judging/judgement.svelte';
	import { createMutation, useQueryClient } from '@tanstack/svelte-query';
	import Button from '../ui/button/button.svelte';
	import CardContent from '../ui/card/card-content.svelte';
	import CardHeader from '../ui/card/card-header.svelte';
	import CardTitle from '../ui/card/card-title.svelte';
	import Card from '../ui/card/card.svelte';
	import DebugJudgement from './debug-judgement.svelte';
	import { updateTournament } from '#lib/api/tournaments';
	import type { Entrant, Tournament } from '../../../types/tournament';

	let {
		tournament = $bindable(),
		judgement,
		entrants
	}: {
		tournament: Tournament;
		judgement: Judgement;
		entrants: Entrant[];
	} = $props();

	const queryClient = useQueryClient();

	let submitting = $state(false);

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
		if (submitting) return;

		judgement.processing = true;
		submitting = true;

		try {
			const { lead } = judgement.finalizedScore();

			await handleQualifying(lead);

			judgement.clear();
		} finally {
			submitting = false;
			judgement.processing = false;
		}
	}

	async function handleQualifying(score: number) {
		const front = orderedEntrants[0];

		await updateEntrantMutation.mutateAsync({
			entrantId: front.id,
			score,
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
	}
</script>

<Card>
	<CardHeader class="pb-3">
		<CardTitle>Judgement</CardTitle>
	</CardHeader>

	<CardContent class="space-y-3">
		<DebugJudgement show results={judgement.results} />

		<Button class="w-full" onclick={onSubmitJudgement} disabled={submitting}>
			{#if submitting}
				Submitting...
			{:else}
				Submit Judgment
			{/if}
		</Button>

		<Button
			variant="destructive"
			type="reset"
			class="w-full"
			onclick={() => judgement.clear()}
			disabled={submitting}
		>
			Reset
		</Button>
	</CardContent>
</Card>
