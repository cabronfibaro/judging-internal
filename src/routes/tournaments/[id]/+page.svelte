<script lang="ts">
	import StageCard from '$lib/components/judging/stage-card.svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import type { Entrant, GroupedMatches, Tournament } from '../../../types/tournament.js';
	import { getEntrants } from '$lib/api/tournaments.js';
	import EntrantList from '$lib/components/tournament/entrant-list.svelte';
	import Matches from '$lib/components/tournament/matches.svelte';
	import Panel from '$lib/components/judging/panel.svelte';

	let { data } = $props();
	const criteriaConfig = $derived(data.criteriaConfig);

	let tournament = $state<Tournament>({} as Tournament);

	const entrantsQuery = createQuery(() => ({
		queryKey: ['entrants', tournament.id],
		queryFn: async (): Promise<Entrant[]> => {
			return getEntrants(fetch, tournament.id);
		},
		initialData: data.entrants
	}));

	let matches = $state<GroupedMatches>([]);

	let initialized = $state(false);
	$effect(() => {
		if (initialized) {
			return;
		}

		tournament = data.tournament;
		matches = data.matches;
		initialized = true;
	});

	const categories = $derived(Object.entries(criteriaConfig.categories));
	const isQualifying = $derived(tournament.status === 'qualifying');
</script>

<div class="mx-auto max-w-4xl space-y-4 p-6">
	<StageCard>
		{#if isQualifying}
			<EntrantList bind:tournament entrants={entrantsQuery.data ?? []} bind:matches />
		{:else}
			<Matches {matches} />
		{/if}
	</StageCard>
	<Panel entrants={entrantsQuery.data ?? []} {isQualifying} {categories} bind:tournament />
</div>
