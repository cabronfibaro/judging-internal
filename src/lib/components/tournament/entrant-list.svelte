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
	import type { Stage } from '$lib/judging/stage.svelte';
	import { createMutation } from '@tanstack/svelte-query';
	import type { Entrant, Tournament } from '../../../types/tournament';
	import { updateTournament } from '$lib/api/tournaments';

	let { stage }: { stage: Stage } = $props();

	const updateTournamentMutation = createMutation(() => ({
		mutationFn: async ({
			tournamentId,
			qualifyingOrder
		}: {
			tournamentId: string;
			qualifyingOrder: number[];
		}) => updateTournament(tournamentId, qualifyingOrder)
	}));

	function moveToFirst(idx: number) {
		stage.move(idx, 0);
		updateTournamentMutation.mutate({
			tournamentId: stage.tournamentId,
			qualifyingOrder: stage.orders
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
			const response = await fetch(
				`/api/tournaments/${stage.tournamentId}/entrants/${qualifierId}`,
				{
					method: 'PATCH',
					headers: {
						'Content-Type': 'application/json'
					},
					body: JSON.stringify({ status })
				}
			);

			const data = await response.json();

			if (!response.ok) {
				throw new Error(data.message ?? 'Failed to update qualifier.');
			}

			return data;
		},

		onSuccess: ({ qualifier }) => {
			stage.updateEntrants(qualifier);
		}
	}));

	function markAs(id: string, status: Entrant['status']) {
		updateQualifierMutation.mutate({
			qualifierId: id,
			status: status
		});
	}
	const concludeQualifiersMutation = createMutation(() => ({
		mutationFn: async () => {
			const response = await fetch(
				`/api/tournaments/${stage.tournamentId}/entrants/concludeQualifiers`,
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

		onSuccess: ({ tournament }: { tournament: Tournament }) => {
			if (tournament.status === 'elimination') {
				stage.clear(tournament.id);
			}
		}
	}));
</script>

{#if stage.isQualifying}
	<div class="space-y-3">
		<div>
			<h3 class="font-semibold">Entrants</h3>
			<p class="text-sm text-muted-foreground">Current qualification scores</p>
		</div>

		<div class="divide-y rounded-md border">
			{#each stage.queue as entrant, idx (entrant.id)}
				<button
					type="button"
					class={[
						'flex w-full items-center justify-between px-4 py-3 text-left transition-colors',
						stage.onStage === idx ? 'bg-accent text-accent-foreground' : 'hover:bg-muted/50'
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
										disabled={idx === 0 || entrant.status === 'not_qualified'}
									>
										Move to front
									</DropdownMenuItem>

									<DropdownMenuItem onclick={() => markAs(entrant.id, 'not_qualified')}
										>Set as not qualified</DropdownMenuItem
									>

									<DropdownMenuItem onclick={() => markAs(entrant.id, 'pending')}
										>Set as pending</DropdownMenuItem
									>
								</DropdownMenuContent>
							</DropdownMenu>
						</div>
					</div>
				</button>
			{/each}
		</div>
		{#if stage.isQualifyingComplete()}
			<Button
				class="w-full"
				variant="default"
				onclick={() => concludeQualifiersMutation.mutate()}
				disabled={concludeQualifiersMutation.isPending}
			>
				{concludeQualifiersMutation.isPending ? 'Concluding...' : 'Conclude'}
			</Button>
		{/if}
	</div>
{/if}
