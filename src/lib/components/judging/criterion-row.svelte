<script lang="ts">
	import Criteriacontrol from './criteria-control.svelte';
	import { Badge } from '#lib/components/ui/badge';
	import { Separator } from '#lib/components/ui/separator';
	import type { Tag } from '../../../types';

	let {
		tag,
		value,
		onValueChange,
		showSeparator = false
	}: {
		tag: Tag;
		value: number | undefined;
		onValueChange: (value: number) => void;
		showSeparator?: boolean;
	} = $props();

	function getBadgeVariant(type: string) {
		if (type === 'bonus') return 'default';
		if (type === 'deduction') return 'destructive';
		return 'secondary';
	}
</script>

{#if showSeparator}
	<Separator />
{/if}

<div class="flex flex-col items-stretch gap-3 py-3 sm:flex-row sm:items-center sm:justify-between sm:gap-4">
	<div class="min-w-0">
		<div class="flex flex-wrap items-center gap-x-3 gap-y-1">
			<span class="text-sm font-medium leading-5">
				{tag.label}
			</span>

			{#if tag.max !== undefined}
				<span class="shrink-0 text-xs text-muted-foreground tabular-nums">
					{Math.abs((value ?? 0) / (tag.multiplier ?? 1))}/{(tag.max ?? 0) / (tag.multiplier ?? 1)}
				</span>
			{/if}

			{#if tag.type !== 'neutral'}
				<Badge variant={getBadgeVariant(tag.type)} class="shrink-0 text-xs capitalize">
					{tag.type}
				</Badge>
			{/if}
		</div>
	</div>

	<div class="shrink-0 self-end sm:self-center">
		<Criteriacontrol {tag} {value} onValueChange={(value) => onValueChange(Number(value))} />
	</div>
</div>
