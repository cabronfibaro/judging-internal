<script lang="ts">
	import CommoncriteriaCard from '$lib/components/judging/commoncriteria-card.svelte';
	import CriteriaCard from '$lib/components/judging/criteria-card.svelte';
	import DebugJudgement from '$lib/components/judging/debug-judgement.svelte';
	import StageCard from '$lib/components/judging/stage-card.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import CardContent from '$lib/components/ui/card/card-content.svelte';
	import CardHeader from '$lib/components/ui/card/card-header.svelte';
	import CardTitle from '$lib/components/ui/card/card-title.svelte';
	import Card from '$lib/components/ui/card/card.svelte';
	import { Stage } from '$lib/judging/stage.svelte';
	import { createMutation } from '@tanstack/svelte-query';
	import type { Entrant, GroupedMatches } from '../../../types/tournament.js';
	import { updateTournament } from '$lib/api/tournaments.js';
	import { Judgement } from '$lib/judging/judgement.svelte.js';

	let { data } = $props();

	const tournament = $derived(data.tournament);
	const criteriaConfig = $derived(data.criteriaConfig);
	const entrants = $derived(data.entrants);
	let matches = $state<GroupedMatches>([]);

	let stage = new Stage();
	let judgement = new Judgement();

	let initialized = $state(false);
	$effect(() => {
		if (initialized) {
			return;
		}

		matches = data.matches;
		stage.tournamentId = tournament.id;
		if (tournament.status !== 'qualifying') return;

		stage.enqueueEntrants(entrants, tournament.qualifyingOrder);
		initialized = true;
	});

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
			const response = await fetch(`/api/tournaments/${stage.tournamentId}/entrants/${entrantId}`, {
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

		onSuccess: ({ entrant }: { entrant: Entrant }) => {
			stage.frontQueue().score = entrant.score;
			stage.frontQueue().status = entrant.status;
			stage.move(0, stage.queue.length - 1); // move front to back
		}
	}));

	const updateTournamentMutation = createMutation(() => ({
		mutationFn: async ({
			tournamentId,
			qualifyingOrder
		}: {
			tournamentId: string;
			qualifyingOrder: number[];
		}) => updateTournament(tournamentId, { qualifyingOrder: qualifyingOrder.join(', ') })
	}));

	async function onSubmitJudgement() {
		const { lead, chase } = judgement.finalizedScore();
		stage.setScore({ lead, chase });
		judgement.clear();

		if (!stage.isQualifying) {
			handleBattle();
			return;
		}

		const front = stage.frontQueue();

		await updateEntrantMutation.mutateAsync({
			entrantId: front.id,
			score: lead,
			status: 'scored'
		});

		await updateTournamentMutation.mutateAsync({
			tournamentId: stage.tournamentId,
			qualifyingOrder: stage.orders
		});
	}

	function handleBattle() {
		if (stage.isConcluded()) {
			console.log(stage.winner());
			return;
		}

		stage.switchLeader();
	}

	const categories = $derived(Object.entries(criteriaConfig.categories));
	let common = $derived(categories.find(([id]) => id === 'common'));
	let lead = $derived(categories.find(([id]) => id === 'lead'));
	let chase = $derived(categories.find(([id]) => id === 'chase'));
</script>

<div class="mx-auto max-w-4xl space-y-4 p-6">
	<StageCard {stage} bind:matches />
	<div class="grid grid-cols-2 gap-4">
		{#if common}
			{@const [, category] = common}
			<div class:col-span-2={!stage.isQualifying}>
				<CommoncriteriaCard {category} {judgement} isQualifying={stage.isQualifying} />
			</div>
		{/if}

		{#if lead}
			{@const [, category] = lead}
			<CriteriaCard {category} role="lead" {judgement} />
		{/if}

		{#if chase && !stage.isQualifying}
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
</div>
