<script lang="ts">
	import type { CriteriaCatrgory } from '../../../types';
	import CommoncriteriaCard from './commoncriteria-card.svelte';
	import CriteriaCard from './criteria-card.svelte';
	import type { Judgement } from '#lib/judging/judgement.svelte';
	import type { Snippet } from 'svelte';
	import type { Stage } from '#lib/judging/stage.svelte';

	import Empty from '#lib/components/ui/empty/empty.svelte';
	import EmptyHeader from '#lib/components/ui/empty/empty-header.svelte';
	import EmptyMedia from '#lib/components/ui/empty/empty-media.svelte';
	import EmptyTitle from '#lib/components/ui/empty/empty-title.svelte';
	import EmptyDescription from '#lib/components/ui/empty/empty-description.svelte';
	import Spinner from '#lib/components/ui/spinner/spinner.svelte';

	let {
		judgement,
		categories,
		isQualifying,
		stage,
		children
	}: {
		judgement: Judgement;
		categories: [string, CriteriaCatrgory][];
		isQualifying: boolean;
		stage: Stage;
		children: Snippet<[]>;
	} = $props();

	let common = $derived(categories.find(([id]) => id === 'common'));
	let lead = $derived(categories.find(([id]) => id === 'lead'));
	let chase = $derived(categories.find(([id]) => id === 'chase'));
</script>

<div class="relative flex flex-col gap-y-4">
	{#if isQualifying || stage.matchId}
		<div class="grid grid-cols-2 gap-4">
			{#if common}
				{@const [, category] = common}
				<div class:col-span-2={!isQualifying}>
					<CommoncriteriaCard {category} {judgement} {isQualifying} />
				</div>
			{/if}

			{#if lead}
				{@const [, category] = lead}
				<CriteriaCard {category} role="lead" {judgement} />
			{/if}

			{#if chase && !isQualifying}
				{@const [, category] = chase}
				<CriteriaCard {category} role="chase" {judgement} />
			{/if}
		</div>

		{@render children()}
	{/if}

	{#if judgement.processing}
		<div
			class="absolute -inset-1 z-50 flex items-center justify-center rounded-lg bg-background/80 backdrop-blur-sm"
		>
			<Empty class="w-full border-0 shadow-none ring-0">
				<EmptyHeader>
					<EmptyMedia variant="icon">
						<Spinner />
					</EmptyMedia>

					<EmptyTitle>Processing the scores</EmptyTitle>
					<EmptyDescription>
						Please wait while we process your request. Do not refresh the page.
					</EmptyDescription>
				</EmptyHeader>
			</Empty>
		</div>
	{/if}
</div>
