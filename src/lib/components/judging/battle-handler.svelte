<script lang="ts">
	import type { Judgement } from '$lib/judging/judgement.svelte';
	import Button from '../ui/button/button.svelte';
	import CardContent from '../ui/card/card-content.svelte';
	import CardHeader from '../ui/card/card-header.svelte';
	import CardTitle from '../ui/card/card-title.svelte';
	import Card from '../ui/card/card.svelte';
	import DebugJudgement from './debug-judgement.svelte';
	import type { Stage } from '$lib/judging/stage.svelte';

	let {
		judgement,
		stage
	}: {
		judgement: Judgement;
		stage: Stage;
	} = $props();

	async function onSubmitJudgement() {
		const { lead, chase } = judgement.finalizedScore();

		// set score for the run
		// switch leader
		stage.setScore({ lead, chase });
		stage.switchLeader();

		Promise.all([]);

		console.log(stage.state);
		judgement.clear();
	}
</script>

<Card>
	<CardHeader class="pb-3">
		<CardTitle>Judgement</CardTitle>
	</CardHeader>

	<CardContent class="space-y-3">
		<DebugJudgement show results={judgement.results} />
		<Button class="w-full" onclick={onSubmitJudgement}>Submit Judgment</Button>
		<Button variant="destructive" type="reset" class="w-full" onclick={() => judgement.clear()}
			>Reset</Button
		>
	</CardContent>
</Card>
