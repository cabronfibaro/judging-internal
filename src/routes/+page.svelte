<script lang="ts">
	import Criteriacontrol from '$lib/components/judging/criteriacontrol.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
	import Input from '$lib/components/ui/input/input.svelte';
	import { Label } from '$lib/components/ui/label';

	import { Separator } from '$lib/components/ui/separator';
	import { Judgement } from '$lib/judging/judgement.svelte.js';
	import { Stage } from '$lib/judging/stage.svelte.js';

	// TODO: handle qualification mode by disabling the chase buttons
	// TODO: add note for lead and chase

	let { data } = $props();
	const criteriaConfig = $derived(data.criteriaConfig);

	const judgement = new Judgement();
	const stage = new Stage();

	function onSubmitJudgement() {
		stage.setScore(judgement.finalizedScore());
		judgement.clear();

		if (stage.isConcluded()) {
			console.log(stage.winner());
			return;
		}

		stage.switchLeader();
	}

	function getBadgeVariant(type: string) {
		if (type === 'bonus') return 'default';
		if (type === 'deduction') return 'destructive';
		return 'secondary';
	}

	const categories = $derived(Object.entries(criteriaConfig.categories));
	let common = $derived(categories.find(([id]) => id === 'common'));
	let lead = $derived(categories.find(([id]) => id === 'lead'));
	let chase = $derived(categories.find(([id]) => id === 'chase'));
</script>

<div class="mx-auto max-w-4xl space-y-4 p-6">
	<div>
		<h1 class="text-2xl font-bold">Using: {criteriaConfig.name} config</h1>
		<p class="text-muted-foreground">Lorem ipsum jembut delemut</p>
	</div>

	<Card>
		<CardHeader>
			<CardTitle>Stage</CardTitle>
		</CardHeader>

		<CardContent class="space-y-4">
			<div class="grid grid-cols-2 gap-4">
				<!-- P1 -->
				<div class="space-y-3">
					<div class="flex items-center justify-between">
						<Label>P1</Label>

						<Badge variant={stage.state.p1.state === 'lead' ? 'default' : 'secondary'}>
							{stage.state.p1.state}
						</Badge>
					</div>
					<Input
						placeholder="Driver name"
						value={stage.state.p1.name}
						oninput={(event) => stage.updateDriver('p1', 'name', event.currentTarget.value)}
					/>
				</div>

				<!-- P2 -->
				<div class="space-y-3">
					<div class="flex items-center justify-between">
						<Label>P2</Label>

						<Badge variant={stage.state.p2.state === 'lead' ? 'default' : 'secondary'}>
							{stage.state.p2.state}
						</Badge>
					</div>

					<Input
						placeholder="Driver name"
						value={stage.state.p2.name}
						oninput={(event) => stage.updateDriver('p2', 'name', event.currentTarget.value)}
					/>
				</div>
			</div>
		</CardContent>
	</Card>

	<div class="grid grid-cols-2 gap-4">
		<!-- Common -->
		{#if common}
			{@const [, category] = common}

			<Card class="col-span-2">
				<CardHeader class="pb-3">
					<CardTitle>{category.label}</CardTitle>
				</CardHeader>

				<CardContent class="space-y-3">
					<div class="grid grid-cols-[1fr_224px_224px] items-center gap-4 px-1">
						<div></div>
						<div class="w-[144px] justify-self-end text-center text-sm font-medium">Lead</div>
						<div class="w-[144px] justify-self-end text-center text-sm font-medium">Chase</div>
					</div>

					{#each category.tags as tag, index (tag.id)}
						{#if index > 0}
							<Separator />
						{/if}

						<div class="grid grid-cols-[1fr_224px_224px] items-center gap-4 py-2">
							<!-- Criterion -->
							<div class="min-w-0">
								<div class="flex items-center gap-2">
									<Label class="text-sm font-medium">
										{tag.label}
									</Label>

									<Badge variant={getBadgeVariant(tag.type)} class="text-xs">
										{tag.type}
									</Badge>
								</div>
							</div>

							<div class="w-[144px] justify-self-end">
								<Criteriacontrol
									{tag}
									value={judgement.get(tag.id)?.lead}
									onValueChange={(value) => judgement.set(tag.id, 'lead', Number(value))}
								/>
							</div>

							<div class="w-[144px] justify-self-end">
								<Criteriacontrol
									{tag}
									value={judgement.get(tag.id)?.chase}
									onValueChange={(value) => judgement.set(tag.id, 'chase', Number(value))}
								/>
							</div>
						</div>
					{/each}
				</CardContent>
			</Card>
		{/if}

		<!-- Lead -->
		{#if lead}
			{@const [, category] = lead}
			<Card>
				<CardHeader class="pb-3">
					<CardTitle>{category.label}</CardTitle>
				</CardHeader>

				<CardContent class="space-y-3">
					{#each category.tags as tag, index (tag.id)}
						{#if index > 0}
							<Separator />
						{/if}

						<div class="flex items-center justify-between space-y-2 py-2">
							<div class="flex items-center gap-2">
								<Label class="text-sm font-medium">
									{tag.label}
								</Label>

								<Badge variant={getBadgeVariant(tag.type)} class="text-xs">
									{tag.type}
								</Badge>
							</div>

							<Criteriacontrol
								{tag}
								value={judgement.get(tag.id)?.lead}
								onValueChange={(value) => judgement.set(tag.id, 'lead', Number(value))}
							/>
						</div>
					{/each}
				</CardContent>
			</Card>
		{/if}

		<!-- Chase -->
		{#if chase}
			{@const [, category] = chase}
			<Card>
				<CardHeader class="pb-3">
					<CardTitle>{category.label}</CardTitle>
				</CardHeader>

				<CardContent class="space-y-3">
					{#each category.tags as tag, index (tag.id)}
						{#if index > 0}
							<Separator />
						{/if}

						<div class="flex items-center justify-between space-y-2 py-2">
							<div class="flex items-center gap-2">
								<Label class="text-sm font-medium">
									{tag.label}
								</Label>

								<Badge variant={getBadgeVariant(tag.type)} class="text-xs">
									{tag.type}
								</Badge>
							</div>

							<Criteriacontrol
								{tag}
								value={judgement.get(tag.id)?.chase}
								onValueChange={(value) => judgement.set(tag.id, 'chase', Number(value))}
							/>
						</div>
					{/each}
				</CardContent>
			</Card>
		{/if}
	</div>

	<Card>
		<CardHeader class="pb-3">
			<CardTitle>Judgment</CardTitle>
		</CardHeader>

		<CardContent class="space-y-3">
			<pre class="overflow-auto rounded-md bg-muted p-3 text-sm">{JSON.stringify(
					judgement.results,
					null,
					2
				)}</pre>

			<Button class="w-full" onclick={onSubmitJudgement}>Submit Judgment</Button>
			<Button variant="destructive" type="reset" class="w-full" onclick={() => judgement.clear()}
				>Reset</Button
			>
		</CardContent>
	</Card>
</div>
