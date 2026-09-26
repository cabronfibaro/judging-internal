<script lang="ts">
	import Criteriacontrol from './criteria-control.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { Separator } from '$lib/components/ui/separator';
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

<div class="flex items-center justify-between gap-4 py-2">
	<div class="min-w-0">
		<div class="flex items-center gap-3">
			<span class="truncate text-sm font-medium">
				{tag.label}
			</span>

			{#if tag.max !== undefined}
				<span class="shrink-0 text-xs text-muted-foreground tabular-nums">
					{Math.abs((value ?? 0) / (tag.multiplier ?? 1))}/{(tag.max ?? 0) / (tag.multiplier ?? 1)}
				</span>
			{/if}

			{#if tag.type !== 'neutral'}
				<Badge variant={getBadgeVariant(tag.type)} class="shrink-0 text-xs">
					{tag.type}
				</Badge>
			{/if}
		</div>
	</div>

	<Criteriacontrol {tag} {value} onValueChange={(value) => onValueChange(Number(value))} />
</div>
