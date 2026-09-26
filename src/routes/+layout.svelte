<script lang="ts">
	import { ModeWatcher } from 'mode-watcher';
	import './layout.css';

	let { children } = $props();
</script>

<ModeWatcher defaultMode="dark" />
{@render children()}
