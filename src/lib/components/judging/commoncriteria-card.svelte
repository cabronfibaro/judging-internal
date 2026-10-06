<script lang="ts">
	import Criteriacontrol from './criteria-control.svelte';
	import { Badge } from '#lib/components/ui/badge';
	import { Card, CardContent, CardHeader, CardTitle } from '#lib/components/ui/card';
	import { Separator } from '#lib/components/ui/separator';
	import type { CriteriaCatrgory } from '../../../types';
	import type { Judgement } from '#lib/judging/judgement.svelte';

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

<Card class="col-span-2 gap-0 border-border/70 shadow-none">
	<CardHeader class="pb-4">
		<CardTitle>{category.label}</CardTitle>
	</CardHeader>
	<Separator />
	<CardContent class="space-y-1 pt-2">
		<!-- Column headers -->

		{#if !isQualifying}
			<div
				class="hidden grid-cols-[minmax(0,1fr)_10rem_10rem] items-center gap-4 px-1 pb-1 md:grid"
			>
				<div></div>
				<div
					class="w-36 justify-self-end text-center text-xs font-medium tracking-wide text-muted-foreground uppercase"
				>
					Lead
				</div>
				<div
					class="w-36 justify-self-end text-center text-xs font-medium tracking-wide text-muted-foreground uppercase"
				>
					Chase
				</div>
			</div>
		{/if}

		{#each category.tags as tag, index (tag.id)}
			{#if index > 0}
				<Separator />
			{/if}

			<div
				class={`grid gap-3 py-3 md:items-center md:gap-4 md:py-2 ${
					isQualifying
						? 'md:grid-cols-[minmax(0,1fr)_10rem]'
						: 'md:grid-cols-[minmax(0,1fr)_10rem_10rem]'
				}`}
			>
				<div class="min-w-0">
					<div class="flex flex-wrap items-center gap-x-3 gap-y-1">
						<span class="text-sm leading-5 font-medium">
							{tag.label}
						</span>

						{#if tag.max !== undefined}
							<span class="shrink-0 text-xs text-muted-foreground tabular-nums">
								{Math.abs((judgement.get(tag.id)?.lead ?? 0) / (tag.multiplier ?? 1))}/{(tag.max ??
									0) / (tag.multiplier ?? 1)}
							</span>
							{#if !isQualifying}
								<span class="shrink-0 text-xs text-muted-foreground md:hidden">•</span>
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
				<div class="flex items-center justify-between gap-3 md:w-36 md:justify-self-end">
					<span class="text-xs font-medium tracking-wide text-muted-foreground uppercase md:hidden"
						>Lead</span
					>
					<Criteriacontrol
						{tag}
						value={judgement.get(tag.id)?.lead}
						onValueChange={(value) => judgement.set(tag.id, 'lead', Number(value))}
					/>
				</div>

				{#if !isQualifying}
					<!-- Chase -->
					<div class="flex items-center justify-between gap-3 md:w-36 md:justify-self-end">
						<span
							class="text-xs font-medium tracking-wide text-muted-foreground uppercase md:hidden"
							>Chase</span
						>
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
