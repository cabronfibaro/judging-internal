<script lang="ts">
	import { createMutation, createQuery, useQueryClient } from '@tanstack/svelte-query';

	import { Button } from '$lib/components/ui/button';
	import {
		Dialog,
		DialogContent,
		DialogDescription,
		DialogFooter,
		DialogHeader,
		DialogTitle
	} from '$lib/components/ui/dialog';

	import type { Tournament } from '../../../types/tournament';

	import { deleteTournament, getTournaments } from '$lib/api/tournaments';
	import TournamentCardSkel from './tournament-card-skel.svelte';
	import TournamentCard from './tournament-card.svelte';

	let { tournaments }: { tournaments: Tournament[] } = $props();

	const queryClient = useQueryClient();

	const tournamentsQuery = createQuery(() => ({
		queryKey: ['tournaments'],
		queryFn: () => getTournaments(),
		initialData: tournaments
	}));

	const deleteTournamentMutation = createMutation(() => ({
		mutationFn: ({ id, status }: { id: string; status: Tournament['status'] }) =>
			deleteTournament(id, status),

		onSuccess: async () => {
			await queryClient.invalidateQueries({
				queryKey: ['tournaments']
			});

			deleteDialogOpen = false;
			selectedTournament = null;
		}
	}));
	let deleteDialogOpen = $state(false);
	let selectedTournament = $state<Tournament | null>(null);

	function openDeleteDialog(tournament: Tournament) {
		selectedTournament = tournament;
		deleteTournamentMutation.reset();
		deleteDialogOpen = true;
	}

	function onDeleteTournament() {
		if (!selectedTournament) return;
		const { id, status } = selectedTournament;
		deleteTournamentMutation.mutate({ id, status });
	}
</script>

<Dialog bind:open={deleteDialogOpen}>
	<DialogContent>
		<DialogHeader>
			<DialogTitle>Delete tournament?</DialogTitle>

			<DialogDescription>
				This will permanently delete
				<strong>{selectedTournament?.name}</strong>
				and its qualifier records. This action cannot be undone.
			</DialogDescription>
		</DialogHeader>

		{#if deleteTournamentMutation.isError}
			<p class="text-sm text-destructive">
				{deleteTournamentMutation.error.message}
			</p>
		{/if}

		<DialogFooter>
			<Button
				variant="outline"
				disabled={deleteTournamentMutation.isPending}
				onclick={() => (deleteDialogOpen = false)}
			>
				Cancel
			</Button>

			<Button
				variant="destructive"
				disabled={deleteTournamentMutation.isPending}
				onclick={onDeleteTournament}
			>
				{deleteTournamentMutation.isPending ? 'Deleting...' : 'Delete'}
			</Button>
		</DialogFooter>
	</DialogContent>
</Dialog>

<div class="grid grid-cols-2 gap-4">
	{#if tournamentsQuery.isPending}
		<div class="grid grid-cols-2 gap-4">
			{#each Array.from({ length: 4 }, (_, i) => i) as idx (idx)}
				<TournamentCardSkel />
			{/each}
		</div>
	{:else if tournamentsQuery.isError}
		<p class="text-destructive">
			{tournamentsQuery.error.message}
		</p>
	{:else if tournamentsQuery.data.length === 0}
		<p class="text-muted-foreground">No tournaments found.</p>
	{:else}
		{#each tournamentsQuery.data as tournament (tournament.id)}
			<TournamentCard {tournament} {openDeleteDialog} />
		{/each}
	{/if}
</div>
