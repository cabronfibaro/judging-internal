<script lang="ts">
  import { mode, toggleMode } from "mode-watcher";
  import Moon from "@lucide/svelte/icons/moon";
  import Sun from "@lucide/svelte/icons/sun";
  import Button from "$src/components/ui/button/button.svelte";

  let current = $derived(mode.current);
</script>

<Button
  variant="ghost"
  size="icon"
  onclick={toggleMode}
  aria-label={current === "dark"
    ? "Switch to light mode"
    : "Switch to dark mode"}
  title={current === "dark" ? "Switch to light mode" : "Switch to dark mode"}
>
  {#if current === "dark"}
    <Sun aria-hidden="true" />
  {:else}
    <Moon aria-hidden="true" />
  {/if}
</Button>
