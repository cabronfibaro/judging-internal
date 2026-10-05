<script lang="ts">
	import type { Snippet } from 'svelte';
	import { cn } from '#lib/utils';

	let {
		fixed = false,
		fluid = false,
		class: className = '',
		children
	}: {
		fixed?: boolean;
		fluid?: boolean;
		class?: string;
		children?: Snippet;
	} = $props();
</script>

<main
	data-layout={fixed ? 'fixed' : 'auto'}
	class={cn(
		'px-4 py-6',
		fixed && 'flex grow flex-col overflow-hidden',
		!fluid && '@7xl/content:mx-auto @7xl/content:w-full @7xl/content:max-w-7xl',
		className
	)}
>
	{@render children?.()}
</main>
