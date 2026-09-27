<script lang="ts">
	import Criteriacontrol from './criteria-control.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Separator } from '$lib/components/ui/separator';
	import type { CriteriaCatrgory } from '../../../types';
	import type { Judgement } from '$lib/judging/judgement.svelte';

	let {
		category,
		judgement,
		isQualifying
	}: {
		category: CriteriaCatrgory;
		judgement: Judgement;
		isQualifying: boolean;
	} = $props();

	function getBadgeVariant(type: string) {
		if (type === 'bonus') return 'default';
		if (type === 'deduction') return 'destructive';
		return 'secondary';
	}
</script>

<Card class="col-span-2">
	<CardHeader>
		<CardTitle>{category.label}</CardTitle>
	</CardHeader>
	<Separator />
	<CardContent class="space-y-3">
		<!-- Column headers -->

		{#if !isQualifying}
			<div class="grid grid-cols-[1fr_224px_224px] items-center gap-4 px-1">
				<div></div>
				<div class="w-36 justify-self-end text-center text-sm font-medium">Lead</div>

				<div class="w-36 justify-self-end text-center text-sm font-medium">Chase</div>
			</div>
		{/if}

		{#each category.tags as tag, index (tag.id)}
			{#if index > 0}
				<Separator />
			{/if}

			<div
				class:grid-cols-[1fr_224px]={isQualifying}
				class:grid-cols-[1fr_224px_224px]={!isQualifying}
				class="grid items-center gap-4 py-2"
			>
				<div class="min-w-0">
					<div class="flex items-center gap-3">
						<span class="truncate text-sm font-medium">
							{tag.label}
						</span>

						{#if tag.max !== undefined}
							<span class="shrink-0 text-xs text-muted-foreground tabular-nums">
								{Math.abs((judgement.get(tag.id)?.lead ?? 0) / (tag.multiplier ?? 1))}/{(tag.max ??
									0) / (tag.multiplier ?? 1)}
							</span>
							{#if !isQualifying}
								<span class="shrink-0 text-xs text-muted-foreground tabular-nums"> | </span>
								<span class="shrink-0 text-xs text-muted-foreground tabular-nums">
									{Math.abs(
										(judgement.get(tag.id)?.chase ?? 0) / (tag.multiplier ?? 1)
									)}/{(tag.max ?? 0) / (tag.multiplier ?? 1)}
								</span>
							{/if}
						{/if}

						{#if tag.type !== 'neutral'}
							<Badge variant={getBadgeVariant(tag.type)} class="shrink-0 text-xs">
								{tag.type}
							</Badge>
						{/if}
					</div>
				</div>

				<!-- Lead -->
				<div class="w-36 justify-self-end">
					<Criteriacontrol
						{tag}
						value={judgement.get(tag.id)?.lead}
						onValueChange={(value) => judgement.set(tag.id, 'lead', Number(value))}
					/>
				</div>

				{#if !isQualifying}
					<!-- Chase -->
					<div class="w-36 justify-self-end">
						<Criteriacontrol
							{tag}
							value={judgement.get(tag.id)?.chase}
							onValueChange={(value) => judgement.set(tag.id, 'chase', Number(value))}
						/>
					</div>
				{/if}
			</div>
		{/each}
	</CardContent>
</Card>
