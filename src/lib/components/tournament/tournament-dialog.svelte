<script lang="ts">
	import { createMutation, useQueryClient } from '@tanstack/svelte-query';

	import { Button } from '$lib/components/ui/button';
	import {
		Dialog,
		DialogContent,
		DialogDescription,
		DialogFooter,
		DialogHeader,
		DialogTitle,
		DialogTrigger
	} from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Textarea } from '$lib/components/ui/textarea';

	let open = $state(false);
	let error = $state('');

	let name = $state('');
	let startingDate = $state('');
	let trackName = $state('');
	let qualifierNames = $state('');

	const queryClient = useQueryClient();

	const createTournamentMutation = createMutation(() => ({
		mutationFn: async ({
			name,
			startingDate,
			track,
			qualifiers
		}: {
			name: string;
			startingDate: string;
			track: string;
			qualifiers: string[];
		}) => {
			const response = await fetch('/api/tournaments', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					name,
					startingDate,
					track,
					qualifiers
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
			qualifierNames = '';
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

		const qualifiers = qualifierNames
			.split('\n')
			.map((name) => name.trim())
			.filter(Boolean);

		if (qualifiers.length === 0) {
			error = 'At least one qualifier is required.';
			return;
		}

		await createTournamentMutation.mutateAsync({
			name: name.trim(),
			startingDate,
			track: trackName.trim(),
			qualifiers
		});
	}
</script>

<Dialog bind:open>
	<DialogTrigger>
		<Button>Create Tournament</Button>
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

			<div class="space-y-2">
				<Label for="starting-date">Starting date</Label>

				<Input id="starting-date" type="datetime-local" bind:value={startingDate} />
			</div>

			<div class="space-y-2">
				<Label for="track">Track</Label>

				<Input id="track" placeholder="Enter or select a track" bind:value={trackName} />
			</div>

			<div class="space-y-2">
				<Label for="qualifiers">Qualifiers</Label>

				<Textarea
					id="qualifiers"
					// eslint-disable-next-line svelte/no-useless-mustaches
					placeholder={'One participant per line \nJohn Doe\nJane Doe\nAlex Smith'}
					rows={8}
					bind:value={qualifierNames}
				/>

				<p class="text-sm text-muted-foreground">Enter one participant name per line.</p>
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
				{createTournamentMutation.isPending ? 'Creating...' : 'Create Tournament'}
			</Button>
		</DialogFooter>
	</DialogContent>
</Dialog>
