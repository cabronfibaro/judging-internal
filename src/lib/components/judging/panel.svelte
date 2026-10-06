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

<div class="relative flex flex-col gap-5" aria-busy={judgement.processing}>
	{#if isQualifying || stage.matchId}
		<div class="grid items-start gap-5 xl:grid-cols-[minmax(0,1fr)_18rem]">
			<div class="grid min-w-0 grid-cols-1 gap-3 md:grid-cols-2 md:gap-4">
				{#if common}
					{@const [, category] = common}
					<div class={!isQualifying ? 'md:col-span-2' : ''}>
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

			<div class="min-w-0 xl:sticky xl:top-24">{@render children()}</div>
		</div>
	{/if}

	{#if judgement.processing}
		<div
			class="absolute inset-0 z-50 flex items-center justify-center rounded-xl bg-background/95 px-4"
		>
			<Empty class="w-full max-w-md border-0 bg-transparent shadow-none ring-0">
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
