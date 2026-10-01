<script lang="ts">
	import StageCard from '$lib/components/judging/stage-card.svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import type { Tournament } from '../../../types/tournament.js';
	import { getEntrants, getMatches } from '$lib/api/tournaments.js';
	import EntrantList from '$lib/components/tournament/entrant-list.svelte';
	import Matches from '$lib/components/tournament/matches.svelte';
	import Panel from '$lib/components/judging/panel.svelte';
	import { Judgement } from '$lib/judging/judgement.svelte.js';
	import { Stage } from '$lib/judging/stage.svelte.js';
	import QualifyingHandler from '$lib/components/judging/qualifying-handler.svelte';
	import BattleHandler from '$lib/components/judging/battle-handler.svelte';

	let { data } = $props();
	const criteriaConfig = $derived(data.criteriaConfig);

	let stage = new Stage();
	let judgement = new Judgement();
	let tournament = $state<Tournament>({} as Tournament);

	const entrantsQuery = createQuery(() => ({
		queryKey: ['entrants', tournament.id],
		queryFn: async () => {
			return getEntrants(fetch, tournament.id);
		},
		initialData: data.entrants
	}));

	const matchesQuery = createQuery(() => ({
		queryKey: ['matches', tournament.id],
		queryFn: async () => {
			return getMatches(tournament.id);
		},
		initialData: data.matches
	}));

	let initialized = $state(false);
	$effect(() => {
		if (initialized) {
			return;
		}

		tournament = data.tournament;
		stage.tournamentId = tournament.id;
		initialized = true;
	});

	const categories = $derived(Object.entries(criteriaConfig.categories));
	const isQualifying = $derived(tournament.status === 'qualifying');
</script>

<div class="mx-auto max-w-4xl space-y-4 p-6">
	<StageCard>
		{#if isQualifying}
			<EntrantList bind:tournament entrants={entrantsQuery.data} />
		{:else}
			<Matches {stage} matches={matchesQuery.data} entrants={entrantsQuery.data} {judgement} />
		{/if}
	</StageCard>
	<Panel {judgement} {categories} {isQualifying}>
		{#if isQualifying}
			<QualifyingHandler bind:tournament {judgement} entrants={entrantsQuery.data} />
		{:else}
			<BattleHandler {judgement} {stage} />
		{/if}
	</Panel>
</div>
