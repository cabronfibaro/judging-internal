<script lang="ts">
	import CreateTournamentDialog from '$lib/components/tournament/create-tournament-dialog.svelte';
	import TournamentGrid from '$lib/components/tournament/tournament-grid.svelte';

	let { data } = $props();
	const tournaments = $derived(data.tournaments);
</script>

<div class="mx-auto max-w-4xl space-y-4 p-6">
	<CreateTournamentDialog />
	<TournamentGrid {tournaments} />
</div>
