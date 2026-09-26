<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import {
		Select,
		SelectContent,
		SelectItem,
		SelectTrigger,
		SelectValue
	} from '$lib/components/ui/select';
	import type { Tag } from '../../../types';
	import Switch from '../ui/switch/switch.svelte';

	let {
		tag,
		value = 0,
		onValueChange
	}: {
		tag: Tag;
		value?: unknown;
		onValueChange: (value: unknown) => void;
	} = $props();

	function changeValue(amount: number, crit_type: (typeof tag)['type']) {
		const current = Number(value ?? 0);

		let multiplier = 1;
		if (tag.multiplier) {
			multiplier = tag.multiplier;
		}

		amount = amount * multiplier;
		const final_amount = crit_type === 'deduction' ? -amount : amount;
		onValueChange(current + final_amount);
	}

	let holdInterval: ReturnType<typeof setInterval> | undefined;
	let holdDelay = 150;

	function startHold(amount: number, crit_type: (typeof tag)['type']) {
		changeValue(amount, crit_type);

		holdDelay = 150;

		function tick() {
			changeValue(amount, crit_type);

			holdDelay = Math.max(30, holdDelay - 10);

			holdInterval = setTimeout(tick, holdDelay);
		}

		holdInterval = setTimeout(tick, holdDelay);
	}

	function stopHold() {
		if (holdInterval) {
			clearTimeout(holdInterval);
			holdInterval = undefined;
		}
	}
</script>

{#if tag.value_type === 'enum'}
	<Select type="single" onValueChange={(selected) => onValueChange(Number(selected))}>
		<SelectTrigger>
			<SelectValue placeholder="Select" />
		</SelectTrigger>

		<SelectContent>
			{#each tag.options ?? [] as option (option.value)}
				<SelectItem value={String(option.value)}>
					{option.label}
				</SelectItem>
			{/each}
		</SelectContent>
	</Select>
{:else if tag.value_type === 'dynamic'}
	<div class="flex justify-center">
		<Button
			variant="outline"
			size="icon"
			class="h-9 w-9 rounded-r-none"
			onpointerdown={() => startHold(-1, tag.type)}
			onpointerup={stopHold}
			onpointerleave={stopHold}
			onpointercancel={stopHold}
			disabled={Math.abs(Number(value ?? 0)) <= 0}
		>
			−
		</Button>

		<Input
			autocomplete="off"
			class="h-9 w-12 [appearance:textfield] rounded-none text-center [&::-webkit-inner-spin-button]:appearance-none [&::-webkit-outer-spin-button]:appearance-none"
			type="number"
			disabled
			value={Math.abs(Number(value ?? 0))}
			oninput={(event) => {
				onValueChange(Math.max(0, Number(event.currentTarget.value)));
			}}
		/>

		<Button
			variant="outline"
			size="icon"
			class="h-9 w-9 rounded-l-none"
			onpointerdown={() => startHold(1, tag.type)}
			onpointerup={stopHold}
			onpointerleave={stopHold}
			onpointercancel={stopHold}
			disabled={tag.max !== undefined && Math.abs(Number(value)) >= tag.max}
		>
			+
		</Button>
	</div>
{:else if tag.value_type === 'switch'}
	<Switch
		checked={Number(value ?? 0) !== 0}
		onCheckedChange={(checked) => {
			const amount = Number(tag.default_value ?? 0);

			if (checked) {
				onValueChange(Number(value ?? 0) + amount);
			} else {
				onValueChange(Number(value ?? 0) - amount);
			}
		}}
	/>
{/if}
