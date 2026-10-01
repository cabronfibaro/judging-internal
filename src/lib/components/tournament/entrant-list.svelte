<script lang="ts">
	import EllipsisVertical from '@lucide/svelte/icons/ellipsis-vertical';
	import { Badge } from '$lib/components/ui/badge';
	import {
		DropdownMenu,
		DropdownMenuContent,
		DropdownMenuItem,
		DropdownMenuTrigger
	} from '$lib/components/ui/dropdown-menu';
	import { Button } from '$lib/components/ui/button';
	import { createMutation, useQueryClient } from '@tanstack/svelte-query';
	import type { Entrant, GroupedMatches, Tournament } from '../../../types/tournament';
	import { syncTournament, updateTournament } from '$lib/api/tournaments';
	import AlertDialog from '../ui/alert-dialog/alert-dialog.svelte';
	import AlertDialogTrigger from '../ui/alert-dialog/alert-dialog-trigger.svelte';
	import AlertDialogContent from '../ui/alert-dialog/alert-dialog-content.svelte';
	import AlertDialogHeader from '../ui/alert-dialog/alert-dialog-header.svelte';
	import AlertDialogTitle from '../ui/alert-dialog/alert-dialog-title.svelte';
	import AlertDialogDescription from '../ui/alert-dialog/alert-dialog-description.svelte';
	import AlertDialogFooter from '../ui/alert-dialog/alert-dialog-footer.svelte';
	import AlertDialogCancel from '../ui/alert-dialog/alert-dialog-cancel.svelte';
	import AlertDialogAction from '../ui/alert-dialog/alert-dialog-action.svelte';

	let {
		tournament = $bindable(),
		entrants,
		// eslint-disable-next-line no-useless-assignment
		matches = $bindable()
	}: {
		tournament: Tournament;
		entrants: Entrant[];
		matches: GroupedMatches;
	} = $props();

	const queryClient = useQueryClient();

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

	function moveToFirst(idx: number) {
		const [entrant] = orderedEntrants.splice(idx, 1);

		if (!entrant) return;

		const reordered = [entrant, ...orderedEntrants];

		updateTournamentMutation.mutate({
			tournamentId: tournament.id,
			qualifyingOrder: reordered.map((entrant) => entrant.participantId)
		});
	}

	const updateQualifierMutation = createMutation(() => ({
		mutationFn: async ({
			qualifierId,
			status
		}: {
			qualifierId: string;
			status: Entrant['status'];
		}) => {
			const response = await fetch(`/api/tournaments/${tournament.id}/entrants/${qualifierId}`, {
				method: 'PATCH',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ status })
			});

			const data = await response.json();

			if (!response.ok) {
				throw new Error(data.message ?? 'Failed to update qualifier.');
			}

			return data;
		},

		onSuccess: async () => {
			await queryClient.invalidateQueries({
				queryKey: ['entrants', tournament.id]
			});
		}
	}));

	function markAs(id: string, status: Entrant['status']) {
		updateQualifierMutation.mutate({
			qualifierId: id,
			status
		});
	}

	let concludeDialogOpen = $state(false);
	const concludeQualifiersMutation = createMutation(() => ({
		mutationFn: async () => {
			const response = await fetch(
				`/api/tournaments/${tournament.id}/entrants/concludeQualifiers`,
				{
					method: 'GET'
				}
			);

			const data = await response.json();

			if (!response.ok) {
				throw new Error(data.message ?? 'Failed to conclude qualifiers.');
			}

			return data;
		},

		onSuccess: async ({ tournament: newTournament }: { tournament: Tournament }) => {
			matches = await syncTournament(fetch, tournament.id);
			tournament = newTournament;
			await queryClient.invalidateQueries({
				queryKey: ['entrants', newTournament.id]
			});
			concludeDialogOpen = false;
		}
	}));

	const orderedEntrants = $derived.by(() => {
		const order = tournament.qualifyingOrder.split(',').map((id) => id.trim());
		const byParticipantId = new Map(entrants.map((entrant) => [entrant.participantId, entrant]));
		return order
			.map((id) => byParticipantId.get(id))
			.filter((entrant): entrant is Entrant => entrant !== undefined);
	});

	const front = $derived(orderedEntrants[0].participantId);

	const concluding = $derived.by(() => {
		return !entrants.some((entrant) => entrant.status === 'pending');
	});
