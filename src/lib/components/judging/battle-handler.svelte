<script lang="ts">
	import type { Judgement } from '$lib/judging/judgement.svelte';
	import Button from '../ui/button/button.svelte';
	import CardContent from '../ui/card/card-content.svelte';
	import CardHeader from '../ui/card/card-header.svelte';
	import CardTitle from '../ui/card/card-title.svelte';
	import Card from '../ui/card/card.svelte';
	import DebugJudgement from './debug-judgement.svelte';
	import type { Stage } from '$lib/judging/stage.svelte';
	import { useQueryClient } from '@tanstack/svelte-query';
	import { getJson } from '$lib/utils';

	let {
		judgement,
		stage
	}: {
		judgement: Judgement;
		stage: Stage;
	} = $props();

	const queryClient = useQueryClient();
	let submitting = $state(false);

	async function onSubmitJudgement() {
		if (submitting) return;

		judgement.processing = true;
		submitting = true;

		try {
			const scores = judgement.finalizedScore();
			stage.setScore(scores);
			const [p1, p2] = stage.getScores();

			await Promise.all([
				getJson(fetch, `/api/tournaments/${stage.tournamentId}/bracket/matches/${stage.matchId}`, {
					method: 'PATCH',
					headers: {
						'Content-Type': 'application/json'
					},
					body: JSON.stringify({
						...p1,
						appendScore: true
					})
				}),
				getJson(fetch, `/api/tournaments/${stage.tournamentId}/bracket/matches/${stage.matchId}`, {
					method: 'PATCH',
					headers: {
						'Content-Type': 'application/json'
					},
					body: JSON.stringify({
						...p2,
						appendScore: true
					})
				})
			]);

			await getJson(
				fetch,
				`/api/tournaments/${stage.tournamentId}/bracket/matches/${stage.matchId}/changeLeader`
			);

			await queryClient.invalidateQueries({
				queryKey: ['activeBattles', stage.tournamentId]
			});

			await queryClient.invalidateQueries({
				queryKey: ['matches', stage.tournamentId]
			});

			stage.switchLeader();
			judgement.clear();
		} finally {
			submitting = false;
			judgement.processing = false;
		}
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
