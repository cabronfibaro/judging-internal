<script lang="ts">
	import CommoncriteriaCard from '$lib/components/judging/commoncriteria-card.svelte';
	import CriteriaCard from '$lib/components/judging/criteria-card.svelte';
	import DebugJudgement from '$lib/components/judging/debug-judgement.svelte';
	import StageCard from '$lib/components/judging/stage-card.svelte';
	import TournamentDialog from '$lib/components/tournament/tournament-dialog.svelte';
	import TournamentGrid from '$lib/components/tournament/tournament-grid.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';

	import { Judgement } from '$lib/judging/judgement.svelte.js';
	import { Stage } from '$lib/judging/stage.svelte.js';

	// TODO: handle qualification mode by disabling the chase buttons
	// TODO: add note for lead and chase

	let { data } = $props();
	const criteriaConfig = $derived(data.criteriaConfig);

	let judgement = new Judgement();
	let stage = new Stage();

	function onSubmitJudgement() {
		const { lead, chase } = judgement.finalizedScore();
		stage.setScore({ lead, chase });
		judgement.clear();

		if (!stage.isQualifying) {
			handleBattle();
			return;
		}
	}

	function handleBattle() {
		if (stage.isConcluded()) {
			console.log(stage.winner());
			return;
		}

		stage.switchLeader();
	}

	const categories = $derived(Object.entries(criteriaConfig.categories));
	let common = $derived(categories.find(([id]) => id === 'common'));
	let lead = $derived(categories.find(([id]) => id === 'lead'));
	let chase = $derived(categories.find(([id]) => id === 'chase'));
</script>

<div class="mx-auto max-w-4xl space-y-4 p-6">
	<TournamentDialog />
	<TournamentGrid {stage} />
	<div>
		<h1 class="text-2xl font-bold">Using: {criteriaConfig.name} config</h1>
		<p class="text-muted-foreground">Lorem ipsum jembut delemut</p>
	</div>

	<StageCard {stage} />

	<div
		class="grid"
		class:grid-cols-1={stage.isQualifying}
		class:gap-0={stage.isQualifying}
		class:gap-y-4={stage.isQualifying}
		class:gap-4={!stage.isQualifying}
		class:grid-cols-2={!stage.isQualifying}
	>
		{#if common}
			{@const [, category] = common}
			<CommoncriteriaCard {category} {judgement} isQualifying={stage.isQualifying} />
		{/if}

		{#if lead}
			{@const [, category] = lead}
			<CriteriaCard {category} role="lead" {judgement} />
		{/if}

		{#if chase && !stage.isQualifying}
			{@const [, category] = chase}
			<CriteriaCard {category} role="lead" {judgement} />
		{/if}
	</div>

	<Card>
		<CardHeader class="pb-3">
			<CardTitle>Judgement</CardTitle>
		</CardHeader>

		<!-- TODO: show score between participants -->

		<CardContent class="space-y-3">
			<DebugJudgement show results={judgement.results} />
			<Button class="w-full" onclick={onSubmitJudgement}>Submit Judgment</Button>
			<Button variant="destructive" type="reset" class="w-full" onclick={() => judgement.clear()}
				>Reset</Button
			>
		</CardContent>
	</Card>
</div>
