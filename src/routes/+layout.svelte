<script lang="ts">
	import { ModeWatcher } from 'mode-watcher';
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
	import './layout.css';
	import Main from '#lib/components/layout/main.svelte';
	import { browser } from '$app/env';

	let { children } = $props();
	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				enabled: browser
			}
		}
	});
</script>

<ModeWatcher defaultMode="dark" />
<QueryClientProvider client={queryClient}>
	<Main>
		{@render children()}
	</Main>
</QueryClientProvider>
