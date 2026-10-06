<script lang="ts">
	import { createMutation, createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { AlertCircle, RefreshCw, Trophy, Plus } from '@lucide/svelte';

	import { Button } from '#lib/components/ui/button';
	import {
		Empty,
		EmptyDescription,
		EmptyHeader,
		EmptyMedia,
		EmptyTitle
	} from '#lib/components/ui/empty';
	import {
		Dialog,
		DialogContent,
		DialogDescription,
		DialogFooter,
		DialogHeader,
		DialogTitle
	} from '#lib/components/ui/dialog';
	import { Badge } from '#lib/components/ui/badge';

	import type { Tournament } from '../../../types/tournament';

	import { deleteTournament, getTournaments } from '#lib/api/tournaments';
	import TournamentCardSkel from './tournament-card-skel.svelte';
	import TournamentCard from './tournament-card.svelte';

	let {
		tournaments,
		onCreateTournament
	}: {
		tournaments: Tournament[];
		onCreateTournament: () => void;
	} = $props();
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

	const tournamentCountLabel = $derived(
		tournamentsQuery.data.length === 1 ? '1 tournament' : `${tournamentsQuery.data.length} tournaments`
	);
	const activeCount = $derived(
		tournamentsQuery.data.filter((tournament) => tournament.status === 'elimination').length
	);
	const qualifyingCount = $derived(
		tournamentsQuery.data.filter((tournament) => tournament.status === 'qualifying').length
	);
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

{#if tournamentsQuery.isPending}
	<div class="w-full space-y-4" aria-label="Loading tournaments">
		{#each Array.from({ length: 4 }, (_, i) => i) as idx (idx)}
			<TournamentCardSkel />
		{/each}
	</div>
{:else if tournamentsQuery.isError}
	<Empty class="min-h-64 border border-dashed border-destructive/40">
		<EmptyHeader>
			<EmptyMedia variant="icon"><AlertCircle aria-hidden="true" /></EmptyMedia>
			<EmptyTitle>We couldn't load tournaments</EmptyTitle>
			<EmptyDescription>
				{tournamentsQuery.error.message || 'The tournament list could not be loaded.'} Try again to refresh the list.
			</EmptyDescription>
		</EmptyHeader>
		<Button variant="outline" onclick={() => tournamentsQuery.refetch()}>
			<RefreshCw aria-hidden="true" />
			Try again
		</Button>
	</Empty>
{:else if tournamentsQuery.data.length === 0}
	<Empty class="min-h-80 border border-dashed">
		<EmptyHeader>
			<EmptyMedia variant="icon"><Trophy aria-hidden="true" /></EmptyMedia>
			<EmptyTitle>Your competition board is empty</EmptyTitle>
			<EmptyDescription>
				Create your first tournament to start qualifying drivers and running head-to-head battles.
			</EmptyDescription>
		</EmptyHeader>
		<Button onclick={onCreateTournament}>
			<Plus aria-hidden="true" />
			Create tournament
		</Button>
	</Empty>
{:else}
	<div class="w-full space-y-6">
		<div class="flex flex-col gap-4 border-b border-border pb-5 sm:flex-row sm:items-end sm:justify-between">
			<div class="min-w-0">
				<div class="flex flex-wrap items-center gap-3">
					<h2 class="text-lg font-semibold tracking-tight">Competition board</h2>
					<Badge variant="outline">{tournamentCountLabel}</Badge>
				</div>
				<p class="mt-1 max-w-2xl text-sm text-muted-foreground">
					A live view of every tournament and its bracket progress.
				</p>
			</div>
			<div class="flex shrink-0 flex-wrap items-center gap-2 text-xs text-muted-foreground">
				<span class="rounded-full bg-primary/10 px-2.5 py-1 text-primary">{activeCount} active</span>
				<span class="rounded-full bg-muted px-2.5 py-1">{qualifyingCount} qualifying</span>
			</div>
		</div>

		<div class="w-full space-y-3">
			{#each tournamentsQuery.data as tournament, index (tournament.id)}
				<TournamentCard {tournament} {openDeleteDialog} position={index + 1} />
			{/each}
		</div>
	</div>
{/if}
