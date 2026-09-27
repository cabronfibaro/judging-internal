<script lang="ts">
	/* eslint-disable @typescript-eslint/no-unused-vars */
	import { createMutation, createQuery, useQueryClient } from '@tanstack/svelte-query';

	import { Button } from '$lib/components/ui/button';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
	import {
		Dialog,
		DialogContent,
		DialogDescription,
		DialogFooter,
		DialogHeader,
		DialogTitle
	} from '$lib/components/ui/dialog';
	import Skeleton from '../ui/skeleton/skeleton.svelte';
	import type { Stage } from '$lib/judging/stage.svelte';

	let { stage }: { stage: Stage } = $props();

	type Tournament = {
		id: string;
		name: string;
		startingDate: Date | string;
		track: string | null;
		qualifyingOrder: string;
	};

	const queryClient = useQueryClient();

	const tournamentsQuery = createQuery(() => ({
		queryKey: ['tournaments'],
		queryFn: async (): Promise<Tournament[]> => {
			const response = await fetch('/api/tournaments');

			if (!response.ok) {
				throw new Error('Failed to load tournaments.');
			}

			const data = await response.json();

			return data.tournaments;
		}
	}));

	const deleteTournamentMutation = createMutation(() => ({
		mutationFn: async (id: string) => {
			const response = await fetch(`/api/tournaments/${id}`, {
				method: 'DELETE'
			});

			const data = await response.json();

			if (!response.ok) {
				throw new Error(data.message ?? 'Failed to delete tournament.');
			}

			return data;
		},

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

	function deleteTournament() {
		if (!selectedTournament) return;

		deleteTournamentMutation.mutate(selectedTournament.id);
	}

	let selectedTournamentId = $state<string | null>(null);

	type Qualifier = {
		id: string;
		participantId: string;
		name: string;
		status: 'qualified' | 'not_qualified' | 'pending';
		score: number;
	};

	async function fetchEntrants(tournamentId: string): Promise<Qualifier[]> {
		const response = await fetch(`/api/tournaments/${tournamentId}/entrants`);

		if (!response.ok) {
			throw new Error('Failed to load entrants.');
		}

		const data = await response.json();

		return data.entrants;
	}

	// TODO: SHOW TOAST AFTER LOAD
	async function onLoadTournament(tournament: Tournament) {
		const entrants = await queryClient.query({
			queryKey: ['tournament-entrants', tournament.id],
			queryFn: () => fetchEntrants(tournament.id)
		});

		stage.tournamentId = tournament.id;

		const isQualifying = entrants.reduce(
			(result, qualifier) => result || qualifier.status === 'pending',
			false
		);

		if (isQualifying) {
			stage.enqueueEntrants(entrants, tournament.qualifyingOrder);
		}

		selectedTournamentId = tournament.id;
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
				onclick={deleteTournament}
			>
				{deleteTournamentMutation.isPending ? 'Deleting...' : 'Delete'}
			</Button>
		</DialogFooter>
	</DialogContent>
</Dialog>

{#if tournamentsQuery.isPending}
	<div class="grid grid-cols-2 gap-4">
		{#each Array(4) as _, idx (idx)}
			<Card>
				<CardHeader>
					<Skeleton class="h-6 w-48 animate-none" />
				</CardHeader>

				<CardContent class="space-y-4">
					<Skeleton class="aspect-video w-full animate-none rounded-md" />

					<div class="space-y-2">
						<Skeleton class="h-4 w-32 animate-none" />
						<Skeleton class="h-4 w-48 animate-none" />
					</div>

					<Skeleton class="h-9 w-20 animate-none" />
				</CardContent>
			</Card>
		{/each}
	</div>
{:else if tournamentsQuery.isError}
	<p class="text-destructive">
		{tournamentsQuery.error.message}
	</p>
{:else if tournamentsQuery.data.length === 0}
	<p class="text-muted-foreground">No tournaments found.</p>
{:else}
	<div class="grid grid-cols-2 gap-4">
		{#each tournamentsQuery.data as tournament (tournament.id)}
			<Card>
				<CardHeader>
					<CardTitle>{tournament.name}</CardTitle>
				</CardHeader>

				<CardContent class="space-y-4">
					<Skeleton class="aspect-video w-full animate-none rounded-md" />

					<div class="space-y-1 text-sm">
						<p>
							<span class="text-muted-foreground">Track:</span>
							{tournament.track ?? 'Unknown'}
						</p>

						<p>
							<span class="text-muted-foreground">Starting:</span>
							{new Date(tournament.startingDate).toLocaleString()}
						</p>
					</div>

					<div class="flex gap-2">
						<Button onclick={() => onLoadTournament(tournament)}>Load</Button>
						<Button variant="destructive" onclick={() => openDeleteDialog(tournament)}>
							Delete
						</Button>
					</div>
				</CardContent>
			</Card>
		{/each}
	</div>
{/if}
