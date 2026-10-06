<script lang="ts">
	import CreateTournamentDialog from '#lib/components/tournament/create-tournament-dialog.svelte';
	import TournamentGrid from '#lib/components/tournament/tournament-grid.svelte';

	let { data } = $props();
	const tournaments = $derived(data.tournaments);
	let createTournamentOpen = $state(false);
</script>

<div class="w-full space-y-8">
	<header class="flex flex-col gap-5 border-b border-border pb-6 sm:flex-row sm:items-end sm:justify-between">
		<div class="space-y-2">
			<p class="text-sm font-medium text-muted-foreground">Tournament management</p>
			<h1 class="text-3xl font-semibold tracking-tight">Tournaments</h1>
			<p class="max-w-xl text-sm text-muted-foreground">
				Create a competition, manage qualifying participants, and open its bracket when you are ready.
			</p>
		</div>
		<CreateTournamentDialog bind:open={createTournamentOpen} />
	</header>

	<TournamentGrid {tournaments} onCreateTournament={() => (createTournamentOpen = true)} />
</div>
