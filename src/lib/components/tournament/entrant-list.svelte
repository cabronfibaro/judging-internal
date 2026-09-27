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
	import type { Entrants } from '../../../types/tournament';

	let { stage }: { stage: Stage } = $props();

	const updateTournamentMutation = createMutation(() => ({
		mutationFn: async ({
			tournamentId,
			qualifyingOrder
		}: {
			tournamentId: string;
			qualifyingOrder: number[];
		}) => {
			const qualifyingOrderStr = qualifyingOrder.join(', ');
			const response = await fetch(`/api/tournaments/${tournamentId}`, {
				method: 'PATCH',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ qualifyingOrder: qualifyingOrderStr })
			});

			const data = await response.json();

			if (!response.ok) {
				throw new Error(data.message ?? 'Failed to update tournament.');
			}

			return data;
		}
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
			status: 'qualified' | 'not_qualified' | 'pending';
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

	function markAs(id: string, status: Entrants['status']) {
		updateQualifierMutation.mutate({
			qualifierId: id,
			status: status
		});
	}
</script>

{#if stage.isQualifying}
	<div class="space-y-3">
		<div>
			<h3 class="font-semibold">Entrants</h3>
			<p class="text-sm text-muted-foreground">Current qualification scores</p>
		</div>

		<div class="divide-y rounded-md border">
			{#each stage.queue as qualifier, idx (qualifier.id)}
				<button
					type="button"
					class={[
						'flex w-full items-center justify-between px-4 py-3 text-left transition-colors',
						stage.onStage === idx ? 'bg-accent text-accent-foreground' : 'hover:bg-muted/50'
					]}
				>
					<div class="flex items-center gap-3">
						<span class="font-medium">
							{qualifier.name}
						</span>

						{#if qualifier.status === 'pending'}
							<Badge variant="secondary">{idx === 0 ? 'On Going' : 'Pending'}</Badge>
						{:else if qualifier.status === 'qualified'}
							<Badge>Qualified</Badge>
						{:else}
							<Badge variant="destructive">Not Qualified</Badge>
						{/if}
					</div>

					<div class="flex items-center gap-3">
						<span class="text-lg font-semibold tabular-nums">
							{qualifier.score}
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
										disabled={idx === 0 || qualifier.status === 'not_qualified'}
									>
										Move to front
									</DropdownMenuItem>

									<DropdownMenuItem onclick={() => markAs(qualifier.id, 'not_qualified')}
										>Set as not qualified</DropdownMenuItem
									>

									<DropdownMenuItem onclick={() => markAs(qualifier.id, 'pending')}
										>Set as pending</DropdownMenuItem
									>
								</DropdownMenuContent>
							</DropdownMenu>
						</div>
					</div>
				</button>
			{/each}
		</div>
	</div>
{/if}
