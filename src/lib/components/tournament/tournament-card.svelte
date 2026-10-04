<script lang="ts">
	import type { Tournament } from '../../../types/tournament';
	import Button from '../ui/button/button.svelte';
	import CardContent from '../ui/card/card-content.svelte';
	import CardHeader from '../ui/card/card-header.svelte';
	import CardTitle from '../ui/card/card-title.svelte';
	import Card from '../ui/card/card.svelte';
	import Progress from '../ui/progress/progress.svelte';
	import Skeleton from '../ui/skeleton/skeleton.svelte';

	let {
		tournament,
		openDeleteDialog
	}: { tournament: Tournament; openDeleteDialog(tournament: Tournament): void } = $props();
</script>

<Card>
	<CardHeader>
		<CardTitle>
			{tournament.name} ({tournament.status === 'qualifying' ? 'Qualifying' : 'Elimination'})
		</CardTitle>
	</CardHeader>

	<CardContent class="space-y-4">
		<Skeleton class="aspect-video w-full animate-none rounded-md" />

		<div class="flex flex-col gap-4 space-y-1 text-sm">
			<div>
				<p>
					<span class="text-muted-foreground">Track:</span>
					{tournament.track ?? 'Unknown'}
				</p>

				<p>
					<span class="text-muted-foreground">Starting:</span>
					{new Date(tournament.startingDate).toLocaleString()}
				</p>

				<p>
					<span class="text-muted-foreground">Bracket:</span>
					<a
						href={`https://challonge.com/${tournament.url}`}
						target="_blank"
						rel="noopener noreferrer"
						class="text-primary underline underline-offset-4 hover:text-primary/80"
					>
						Link
					</a>
				</p>
			</div>
			<Progress value={tournament.progress * 100} />
		</div>

		<div class="flex gap-2">
			<Button href={`/tournaments/${tournament.id}`}>Load</Button>
			<Button variant="destructive" onclick={() => openDeleteDialog(tournament)}>Delete</Button>
		</div>
	</CardContent>
</Card>
