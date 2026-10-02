<script lang="ts">
	import type { CriteriaCatrgory } from '../../../types';
	import CommoncriteriaCard from './commoncriteria-card.svelte';
	import CriteriaCard from './criteria-card.svelte';
	import type { Judgement } from '$lib/judging/judgement.svelte';
	import type { Snippet } from 'svelte';
	import type { Stage } from '$lib/judging/stage.svelte';

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
