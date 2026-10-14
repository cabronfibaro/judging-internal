<script lang="ts">
  import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
  import { SvelteQueryDevtools } from "@tanstack/svelte-query-devtools";
  import { ModeWatcher } from "mode-watcher";
  import Router from "svelte-spa-router";
  import { wrap } from "svelte-spa-router/wrap";

  import Home from "./components/pages/home.svelte";
  import ChallongeDialog from "./components/challonge-dialog.svelte";
  import { HasChallongeAPIKey } from "$wails/go/main/App";

  const routes = {
    "/": Home,
    "/tournaments/:id": wrap({
      asyncComponent: () => import("./components/pages/tournament.svelte"),
    }),
    "/criteria-configs": wrap({
      asyncComponent: () => import("./components/pages/criteria-config.svelte"),
    }),
  };

  const queryClient = new QueryClient();

  let dialogOpen = $state(false);
  $effect(() => {
    let cancelled = false;

    const checkForKey = async () => {
      try {
        const hasChallongeKey = await HasChallongeAPIKey();

        if (!cancelled) {
          dialogOpen = !hasChallongeKey;
        }
      } catch (error) {
        console.error("Failed to check Challonge API key:", error);

        if (!cancelled) {
          dialogOpen = true;
        }
      }
    };

    checkForKey();

    return () => {
      cancelled = true;
    };
  });
</script>

<ModeWatcher defaultMode="dark" />
<QueryClientProvider client={queryClient}>
  <SvelteQueryDevtools />
  <ChallongeDialog bind:open={dialogOpen} />
  <Router {routes} />
</QueryClientProvider>
