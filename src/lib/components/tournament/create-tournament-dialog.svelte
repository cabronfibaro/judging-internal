<script lang="ts">
	import { createMutation, useQueryClient } from '@tanstack/svelte-query';
	import { Plus } from '@lucide/svelte';
	import { Button } from '#lib/components/ui/button';
	import {
		Dialog,
		DialogContent,
		DialogDescription,
		DialogFooter,
		DialogHeader,
		DialogTitle,
		DialogTrigger
	} from '#lib/components/ui/dialog';
	import { Input } from '#lib/components/ui/input';
	import { Label } from '#lib/components/ui/label';
	import { Textarea } from '#lib/components/ui/textarea';
	import Switch from '../ui/switch/switch.svelte';
	let { open = $bindable(false) } = $props();
	let error = $state('');
	let name = $state('');
	let size = $state(2);
	let startingDate = $state('');
	let trackName = $state('');
	let entrantNames = $state('');
	let randomizeQualifyingOrder = $state(true);
	const queryClient = useQueryClient();
	const createTournamentMutation = createMutation(() => ({
		mutationFn: async ({
			name,
			size,
			startingDate,
			track,
			entrants,
			randomizeQualifyingOrder
		}: {
			name: string;
			size: number;
			startingDate: string;
			track: string;
			entrants: string[];
			randomizeQualifyingOrder: boolean;
		}) => {
			const response = await fetch('/api/tournaments', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					name,
					size,
					startingDate,
					track,
					entrants,
					randomizeQualifyingOrder
				})
			});
			const result = await response.json();
			if (!response.ok) {
				throw new Error(result.message ?? 'Failed to create tournament.');
			}
			return result;
		},
		onSuccess: async () => {
			await queryClient.invalidateQueries({
				queryKey: ['tournaments']
			});
			open = false;
			name = '';
			startingDate = '';
			trackName = '';
			entrantNames = '';
			error = '';
		}
	}));
	async function createTournament() {
		error = '';
		if (!name.trim()) {
			error = 'Tournament name is required.';
			return;
		}
		if (!startingDate) {
			error = 'Starting date is required.';
			return;
		}
		if (!trackName.trim()) {
			error = 'Track is required.';
			return;
		}
		const entrants = entrantNames
			.split('\n')
			.map((name) => name.trim())
			.filter(Boolean);
		if (entrants.length === 0) {
			error = 'At least one qualifier is required.';
			return;
		}
		await createTournamentMutation.mutateAsync({
			name: name.trim(),
			size,
			startingDate,
			track: trackName.trim(),
			entrants,
			randomizeQualifyingOrder
		});
	}
</script>

<Dialog bind:open>
	<DialogTrigger>
		{#snippet child({ props })}
			<Button {...props}>
				<Plus />
				New tournament
			</Button>
		{/snippet}
	</DialogTrigger>
	<DialogContent>
		<DialogHeader>
			<DialogTitle>Create Tournament</DialogTitle>
			<DialogDescription>
				Create a tournament and add its qualifying participants.
			</DialogDescription>
		</DialogHeader>
		<div class="space-y-4">
			<div class="space-y-2">
				<Label for="name">Tournament name</Label>
				<Input id="name" placeholder="Summer Drift Championship" bind:value={name} />
			</div>
			<div class="grid gap-4 sm:grid-cols-2">
				<div class="space-y-2">
					<Label for="size">Size</Label>
					<Input
						type="number"
						min="2"
						id="size"
						placeholder="Size of the bracket"
						bind:value={size}
					/>
				</div>
				<div class="space-y-2">
					<Label for="starting-date">Starting date</Label>
					<Input id="starting-date" type="datetime-local" bind:value={startingDate} />
				</div>
			</div>
			<div class="space-y-2">
				<Label for="track">Track</Label>
				<Input id="track" placeholder="Enter or select a track" bind:value={trackName} />
			</div>
			<div class="space-y-2">
				<Label for="Entrants">Entrants</Label>
				<Textarea
					id="Entrants"
					// eslint-disable-next-line svelte/no-useless-mustaches
					placeholder={'One participant per line \nJohn Doe\nJane Doe\nAlex Smith'}
					rows={8}
					bind:value={entrantNames}
				/>
				<p class="text-sm text-muted-foreground">Enter one participant name per line.</p>
			</div>
			<div class="flex items-center justify-between gap-4 rounded-lg border border-border px-3 py-2.5">
				<Label>Randomize Qualifying Order</Label>
				<Switch bind:checked={randomizeQualifyingOrder} />
			</div>
			{#if error}
				<p class="text-sm text-destructive">
					{error}
				</p>
			{/if}
			{#if createTournamentMutation.isError}
				<p class="text-sm text-destructive">
					{createTournamentMutation.error.message}
				</p>
			{/if}
		</div>
		<DialogFooter>
			<Button disabled={createTournamentMutation.isPending} onclick={createTournament}>
				{createTournamentMutation.isPending ? 'Creating...' : 'Create tournament'}
			</Button>
		</DialogFooter>
	</DialogContent>
</Dialog>
