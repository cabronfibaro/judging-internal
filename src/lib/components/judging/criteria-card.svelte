<script lang="ts">
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
	import type { Judgement } from '$lib/judging/judgement.svelte';
	import type { CriteriaCatrgory } from '../../../types';
	import Separator from '../ui/separator/separator.svelte';
	import CriterionRow from './criterion-row.svelte';

	let {
		category,
		role,
		judgement
	}: {
		category: CriteriaCatrgory;
		role: 'lead' | 'chase';
		judgement: Judgement;
	} = $props();
</script>

<Card>
	<CardHeader>
		<CardTitle>{category.label}</CardTitle>
	</CardHeader>
	<Separator />
	<CardContent class="space-y-3">
		{#each category.tags as tag, index (tag.id)}
			<CriterionRow
				{tag}
				value={judgement.get(tag.id)?.[role]}
				showSeparator={index > 0}
				onValueChange={(value) => judgement.set(tag.id, role, value)}
			/>
		{/each}
	</CardContent>
</Card>
