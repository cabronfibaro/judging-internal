<script lang="ts">
	import CardContent from '../ui/card/card-content.svelte';
	import CardHeader from '../ui/card/card-header.svelte';
	import Card from '../ui/card/card.svelte';
	import Skeleton from '../ui/skeleton/skeleton.svelte';
</script>

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