</script>

<div class="space-y-3">
	<div>
		<h3 class="font-semibold">Entrants</h3>
		<p class="text-sm text-muted-foreground">Current qualification scores</p>
	</div>

	<div class="divide-y rounded-md border">
		{#each orderedEntrants as entrant, idx (entrant.id)}
			<button
				type="button"
				class={[
					'flex w-full items-center justify-between px-4 py-3 text-left transition-colors',
					front === entrant.participantId ? 'bg-accent text-accent-foreground' : 'hover:bg-muted/50'
				]}
			>
				<div class="flex items-center gap-3">
					<span class="font-medium">
						{entrant.name}
					</span>

					{#if entrant.status === 'pending'}
						{#if idx === 0}
							<Badge variant="default">On Going</Badge>
						{/if}
					{:else if entrant.status === 'qualified'}
						<Badge>Qualified</Badge>
					{:else if entrant.status === 'not_qualified'}
						<Badge variant="destructive">Not Qualified</Badge>
					{/if}
				</div>

				<div class="flex items-center gap-3">
					<span class="text-lg font-semibold tabular-nums">
						{entrant.status === 'pending' ? 'TBA' : entrant.score}
					</span>

					<!-- svelte-ignore a11y_click_events_have_key_events -->
					<!-- svelte-ignore a11y_no_static_element_interactions -->
					<div onclick={(event) => event.stopPropagation()}>
						<DropdownMenu>
							<DropdownMenuTrigger>
								<Button variant="ghost" size="icon" aria-label="Qualifier actions">
									<EllipsisVertical />
								</Button>
							</DropdownMenuTrigger>

							<DropdownMenuContent align="end" class="w-48">
								<DropdownMenuItem
									onclick={() => moveToFirst(idx)}
									disabled={entrant.status === 'not_qualified' ||
										updateTournamentMutation.isPending}
								>
									Move to front
								</DropdownMenuItem>

								<DropdownMenuItem
									onclick={() => markAs(entrant.id, 'not_qualified')}
									disabled={updateQualifierMutation.isPending}
								>
									Set as not qualified
								</DropdownMenuItem>

								<DropdownMenuItem
									onclick={() => markAs(entrant.id, 'pending')}
									disabled={updateQualifierMutation.isPending}
								>
									Set as pending
								</DropdownMenuItem>
							</DropdownMenuContent>
						</DropdownMenu>
					</div>
				</div>
			</button>
		{/each}
	</div>

	{#if concluding}
		<AlertDialog bind:open={concludeDialogOpen}>
			<AlertDialogTrigger class="w-full">
				<Button class="w-full" variant="default" disabled={concludeQualifiersMutation.isPending}>
					Conclude
				</Button>
			</AlertDialogTrigger>

			<AlertDialogContent>
				<AlertDialogHeader>
					<AlertDialogTitle>Conclude qualifying?</AlertDialogTitle>

					<AlertDialogDescription>
						After concluding qualifying, the qualified entrants will be finalized and will be seeded
						to the tournament bracket. This action cannot be undone.
					</AlertDialogDescription>
				</AlertDialogHeader>

				<AlertDialogFooter>
					<AlertDialogCancel>Cancel</AlertDialogCancel>

					<AlertDialogAction
						onclick={() => {
							concludeQualifiersMutation.mutate();
						}}
						disabled={concludeQualifiersMutation.isPending}
					>
						{concludeQualifiersMutation.isPending ? 'Concluding...' : 'Conclude'}
					</AlertDialogAction>
				</AlertDialogFooter>
			</AlertDialogContent>
		</AlertDialog>
	{/if}
</div>
